package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/simulator"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleProfile struct {
	Id            *string `json:"id" validate:"required"`
	Email         *string `json:"email" validate:"required"`
	VerifiedEmail *bool   `json:"verified_email" validate:"required"`
	Name          *string `json:"name" validate:"required"`
}

// Google config:
// https://console.cloud.google.com/auth/clients?highlightClient=750385423567-rsrv4dknuvrts9rv5neab3dl667r5la6.apps.googleusercontent.com&authuser=2&organizationId=0&project=htl-ride-sharing

// TODO: change secrets and put them in ENV

const oauthUrlAPIGoogle = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://127.0.0.1:8000/auth/google/callback",
	ClientID:     "750385423567-rsrv4dknuvrts9rv5neab3dl667r5la6.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func authHandlersGoogle(h simulator.HTTPMux) {
	h.HandleFunc("GET /auth/google/login", oauthLoginGoogle)
	h.HandleFunc("GET /auth/google/callback", oauthCallbackGoogle)
}

func oauthLoginGoogle(w http.ResponseWriter, r *http.Request) {
	oauthState := genRandBase64(16)
	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthCallbackGoogle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	profile, err := getUserProfileFromGoogle(r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !*profile.VerifiedEmail {
		http.Error(w, "Google user has an unverified email address. This is not allowed.", http.StatusBadRequest)
		return
	}

	user, err := state.queries.UsersGetById(context.Background(), *profile.Id)
	if err == sql.ErrNoRows {
		user, err = state.queries.UsersCreate(context.Background(), sqlc.UsersCreateParams{ID: *profile.Id, Name: *profile.Name, Email: *profile.Email, Provider: "google"})
		if err != nil {
			log.Println("Error: Failed to create user after Google authentication.", err.Error(), "user email:", *profile.Email)
			http.Error(w, "Failed to create user after Google authentication.", http.StatusInternalServerError)
			return
		}

		tokens := genAuthTokens(user.ID, user.Email)

		args := sqlc.UsersSetTokensParams{ID: user.ID, AccessToken: sql.NullString{String: tokens.AccessToken}, RefreshToken: sql.NullString{String: tokens.RefreshToken}}
		err = state.queries.UsersSetTokens(context.Background(), args)
		if err != nil {
			log.Println("Error: Failed to update user tokens.", err.Error(), "user id:", *profile.Id, "user email:", *profile.Email)
			http.Error(w, "Failed to update user data.", http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(tokens)
		assert.True(err == nil, "Failed to serialize authentication response.", tokens)

		w.Write(bytes)

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

	user, err = state.queries.UsersUpdateNameAndEmail(context.Background(), sqlc.UsersUpdateNameAndEmailParams{ID: *profile.Id, Name: *profile.Name, Email: *profile.Email})
	if err != nil {
		log.Println("Error: Failed to update user data.", err.Error(), "user id:", *profile.Id, "user email:", *profile.Email)
		http.Error(w, "Failed to update user data.", http.StatusInternalServerError)
		return
	}

	tokens := genAuthTokens(user.ID, user.Email)

	args := sqlc.UsersSetTokensParams{ID: user.ID, AccessToken: utils.SqlNullStr(tokens.AccessToken), RefreshToken: utils.SqlNullStr(tokens.RefreshToken)}
	err = state.queries.UsersSetTokens(context.Background(), args)
	if err != nil {
		log.Println("Error: Failed to update user tokens.", err.Error(), "user id:", *profile.Id, "user email:", *profile.Email)
		http.Error(w, "Failed to update user data.", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(tokens)
	assert.True(err == nil, "Failed to serialize authentication response.", tokens)
	w.Write(bytes)
}

func getUserProfileFromGoogle(code string) (*googleProfile, error) {
	// TODO: validate state
	token, err := googleOauthConfig.Exchange(context.Background(), code) // TODO: simulator
	if err != nil {
		log.Println("Error: Failed to exchange google auth codes.", err.Error())
		return nil, fmt.Errorf("Error during code exchange. Make sure this request was started from '/auth/google/login'.")
	}

	url := oauthUrlAPIGoogle + token.AccessToken
	response, err := simulator.S.HttpGet(url)
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

	var profile googleProfile
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
