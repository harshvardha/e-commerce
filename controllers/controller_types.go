package controllers

import (
	"time"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
	"github.com/twilio/twilio-go"
)

type ApiConfig struct {
	DB        *database.Queries
	JwtSecret string
}

type TwilioConfig struct {
	VERIFY_SERVICE_SID string
	Client             *twilio.RestClient
	DataValidator      utility.Validator
}

type ApiTwilioConfig struct {
	ApiCfg    ApiConfig
	TwilioCfg TwilioConfig
}

type SendOtpRequest struct {
	Email       string `json:"email"`
	Phonenumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type VerifyOtpRequest struct {
	VerificationToken string `json:"verification_token"`
	OTP               string `json:"otp"`
	RequestType       int    `json:"request_type"`
}

type SendOTPResponse struct {
	VerificationToken string `json:"verification_token"`
}

type ResendOTPRequest struct {
	VerificationToken string `json:"verification_token"`
	RequestType       int    `json:"request_type"`
}

type LoginRequest struct {
	Phonenumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserTokenClaims struct {
	UserID   uuid.UUID
	SellerID string
	StoreID  uuid.UUID
}

type ResponseUser struct {
	Email       string    `json:"email"`
	Phonenumber string    `json:"phone_number"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
