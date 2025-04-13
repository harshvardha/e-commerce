package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleCreateOrder(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding the request body
	type Product struct {
		ID       uuid.UUID `json:"id"`
		Quantity int32     `json:"quantity"`
	}
	decoder := json.NewDecoder(r.Body)
	params := struct {
		TotalValue float64   `json:"total_value"`
		Products   []Product `json:"products"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// getting the order status id for status: 'RECIEVED'
	orderStatusRecievedID, err := apiConfig.DB.GetStatusID(r.Context(), "RECIEVED")
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// creating order
	order, err := apiConfig.DB.CreateOrder(r.Context(), database.CreateOrderParams{
		TotalValue: params.TotalValue,
		SellerID:   IDs.SellerID,
		Status:     orderStatusRecievedID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// adding products to order
	for _, product := range params.Products {
		err = apiConfig.DB.AddProductToOrder(r.Context(), database.AddProductToOrderParams{
			UserID:    IDs.UserID,
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  product.Quantity,
		})
		if err != nil {
			err = apiConfig.DB.DeleteOrder(r.Context(), order.ID)
			if err != nil {
				utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	// getting order details to return as response
	orderDetails, err := apiConfig.DB.GetOrderDetails(r.Context(), database.GetOrderDetailsParams{
		UserID:  IDs.UserID,
		OrderID: order.ID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, OrderDetailsResponse{
		OrderDetails: orderDetails,
		AccessToken:  newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleUpdateOrderStatus(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching order id
	orderIDString := r.PathValue("order_id")
	statusIDString := r.PathValue("status_id")
	if orderIDString == "" || statusIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Order ID Or Status ID")
		return
	}
	orderID, err := uuid.Parse(orderIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	statusID, err := uuid.Parse(statusIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// updating order status
	err = apiConfig.DB.UpdateOrderStatus(r.Context(), database.UpdateOrderStatusParams{
		Status: statusID,
		ID:     orderID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetOrderDetails(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the order id
	orderIDString := r.PathValue("order_id")
	if orderIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Order ID")
		return
	}
	orderID, err := uuid.Parse(orderIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// fetching order details
	orderDetails, err := apiConfig.DB.GetOrderDetails(r.Context(), database.GetOrderDetailsParams{
		UserID:  IDs.UserID,
		OrderID: orderID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, OrderDetailsResponse{
		OrderDetails: orderDetails,
		AccessToken:  newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetAllOrders(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	userOrders, err := apiConfig.DB.GetAllOrders(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, struct {
		Orders      []database.GetAllOrdersRow `json:"orders"`
		AccessToken string                     `json:"access_token"`
	}{
		Orders:      userOrders,
		AccessToken: newAccessToken,
	})
}
