package controllers

// add a admin id validation middleware for this controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

func (apiConfig *ApiConfig) HandleAddState(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := AddOrUpdateStateRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// adding a new state
	newState, err := apiConfig.DB.CreateState(r.Context(), params.Name)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, StateResponse{
		ID:          newState.ID,
		Name:        newState.Name,
		AccessToken: newAccessToken,
		CreatedAt:   newState.CreatedAt,
		UpdatedAt:   newState.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateState(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the state id
	stateIDString := r.PathValue("state_id")
	if stateIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid state id")
		return
	}
	stateID, err := uuid.Parse(stateIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := AddOrUpdateStateRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking which fields need to be updated
	existingStateInformation, err := apiConfig.DB.GetStateName(r.Context(), stateID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updateState := database.UpdateStateParams{
		Name: existingStateInformation,
		ID:   stateID,
	}

	if params.Name != "" {
		updateState.Name = params.Name
	}

	// updating state information
	updatedState, err := apiConfig.DB.UpdateState(r.Context(), updateState)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, StateResponse{
		ID:          stateID,
		Name:        updatedState.Name,
		AccessToken: newAccessToken,
		CreatedAt:   updatedState.CreatedAt,
		UpdatedAt:   updatedState.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleRemoveState(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the state id
	stateIDString := r.PathValue("state_id")
	if stateIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid State ID")
		return
	}
	stateID, err := uuid.Parse(stateIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// removing state
	err = apiConfig.DB.RemoveState(r.Context(), stateID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) GetStateInformation(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// fetching the state id
	stateIDString := r.PathValue("state_id")
	if stateIDString == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid State ID")
		return
	}
	stateID, err := uuid.Parse(stateIDString)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// getting the state information
	state, err := apiConfig.DB.GetStateName(r.Context(), stateID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, StateResponse{
		ID:          stateID,
		Name:        state,
		AccessToken: newAccessToken,
	})
}
