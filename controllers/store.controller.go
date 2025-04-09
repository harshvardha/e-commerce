package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/harshvardha/e-commerce/cache"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleCreateStore(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := CreateOrUpdateStoreRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the name param is empty or not
	if params.Name == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Incomplete store information")
		return
	}

	newStore, err := apiConfig.DB.CreateStore(r.Context(), database.CreateStoreParams{
		Name:    params.Name,
		OwnerID: IDs.SellerID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, CreateStoreResponse{
		ID:          newStore.ID,
		Name:        newStore.Name,
		SellerID:    newStore.OwnerID,
		AccessToken: newAccessToken,
		CreatedAt:   newStore.CreatedAt,
		UpdatedAt:   newStore.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateStore(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := CreateOrUpdateStoreRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if name params is empty or not
	if params.Name == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Incomplete store information")
		return
	}

	updatedStore, err := apiConfig.DB.UpdateStoreInformation(r.Context(), database.UpdateStoreInformationParams{
		Name: params.Name,
		ID:   IDs.StoreID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, UpdateStoreResponse{
		Name:        updatedStore.Name,
		AccessToken: newAccessToken,
		UpdatedAt:   updatedStore.UpdatedAt,
	})
}

func (apiTwilioConfig *ApiTwilioConfig) HandleDeleteStore(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching otp verification token and otp
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

	err = apiTwilioConfig.ApiCfg.DB.DeleteStore(r.Context(), IDs.StoreID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetStore(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	store, err := apiConfig.DB.GetStoreInformation(r.Context(), IDs.StoreID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, GetStoreInformationResponse{
		Name:        store.Name,
		AccessToken: newAccessToken,
		CreatedAt:   store.CreatedAt,
		UpdatedAt:   store.UpdatedAt,
	})
}
