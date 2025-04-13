package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleAddProductToCart(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
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

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := AddProductToCartRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if params.Quantity == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Product Quantity Cannot be 0")
		return
	}

	// adding product to cart
	err = apiConfig.DB.AddProductToCart(r.Context(), database.AddProductToCartParams{
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

func (apiConfig *ApiConfig) HandleUpdateProductQuantity(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
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

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := AddProductToCartRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if params.Quantity == 0 {
		utility.RespondWithError(w, http.StatusBadRequest, "Product Quantity Cannot be 0")
		return
	}

	// updating product quantity
	err = apiConfig.DB.UpdateProductQuantity(r.Context(), database.UpdateProductQuantityParams{
		Quantity:  params.Quantity,
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

func (apiConfig *ApiConfig) HandleRemoveProductFromCart(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
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

	// removing product from cart
	err = apiConfig.DB.RemoveProductFromCart(r.Context(), database.RemoveProductFromCartParams{
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

func (apiConfig *ApiConfig) HandleGetAllProductsInCart(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// getting all the products in user cart
	cart, err := apiConfig.DB.GetAllProductsInCart(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, struct {
		Cart        []database.GetAllProductsInCartRow `json:"cart"`
		AccessToken string                             `json:"access_token"`
	}{
		Cart:        cart,
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleEmptyCart(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	err := apiConfig.DB.EmptyCart(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}
