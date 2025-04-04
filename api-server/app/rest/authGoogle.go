package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/common"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google config:
// https://console.cloud.google.com/auth/clients?highlightClient=750385423567-rsrv4dknuvrts9rv5neab3dl667r5la6.apps.googleusercontent.com&authuser=2&organizationId=0&project=htl-ride-sharing

const oauthUrlAPIGoogle = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  utils.GetEnvRequired(common.ENV_GOOGLE_REDIRECT_URL),
	ClientID:     utils.GetEnvRequired(common.ENV_GOOGLE_CLIENT_ID),
	ClientSecret: utils.GetEnvRequired(common.ENV_GOOGLE_CLIENT_SECRET),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func authHandlersGoogle(h *http.ServeMux) {
	h.HandleFunc("GET /auth/google/login", oauthLoginGoogle)
	h.HandleFunc("GET /auth/google/callback", oauthCallbackGoogle)
}

func oauthLoginGoogle(w http.ResponseWriter, r *http.Request) {
	oauthState := genRandBase64(64)

	state.mutex.Lock()
	state.oauthStates[oauthState] = time.Now()
	state.mutex.Unlock()

	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthCallbackGoogle(w http.ResponseWriter, r *http.Request) {
	oauthState := r.FormValue("state")
	state.mutex.Lock()
	createdAt, ok := state.oauthStates[oauthState]
	if !ok || time.Now().Sub(createdAt) > time.Duration(5*time.Minute) {
		state.mutex.Unlock()
		http.Error(w, "Invalid 'state' parameter. Make sure authentication request are only started from '/auth/google/login'.", http.StatusBadRequest)
		return
	}
	delete(state.oauthStates, oauthState)

	for oauthState, createdAt := range state.oauthStates {
		elapsed := time.Now().Sub(createdAt)
		if elapsed > time.Duration(5*time.Minute) {
			delete(state.oauthStates, oauthState)
		}
	}

	state.mutex.Unlock()

	profile, err := getUserProfileFromGoogle(r.Context(), r.FormValue("code"))
	if err != nil {
		log.Println("Error: Failed to get Google user profile.", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !*profile.VerifiedEmail {
		http.Error(w, "Google user has an unverified email address. This is not allowed.", http.StatusBadRequest)
		return
	}

	user, err := state.queries.UsersGetById(r.Context(), *profile.Id)
	if err == sql.ErrNoRows {
		user, err = state.queries.UsersCreate(r.Context(), sqlc.UsersCreateParams{ID: *profile.Id, Name: *profile.Name, Email: *profile.Email, Provider: "google"})
		if err != nil {
			log.Println("Error: Failed to create user after Google authentication.", err.Error(), "user email:", *profile.Email)
			http.Error(w, "Failed to create user after Google authentication.", http.StatusInternalServerError)
			return
		}

		tokens := GenAuthTokens(user.ID, user.Email)

		args := sqlc.UsersSetTokensParams{ID: user.ID, AccessToken: sql.NullString{String: tokens.AccessToken}, RefreshToken: sql.NullString{String: tokens.RefreshToken}}
		err = state.queries.UsersSetTokens(r.Context(), args)
		if err != nil {
			log.Println("Error: Failed to update user tokens.", err.Error(), "user id:", *profile.Id, "user email:", *profile.Email)
			http.Error(w, "Failed to update user data.", http.StatusInternalServerError)
			return
		}

		url := fmt.Sprintf("%s?accessToken=%s&refreshToken=%s", clientUrlAuth, tokens.AccessToken, tokens.RefreshToken)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	} else if err != nil {
		log.Println("Error: Failed to get user from database.", err.Error(), "user id:", *profile.Id, "user email:", *profile.Email)
		http.Error(w, "Failed to get user.", http.StatusInternalServerError)
		return
	}

	if user.Provider != AUTH_PROVIDER_GOOGLE {
		http.Error(w, "The user already exists but was created with a different authentication method than Google.", http.StatusBadRequest)
		return
	}

	user, err = state.queries.UsersUpdateNameAndEmail(r.Context(), sqlc.UsersUpdateNameAndEmailParams{ID: *profile.Id, Name: *profile.Name, Email: *profile.Email})
	if err != nil {
		log.Println("Error: Failed to update user data.", err.Error(), "user id:", *profile.Id, "user email:", *profile.Email)
		http.Error(w, "Failed to update user data.", http.StatusInternalServerError)
		return
	}

	tokens := GenAuthTokens(user.ID, user.Email)

	args := sqlc.UsersSetTokensParams{ID: user.ID, AccessToken: utils.SqlNullStrWrapped(tokens.AccessToken), RefreshToken: utils.SqlNullStrWrapped(tokens.RefreshToken)}
	err = state.queries.UsersSetTokens(r.Context(), args)
	if err != nil {
		log.Println("Error: Failed to update user tokens.", err.Error(), "user id:", *profile.Id, "user email:", *profile.Email)
		http.Error(w, "Failed to update user data.", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s?accessToken=%s&refreshToken=%s", clientUrlAuth, tokens.AccessToken, tokens.RefreshToken)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func getUserProfileFromGoogle(ctx context.Context, code string) (*common.GoogleProfile, error) {
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Println("Error: Failed to exchange google auth codes.", err.Error())
		return nil, fmt.Errorf("Error during code exchange. Make sure this request was started from '/auth/google/login'.")
	}

	url := oauthUrlAPIGoogle + token.AccessToken
	response, err := http.Get(url)
	if err != nil {
		log.Println("Error: Failed to get google user info.", err.Error(), "url:", url)
		return nil, fmt.Errorf("Failed to get user info. This might be an issue with the Google Oauth configuration.")
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error: Failed to read google user info.", err.Error())
		return nil, fmt.Errorf("Failed to read user info. Google might have send invalid data.")
	}

	var profile common.GoogleProfile
	err = json.Unmarshal(contents, &profile)
	if err != nil {
		log.Println("Error: Failed to parse google user info.", err.Error(), "contents:", string(contents))
		return nil, fmt.Errorf("Failed to parse user info. Google might have sent invalid data.")
	}

	err = utils.Validate.Struct(profile)
	if err != nil {
		log.Println("Error: Invalid google user info received.", err.Error())
		return nil, fmt.Errorf("Received invalid user info. Google might have sent invalid data.")
	}

	return &profile, nil
}
