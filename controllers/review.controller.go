package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleAddReview(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
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
	params := CreateOrUpdateReviewRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// adding the review
	newReview, err := apiConfig.DB.CreateReview(r.Context(), database.CreateReviewParams{
		Description: params.Description,
		UserID:      IDs.UserID,
		ProductID:   productID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, ReviewResponse{
		ID:          newReview.ID,
		Description: newReview.Description,
		AccessToken: newAccessToken,
		CreatedAt:   newReview.CreatedAt,
		UpdatedAt:   newReview.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateReview(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the review id
	reviewIDString := r.PathValue("review_id")
	if reviewIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Review ID")
		return
	}
	reviewID, err := uuid.Parse(reviewIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CreateOrUpdateReviewRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking which fields to update
	existingReview, err := apiConfig.DB.GetReviewByID(r.Context(), reviewID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateReview := database.UpdateReviewParams{
		Description: existingReview,
		ID:          reviewID,
		UserID:      IDs.UserID,
	}

	if params.Description != "" {
		updateReview.Description = params.Description
	}

	// updating the review
	updatedReview, err := apiConfig.DB.UpdateReview(r.Context(), updateReview)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ReviewResponse{
		ID:          reviewID,
		Description: updatedReview.Description,
		AccessToken: newAccessToken,
		CreatedAt:   updatedReview.CreatedAt,
		UpdatedAt:   updatedReview.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleRemoveReview(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the review id
	reviewIDString := r.PathValue("review_id")
	if reviewIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Review ID")
		return
	}
	reviewID, err := uuid.Parse(reviewIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// removing the review
	err = apiConfig.DB.RemoveReview(r.Context(), database.RemoveReviewParams{
		ID:     reviewID,
		UserID: IDs.UserID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetReviewsByProductID(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
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

	// getting all the reviews for the product
	productReviews, err := apiConfig.DB.GetReviewsByProductID(r.Context(), productID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, struct {
		Reviews     []database.Review `json:"reviews"`
		AccessToken string            `json:"access_token"`
	}{
		Reviews:     productReviews,
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetReviewsByUserID(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// getting all the reviews for a user
	userReviews, err := apiConfig.DB.GetReviewsByUserID(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, struct {
		Reviews     []database.Review `json:"reviews"`
		AccessToken string            `json:"access_token"`
	}{
		Reviews:     userReviews,
		AccessToken: newAccessToken,
	})
}
