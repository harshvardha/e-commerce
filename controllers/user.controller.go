package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (twilioCfg *TwilioConfig) SendOTPForVerification(w http.ResponseWriter, r *http.Request) {
	// decoding request
	decoder := json.NewDecoder(r.Body)
	params := SendOtpRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// if phone number is available then sending otp to sms
	// otherwise sending otp to mail
	if utility.Validate_Phonenumber(params.Phonenumber) {
		verificationToken, err := utility.SendOtpViaPhonenumber(
			params.Phonenumber,
			twilioCfg.Client,
			twilioCfg.VERIFY_SERVICE_SID,
		)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utility.RespondWithJson(w, http.StatusOK, SendOTPResponse{
			VerificationToken: verificationToken,
		})
		return
	}

	utility.RespondWithError(w, http.StatusBadRequest, "Invalid Phonenumber")
}

func (apiTwilioCfg *ApiTwilioConfig) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := VerifyOtpRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// verifying otp
	isOTPCorrect, phonenumber, err := utility.VerifyOtp(
		apiTwilioCfg.TwilioCfg.Client,
		apiTwilioCfg.TwilioCfg.VERIFY_SERVICE_SID,
		params.VerificationToken,
		params.OTP,
	)
	if !isOTPCorrect {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// registering user if otp is verified
	user, err := apiTwilioCfg.ApiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: sql.NullString{
			String: "",
			Valid:  true,
		},
		PhoneNumber: sql.NullString{
			String: phonenumber,
			Valid:  true,
		},
	})

	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, user)
}
