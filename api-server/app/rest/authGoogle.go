package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	h.HandleFunc("/auth/google/login", oauthLoginGoogle)
	h.HandleFunc("/auth/google/callback", oauthCallbackGoogle)
}

func oauthLoginGoogle(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie()
	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func oauthCallbackGoogle(w http.ResponseWriter, r *http.Request) {
	profile, err := getUserProfileFromGoogle(r.FormValue("code"))
	if err != nil {
		// TODO
		log.Fatalln(err)
		return
	}

	if !*profile.VerifiedEmail {
		// TODO
		return
	}

	user, err := state.queries.UsersGetById(context.Background(), *profile.Id)
	if err == sql.ErrNoRows {
		user, err = state.queries.UsersCreate(context.Background(), sqlc.UsersCreateParams{ID: *profile.Id, Name: *profile.Name, Email: *profile.Email, Provider: "google"})
		// TODO:
	} else if err != nil {
		// TODO
		log.Fatalln(err)
		return
	}

	if user.Provider != AUTH_PROVIDER_GOOGLE {
		// TODO
		return
	}

	state.queries.UsersUpdateNameAndEmail(context.Background(), sqlc.UsersUpdateNameAndEmailParams{ID: *profile.Id, Name: *profile.Name, Email: *profile.Email})
	// TODO
}

// TODO: simulator
func getUserProfileFromGoogle(code string) (*googleProfile, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Error during code exchange: %s", err.Error())
	}

	response, err := simulator.S.HttpGet(oauthUrlAPIGoogle + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read user info: %s", err.Error())
	}

	// TODO: check that all required fields are set
	var profile googleProfile
	err = json.Unmarshal(contents, &profile)
	if err != nil {
		return nil, err
	}

	err = utils.Validate.Struct(profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
