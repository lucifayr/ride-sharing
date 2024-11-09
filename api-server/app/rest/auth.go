package rest

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/simulator"
	"ride_sharing_api/app/utils"
)

const (
	AUTH_PROVIDER_GOOGLE    = "google"
	AUTH_PROVIDER_MICROSOFT = "microsoft"
)

type authTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type accessToken struct {
	Id    *string `json:"id" validate:"required"`
	Email *string `json:"email" validate:"required"`
	Token *string `json:"token" validate:"required"`
}

// TODO: change and move into ENV

// Must be 32 bytes long!
const authTokenEncodingSecretKey = "268aTvg3uNE*xLkB7tYSW%Cl#CmuY5!L"

func genAuthTokens(userId string, email string) authTokens {
	return authTokens{
		AccessToken:  encodeAccessToken(userId, email),
		RefreshToken: genRandBase64(512),
	}
}

func encodeAccessToken(userId string, email string) string {
	token := genRandBase64(512)
	at := accessToken{
		Id:    &userId,
		Email: &email,
		Token: &token,
	}

	bytes, err := json.Marshal(at)
	assert.True(err == nil, "Invalid token JSON.", "access-token:", at)

	plain := []byte(base64.URLEncoding.EncodeToString(bytes))

	block, err := aes.NewCipher([]byte(authTokenEncodingSecretKey))
	assert.True(err == nil, "Invalid aes secret key.")
	assert.True(len(plain) > block.BlockSize(), "Plaintext length too short.")

	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	assert.True(err == nil, "Failed to read block-size from ciphertext")

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plain)

	return string(ciphertext)
}

func decodeAccessToken(token []byte) (*accessToken, error) {
	block, err := aes.NewCipher([]byte(authTokenEncodingSecretKey))
	assert.True(err == nil, "Invalid aes secret key.")

	if len(token) < aes.BlockSize {
		return nil, errors.New("ciphertext too short/invalid.")
	}

	iv := token[:aes.BlockSize]
	token = token[aes.BlockSize:]

	plain := make([]byte, len(token))
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plain, token)

	plain, err = base64.URLEncoding.DecodeString(string(plain))
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

func genRandBase64(size int) string {
	b := make([]byte, size)
	simulator.S.RandCrypto(b)
	return base64.URLEncoding.EncodeToString(b)
}
