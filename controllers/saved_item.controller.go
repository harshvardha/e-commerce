package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleAddProductToSavedItems(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the product id
	productIDString := r.PathValue("product_id")
	if productIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	productID, err := uuid.Parse(productIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = apiConfig.DB.AddProductToSavedItems(r.Context(), database.AddProductToSavedItemsParams{
		ProductID: productID,
		UserID:    IDs.UserID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleRemoveProductFromSavedItems(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching product id
	productIDString := r.PathValue("product_id")
	if productIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	productID, err := uuid.Parse(productIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = apiConfig.DB.RemoveProductFromSavedItems(r.Context(), database.RemoveProductFromSavedItemsParams{
		UserID:    IDs.UserID,
		ProductID: productID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetAllSavedItems(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	savedItems, err := apiConfig.DB.GetAllSavedItems(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, struct {
		SavedItems  []database.GetAllSavedItemsRow `json:"saved_items"`
		AccessToken string                         `json:"access_token"`
	}{
		SavedItems:  savedItems,
		AccessToken: newAccessToken,
	})
}
