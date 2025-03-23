package controllers

import (
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/twilio/twilio-go"
)

type ApiConfig struct {
	DB        *database.Queries
	JwtSecret string
}

type TwilioConfig struct {
	VERIFY_SERVICE_SID string
	Client             *twilio.RestClient
}

type ApiTwilioConfig struct {
	ApiCfg    ApiConfig
	TwilioCfg TwilioConfig
}

type SendOtpRequest struct {
	Phonenumber string `json:"phone_number"`
}

type VerifyOtpRequest struct {
	VerificationToken string `json:"verification_token"`
	OTP               string `json:"otp"`
}

type SendOTPResponse struct {
	VerificationToken string `json:"verification_token"`
}
