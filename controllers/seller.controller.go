package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/harshvardha/e-commerce/cache"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleCreateSeller(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := CreateSellerRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if all the required information is present or not
	if params.GstNumber == "" || params.PanNumber == "" || params.PickupAddress == "" ||
		params.BankAccountHolderName == "" || params.BankAccountNumber == "" || params.IFSCCode == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Incomplete seller information")
		return
	}

	newSeller, err := apiConfig.DB.RegisterSeller(r.Context(), database.RegisterSellerParams{
		ID:                    utility.GenerateSellerID(),
		GstNumber:             params.GstNumber,
		PanNumber:             params.PanNumber,
		PickupAddress:         params.PickupAddress,
		BankAccountHolderName: params.BankAccountHolderName,
		BankAccountNumber:     params.BankAccountNumber,
		IfscCode:              params.IFSCCode,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, CreateSellerResponse{
		ID:                    newSeller.ID,
		GstNumber:             newSeller.GstNumber,
		PanNumber:             newSeller.PanNumber,
		PickupAddress:         newSeller.PickupAddress,
		BankAccountHolderName: newSeller.BankAccountHolderName,
		BankAccountNumber:     newSeller.BankAccountNumber,
		IFSCCode:              newSeller.IfscCode,
		AccessToken:           newAccessToken,
		CreatedAt:             newSeller.CreatedAt,
		UpdatedAt:             newSeller.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleGetSellerContactInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	sellerContactInfo, err := apiConfig.DB.GetSellerContactInfo(r.Context(), IDs.SellerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, SellerContactInfoResponse{
		ID:          sellerContactInfo.ID,
		Email:       sellerContactInfo.Email,
		Phonenumber: sellerContactInfo.PhoneNumber,
		CreatedAt:   sellerContactInfo.CreatedAt,
		UpdatedAt:   sellerContactInfo.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleGetSellerTaxAndAddressInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	sellerTaxAndAddressInfo, err := apiConfig.DB.GetSellerTaxAndAddressInfo(r.Context(), IDs.SellerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, SellerTaxAndAddressInfoResponse{
		ID:            sellerTaxAndAddressInfo.ID,
		GstNumber:     sellerTaxAndAddressInfo.GstNumber,
		PanNumber:     sellerTaxAndAddressInfo.PanNumber,
		PickupAddress: sellerTaxAndAddressInfo.PickupAddress,
		AccessToken:   newAccessToken,
		CreatedAt:     sellerTaxAndAddressInfo.CreatedAt,
		UpdatedAt:     sellerTaxAndAddressInfo.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleGetSellerBankDetails(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	sellerBankDetails, err := apiConfig.DB.GetSellerBankDetails(r.Context(), IDs.SellerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, SellerBankDetailsResponse{
		ID:                    sellerBankDetails.ID,
		BankAccountHolderName: sellerBankDetails.BankAccountHolderName,
		BankAccountNumber:     sellerBankDetails.BankAccountNumber,
		IFSCCode:              sellerBankDetails.IfscCode,
		AccessToken:           newAccessToken,
		CreatedAt:             sellerBankDetails.CreatedAt,
		UpdatedAt:             sellerBankDetails.UpdatedAt,
	})
}

func (apiTwilioConfig *ApiTwilioConfig) HandleUpdateSellerTaxAndAddressInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the otp verification token and otp
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

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateSellerTaxAndAddress{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	existingSellerTaxAndAddressInfo, err := apiTwilioConfig.ApiCfg.DB.GetSellerTaxAndAddressInfo(r.Context(), IDs.SellerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// updating tax and address information
	newTaxAndAddressInfo := database.UpdateSellerTaxAndAddressParams{
		GstNumber:     existingSellerTaxAndAddressInfo.GstNumber,
		PanNumber:     existingSellerTaxAndAddressInfo.PanNumber,
		PickupAddress: existingSellerTaxAndAddressInfo.PickupAddress,
		ID:            IDs.SellerID,
	}

	if params.GstNumber != "" {
		newTaxAndAddressInfo.GstNumber = params.GstNumber
	}
	if params.PanNumber != "" {
		newTaxAndAddressInfo.PanNumber = params.PanNumber
	}
	if params.PickupAddress != "" {
		newTaxAndAddressInfo.PickupAddress = params.PickupAddress
	}

	updatedSellerTaxAndAddressInfo, err := apiTwilioConfig.ApiCfg.DB.UpdateSellerTaxAndAddress(r.Context(), newTaxAndAddressInfo)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, SellerTaxAndAddressInfoResponse{
		ID:            updatedSellerTaxAndAddressInfo.ID,
		GstNumber:     updatedSellerTaxAndAddressInfo.GstNumber,
		PanNumber:     updatedSellerTaxAndAddressInfo.PanNumber,
		PickupAddress: updatedSellerTaxAndAddressInfo.PickupAddress,
		AccessToken:   newAccessToken,
		CreatedAt:     updatedSellerTaxAndAddressInfo.CreatedAt,
		UpdatedAt:     existingSellerTaxAndAddressInfo.UpdatedAt,
	})
}

func (apiTwilioConfig *ApiTwilioConfig) HandleUpdateSellerBankDetails(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the otp verification token and otp
	OTPVerificationToken := r.PathValue("otp_verification_token")
	OTP := r.PathValue("otp")
	if OTPVerificationToken == "" || OTP == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid OTP")
		return
	}

	// verifying the otp
	isOTPCorrect, err := cache.VerifyLoginOTP(apiTwilioConfig.TwilioCfg.Client, apiTwilioConfig.TwilioCfg.VERIFY_SERVICE_SID, OTPVerificationToken, OTP)
	if !isOTPCorrect {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateSellerBankDetails{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	existingSellerBankDetails, err := apiTwilioConfig.ApiCfg.DB.GetSellerBankDetails(r.Context(), IDs.SellerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// updating the bank details
	newBankDetails := database.UpdateSellerBankDetailsParams{
		BankAccountHolderName: existingSellerBankDetails.BankAccountHolderName,
		BankAccountNumber:     existingSellerBankDetails.BankAccountNumber,
		IfscCode:              existingSellerBankDetails.IfscCode,
		ID:                    IDs.SellerID,
	}

	if params.BankAccountHolderName != "" {
		newBankDetails.BankAccountHolderName = params.BankAccountHolderName
	}
	if params.BankAccountNumber != "" {
		newBankDetails.BankAccountNumber = params.BankAccountNumber
	}
	if params.IFSCCode != "" {
		newBankDetails.IfscCode = params.IFSCCode
	}

	updatedBankDetails, err := apiTwilioConfig.ApiCfg.DB.UpdateSellerBankDetails(r.Context(), newBankDetails)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, SellerBankDetailsResponse{
		ID:                    updatedBankDetails.ID,
		BankAccountHolderName: updatedBankDetails.BankAccountHolderName,
		BankAccountNumber:     updatedBankDetails.BankAccountNumber,
		IFSCCode:              updatedBankDetails.IfscCode,
		AccessToken:           newAccessToken,
		CreatedAt:             updatedBankDetails.CreatedAt,
		UpdatedAt:             updatedBankDetails.UpdatedAt,
	})
}

func (apiTwilioConfig *ApiTwilioConfig) HandleRemoveSeller(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the otp verification token and otp
	OTPVerificationToken := r.PathValue("otp_verification_token")
	OTP := r.PathValue("OTP")
	if OTPVerificationToken == "" || OTP == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid OTP")
		return
	}

	// verifying the otp
	isOTPCorrect, err := cache.VerifyLoginOTP(apiTwilioConfig.TwilioCfg.Client, apiTwilioConfig.TwilioCfg.VERIFY_SERVICE_SID, OTPVerificationToken, OTP)
	if !isOTPCorrect {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = apiTwilioConfig.ApiCfg.DB.DeleteSellerAccount(r.Context(), IDs.SellerID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}
