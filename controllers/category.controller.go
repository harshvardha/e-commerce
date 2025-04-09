package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleCreateCategory(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CreateOrUpdateCategoryRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// creating new category
	newCategory, err := apiConfig.DB.CreateCategory(r.Context(), database.CreateCategoryParams{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, CreateCategoryResponse{
		ID:          newCategory.ID,
		Name:        newCategory.Name,
		Description: newCategory.Description,
		AccessToken: newAccessToken,
		CreatedAt:   newCategory.CreatedAt,
		UpdatedAt:   newCategory.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateCategory(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the category id
	categoryIDString := r.PathValue("category_id")
	if categoryIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category id")
		return
	}
	categoryID, err := uuid.Parse(categoryIDString)

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := CreateOrUpdateCategoryRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking which fields need to be updated
	existingCategory, err := apiConfig.DB.GetCateogryInformation(r.Context(), categoryID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateCategory := database.UpdateCategoryParams{
		Name:        existingCategory.Name,
		Description: existingCategory.Description,
		ID:          categoryID,
	}

	if params.Name != "" {
		updateCategory.Name = params.Name
	}
	if params.Description != "" {
		updateCategory.Description = params.Description
	}

	// updating the category
	updatedCategory, err := apiConfig.DB.UpdateCategory(r.Context(), updateCategory)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CreateCategoryResponse{
		ID:          categoryID,
		Name:        updatedCategory.Name,
		Description: updatedCategory.Description,
		AccessToken: newAccessToken,
		CreatedAt:   updatedCategory.CreatedAt,
		UpdatedAt:   existingCategory.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleRemoveCategory(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the category id
	categoryIDString := r.PathValue("category_id")
	if categoryIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category id")
		return
	}
	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// removing the category
	err = apiConfig.DB.RemoveCategory(r.Context(), categoryID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetCategoryInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching category id
	categoryIDString := r.PathValue("category_id")
	if categoryIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid category id")
		return
	}
	categoryID, err := uuid.Parse(categoryIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// getting the category information
	category, err := apiConfig.DB.GetCateogryInformation(r.Context(), categoryID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CreateCategoryResponse{
		ID:          categoryID,
		Name:        category.Name,
		Description: category.Description,
		AccessToken: newAccessToken,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	})
}
