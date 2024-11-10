package rest

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/simulator"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"time"
)

const (
	AUTH_PROVIDER_GOOGLE    = "google"
	AUTH_PROVIDER_MICROSOFT = "microsoft"
)

const clientUrlAuth = "http://127.0.0.1:5173/authenticate"

type authTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type accessToken struct {
	Id        *string   `json:"id" validate:"required"`
	Email     *string   `json:"email" validate:"required"`
	Random    *string   `json:"random" validate:"required"`
	ExpiresAt time.Time `json:"expiresAt" validate:"required"`
}

// TODO: change and move into ENV

// Must be 32 bytes long!
const authTokenEncodingSecretKey = "268aTvg3uNE*xLkB7tYSW%Cl#CmuY5!L"
const authStateEncodingSecretKey = "@j&P4m$Fcq$en*C75six#9dbNBDijJgU"

func authHandlers(h simulator.HTTPMux) {
	h.HandleFunc("POST /auth/refresh", handle(refreshAuthTokens).with(bearerAuth(true)).build())
}

type refreshAuthTokensBody struct {
	RefreshToken *string `json:"refreshToken" validate:"required"`
}

func refreshAuthTokens(w http.ResponseWriter, r *http.Request) {
	// Locking to prevent change of access token during refresh
	state.mutex.Lock()
	defer state.mutex.Unlock()

	user := getMiddlewareData[sqlc.User](r, "user")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusBadRequest)
		return
	}

	var refresh refreshAuthTokensBody
	err = json.Unmarshal(body, &refresh)
	if err != nil {
		http.Error(w, "Failed to read request body as JSON.", http.StatusBadRequest)
		return
	}

	err = utils.Validate.Struct(refresh)
	if err != nil {
		http.Error(w, "Invalid JSON in request body. "+err.Error(), http.StatusBadRequest)
		return
	}

	if !user.RefreshToken.Valid || user.RefreshToken.String != *refresh.RefreshToken {
		http.Error(w, "Invalid refresh token cannot be used to get new tokens.", http.StatusUnauthorized)
		return
	}

	tokens := genAuthTokens(user.ID, user.Email)
	bytes, err := json.Marshal(tokens)
	assert.True(err == nil, "Failed to serialize authentication tokens.", tokens, "error:", func() any { return err })

	args := sqlc.UsersSetTokensParams{ID: user.ID, AccessToken: utils.SqlNullStr(tokens.AccessToken), RefreshToken: utils.SqlNullStr(tokens.RefreshToken)}
	err = state.queries.UsersSetTokens(r.Context(), args)
	if err != nil {
		log.Println("Failed to update authentication tokens.", "error:", err)
		http.Error(w, "Failed to update authentication tokens.", http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func genAuthTokens(userId string, email string) authTokens {
	return authTokens{
		AccessToken:  encodeAccessToken(userId, email),
		RefreshToken: genRandBase64(512),
	}
}

func encodeAccessToken(userId string, email string) string {
	token := genRandBase64(512)
	at := accessToken{
		Id:        &userId,
		Email:     &email,
		Random:    &token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	plain, err := json.Marshal(at)
	assert.True(err == nil, "Invalid token JSON.", "access-token:", at, "error", func() any { return err.Error() })

	ciphertext, err := encrypt(plain, []byte(authTokenEncodingSecretKey))
	assert.True(err == nil, "Encryption error on server defined data.", "error:", func() any { return err.Error() })

	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decodeAccessToken(token []byte) (*accessToken, error) {
	token, err := base64.URLEncoding.DecodeString(string(token))
	if err != nil {
		return nil, err
	}

	plain, err := decrypt(token, []byte(authTokenEncodingSecretKey))
	if err != nil {
		return nil, err
	}

	var at accessToken
	err = json.Unmarshal(plain, &at)
	if err != nil {
		return nil, err
	}

	err = utils.Validate.Struct(at)
	if err != nil {
		return nil, err
	}

	return &at, nil
}

func encrypt(plain []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	if len(plain) < aes.BlockSize {
		return []byte{}, fmt.Errorf("Plaintext length too short.")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return []byte{}, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plain)

	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	if len(ciphertext) < aes.BlockSize {
		return []byte{}, errors.New("Ciphertext length too short.")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func genRandBase64(size int) string {
	b := make([]byte, size)
	simulator.S.RandCrypto(b)
	return base64.URLEncoding.EncodeToString(b)
}
