package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleCreateCustomer(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// creating a buyer
	err := apiConfig.DB.CreateCustomer(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleUpdateCustomerAddress(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	//	decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateCustomerAddress{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking which fields to update
	existingCustomerAddress, err := apiConfig.DB.GetCustomerAddress(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateCustomerAddress := database.UpdateCustomerAddressParams{
		DeliveryAddress: existingCustomerAddress.DeliveryAddress,
		Pincode:         existingCustomerAddress.Pincode,
		City:            existingCustomerAddress.City,
		State:           existingCustomerAddress.State,
	}

	if params.DeliveryAddress != "" {
		updateCustomerAddress.DeliveryAddress = params.DeliveryAddress
	}
	if params.Pincode != "" {
		if len(params.Pincode) != 6 {
			utility.RespondWithError(w, http.StatusBadRequest, "Invalid Pincode")
			return
		}
		updateCustomerAddress.Pincode = params.Pincode
	}
	if params.City != uuid.Nil {
		updateCustomerAddress.City = params.City
	}
	if params.State != uuid.Nil {
		updateCustomerAddress.State = params.State
	}

	// updating customer address
	updatedCustomerAddress, err := apiConfig.DB.UpdateCustomerAddress(r.Context(), updateCustomerAddress)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cityAndState, err := apiConfig.DB.GetCityAndState(r.Context(), updatedCustomerAddress.City)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CustomerResponse{
		ID:              updatedCustomerAddress.ID,
		DeliveryAddress: updateCustomerAddress.DeliveryAddress,
		Pincode:         updatedCustomerAddress.Pincode,
		City:            cityAndState.CityName,
		State:           cityAndState.StateName,
		AccessToken:     newAccessToken,
		CreatedAt:       updatedCustomerAddress.CreatedAt,
		UpdatedAt:       updatedCustomerAddress.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleGetCustomerAddress(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	customerAddress, err := apiConfig.DB.GetCustomerAddress(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cityAndState, err := apiConfig.DB.GetCityAndState(r.Context(), customerAddress.City)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CustomerResponse{
		ID:              IDs.UserID,
		DeliveryAddress: customerAddress.DeliveryAddress,
		Pincode:         customerAddress.Pincode,
		City:            cityAndState.CityName,
		State:           cityAndState.StateName,
		AccessToken:     newAccessToken,
	})
}
