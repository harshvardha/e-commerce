package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/harshvardha/e-commerce/cache"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
	"golang.org/x/crypto/bcrypt"
)

func (apiConfig *ApiConfig) HandleGetAdminInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	admin, err := apiConfig.DB.GetAdminInformation(r.Context(), IDs.AdminID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, AdminResponse{
		Name:        admin.Name,
		Email:       admin.Email,
		Phonenumber: admin.Phonenumber,
		AccessToken: newAccessToken,
		CreatedAt:   admin.CreatedAt,
		UpdatedAt:   admin.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateAdminInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateAdmin{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// updating admin information
	existingAdminInfo, err := apiConfig.DB.GetAdminInformation(r.Context(), IDs.AdminID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newAdminInfo := database.UpdateAdminInformationParams{
		Name:  existingAdminInfo.Name,
		Email: existingAdminInfo.Email,
	}

	updatedAdminInfo, err := apiConfig.DB.UpdateAdminInformation(r.Context(), newAdminInfo)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, AdminResponse{
		Name:        updatedAdminInfo.Name,
		Email:       updatedAdminInfo.Email,
		Phonenumber: updatedAdminInfo.Phonenumber,
		AccessToken: newAccessToken,
		CreatedAt:   updatedAdminInfo.CreatedAt,
		UpdatedAt:   existingAdminInfo.UpdatedAt,
	})
}

func (apiTwilioConfig *ApiTwilioConfig) HandleUpdateAdminPasswordOrPhonenumber(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateAdminPasswordOrPhonenumber{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// if Phonenumber is provided in request body then update only phonenumber
	// if Password is provided in request body then update only password
	// first verify the otp
	// if otp is verified then update the respective information
	if params.OTPVerificationToken != "" && params.OTP != "" {
		isOTPCorrect, err := cache.VerifyLoginOTP(apiTwilioConfig.TwilioCfg.Client, apiTwilioConfig.TwilioCfg.VERIFY_SERVICE_SID, params.OTPVerificationToken, params.OTP)
		if !isOTPCorrect {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if params.Password != "" {
			if !apiTwilioConfig.TwilioCfg.DataValidator.ValidatePassword(params.Password) {
				utility.RespondWithError(w, http.StatusBadRequest, "weak password")
				return
			}
			newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
			if err != nil {
				utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			err = apiTwilioConfig.ApiCfg.DB.UpdateAdminPassword(r.Context(), database.UpdateAdminPasswordParams{
				Password: string(newHashedPassword),
				ID:       IDs.AdminID,
			})
			if err != nil {
				utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
				AccessToken: newAccessToken,
			})
			return
		} else if params.Phonenumber != "" {
			if !utility.ValidatePhonenumber(params.Phonenumber) {
				utility.RespondWithError(w, http.StatusBadRequest, "Invalid Phonenumber")
				return
			}
			udpatedPhonenumber, err := apiTwilioConfig.ApiCfg.DB.UpdateAdminPhonenumber(r.Context(), database.UpdateAdminPhonenumberParams{
				Phonenumber: params.Phonenumber,
				ID:          IDs.AdminID,
			})
			if err != nil {
				utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

			utility.RespondWithJson(w, http.StatusOK, UpdateAdminPhonenumberResponse{
				Phonenumber: udpatedPhonenumber.Phonenumber,
				AccessToken: newAccessToken,
				UpdatedAt:   udpatedPhonenumber.UpdatedAt,
			})
			return
		}
	}

	utility.RespondWithError(w, http.StatusBadRequest, "Invalid otp")
}

func (apiTwilioConfig *ApiTwilioConfig) HandleRemoveAdmin(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching otp verification token and otp from url params
	OTPVerificationToken := r.PathValue("otp_verification_token")
	OTP := r.PathValue("otp")
	if OTPVerificationToken == "" || OTP == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid OTP")
		return
	}

	// verifying otp
	isOTPCorrect, err := cache.VerifyLoginOTP(apiTwilioConfig.TwilioCfg.Client, apiTwilioConfig.TwilioCfg.VERIFY_SERVICE_SID, OTPVerificationToken, OTP)
	if !isOTPCorrect {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// removing admin account
	err = apiTwilioConfig.ApiCfg.DB.RemoveAdmin(r.Context(), IDs.AdminID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}
