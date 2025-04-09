package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
	"golang.org/x/crypto/bcrypt"
)

func (apiConfig *ApiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	user, err := apiConfig.DB.GetUserByID(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, "User not found.")
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ResponseUser{
		Email:       user.Email,
		Phonenumber: user.PhoneNumber,
		AccessToken: newAccessToken,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateUserInfo{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking what needs to be updated (email || phone_number || password)
	existingUser, err := apiConfig.DB.GetUserByID(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, "User does not exist")
		return
	}

	var updatedUserInfo database.UpdateUserParams
	if params.Email != "" {
		updatedUserInfo.Email = params.Email
	} else {
		updatedUserInfo.Email = existingUser.Email
	}

	if params.Phonenumber != "" {
		updatedUserInfo.PhoneNumber = params.Phonenumber
	} else {
		updatedUserInfo.PhoneNumber = existingUser.PhoneNumber
	}

	if params.Password != "" {
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		updatedUserInfo.Password = string(newHashedPassword)
	} else {
		updatedUserInfo.Password = existingUser.Password
	}

	updatedUser, err := apiConfig.DB.UpdateUser(r.Context(), updatedUserInfo)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ResponseUser{
		Email:       updatedUser.Email,
		Phonenumber: updatedUser.PhoneNumber,
		AccessToken: newAccessToken,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleDeleteUser(w http.ResponseWriter, r *http.Request, IDs ID, newAccessToken string) {
	deletedUser, err := apiConfig.DB.DeleteUser(r.Context(), IDs.UserID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ResponseUser{
		Email:       deletedUser.Email,
		Phonenumber: deletedUser.PhoneNumber,
		AccessToken: newAccessToken,
		CreatedAt:   deletedUser.CreatedAt,
		UpdatedAt:   deletedUser.UpdatedAt,
	})
}
