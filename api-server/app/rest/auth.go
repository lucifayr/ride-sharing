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
	Id     *string `json:"id" validate:"required"`
	Email  *string `json:"email" validate:"required"`
	Random *string `json:"random" validate:"required"`
}

// TODO: change and move into ENV

// Must be 32 bytes long!
const authTokenEncodingSecretKey = "268aTvg3uNE*xLkB7tYSW%Cl#CmuY5!L"
const authStateEncodingSecretKey = "@j&P4m$Fcq$en*C75six#9dbNBDijJgU"

func genAuthTokens(userId string, email string) authTokens {
	return authTokens{
		AccessToken:  encodeAccessToken(userId, email),
		RefreshToken: genRandBase64(512),
	}
}

func encodeAccessToken(userId string, email string) string {
	token := genRandBase64(512)
	at := accessToken{
		Id:     &userId,
		Email:  &email,
		Random: &token,
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
