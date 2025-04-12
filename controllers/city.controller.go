package controllers

// add a admin id validation middleware for this controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleAddCity(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := AddOrUpdateCityRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// adding new city
	newCity, err := apiConfig.DB.CreateCity(r.Context(), database.CreateCityParams{
		Name:  params.Name,
		State: params.State,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	state, err := apiConfig.DB.GetStateName(r.Context(), params.State)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, CityResponse{
		ID:          newCity.ID,
		Name:        newCity.Name,
		State:       state,
		AccessToken: newAccessToken,
		CreatedAt:   newCity.CreatedAt,
		UpdatedAt:   newCity.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateCity(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching city id
	cityIDString := r.PathValue("city_id")
	if cityIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid City ID")
		return
	}
	cityID, err := uuid.Parse(cityIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// decoding request params
	decoder := json.NewDecoder(r.Body)
	params := AddOrUpdateCityRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking which fields to update
	existingCityInformation, err := apiConfig.DB.GetCityById(r.Context(), cityID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateCity := database.UpdateCityParams{
		Name:  existingCityInformation.Name,
		State: existingCityInformation.State,
		ID:    cityID,
	}

	if params.Name != "" {
		updateCity.Name = params.Name
	}

	if params.State != uuid.Nil {
		updateCity.State = params.State
	}

	// updating city
	updatedCity, err := apiConfig.DB.UpdateCity(r.Context(), updateCity)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	state, err := apiConfig.DB.GetStateName(r.Context(), updateCity.State)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CityResponse{
		ID:          updatedCity.ID,
		Name:        updatedCity.Name,
		State:       state,
		AccessToken: newAccessToken,
		CreatedAt:   updatedCity.CreatedAt,
		UpdatedAt:   updatedCity.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleRemoveCity(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching city id
	cityIDString := r.PathValue("city_id")
	if cityIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid City ID")
		return
	}
	cityID, err := uuid.Parse(cityIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// removing city
	err = apiConfig.DB.RemoveCity(r.Context(), cityID)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetCityInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching city id
	cityIDString := r.PathValue("city_id")
	if cityIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid City ID")
		return
	}
	cityID, err := uuid.Parse(cityIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// getting city information
	cityInformation, err := apiConfig.DB.GetCityById(r.Context(), cityID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	state, err := apiConfig.DB.GetStateName(r.Context(), cityInformation.State)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, CityResponse{
		ID:          cityID,
		Name:        cityInformation.Name,
		State:       state,
		AccessToken: newAccessToken,
		CreatedAt:   cityInformation.CreatedAt,
		UpdatedAt:   cityInformation.UpdatedAt,
	})
}
