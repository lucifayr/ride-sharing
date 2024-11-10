package common

type GoogleProfile struct {
	Id            *string `json:"id" validate:"required"`
	Email         *string `json:"email" validate:"required"`
	VerifiedEmail *bool   `json:"verified_email" validate:"required"`
	Name          *string `json:"name" validate:"required"`
}
