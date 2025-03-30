package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/cache"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
	"golang.org/x/crypto/bcrypt"
)

func (twilioConfig *TwilioConfig) HandleSendOTP(w http.ResponseWriter, r *http.Request) {
	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := SendOtpRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// validating phonenumber param as it will be present in both login and registration request
	if len(params.Phonenumber) > 13 || !utility.ValidatePhonenumber(params.Phonenumber) {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Phonenumber")
		return
	}

	var verificationToken string

	// checking if the request is a registration or login type
	if len(params.Email) > 0 && len(params.Password) > 0 {
		// request is registration type
		// validating email, phonenumber and password
		if !twilioConfig.DataValidator.ValidateEmail(params.Email) {
			utility.RespondWithError(w, http.StatusBadRequest, "Invalid Email")
			return
		}

		if !twilioConfig.DataValidator.ValidatePassword(params.Password) {
			utility.RespondWithError(w, http.StatusBadRequest, "Weak Password")
			return
		}

		// sending otp
		verificationToken, err = cache.SendOtpViaPhonenumber(twilioConfig.Client, twilioConfig.VERIFY_SERVICE_SID, params.Phonenumber, params.Email, params.Password)
	} else {
		// sending otp for login request type
		verificationToken, err = cache.SendOtpViaPhonenumber(twilioConfig.Client, twilioConfig.VERIFY_SERVICE_SID, params.Phonenumber)
	}

	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, SendOTPResponse{
		VerificationToken: verificationToken,
	})
}

func (apiTwilioConfig *ApiTwilioConfig) HandleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := VerifyOtpRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// if requestType == 1 then it is login otp verification request
	if params.RequestType == 1 {
		isCorrectOTP, err := cache.VerifyLoginOTP(
			apiTwilioConfig.TwilioCfg.Client,
			apiTwilioConfig.TwilioCfg.VERIFY_SERVICE_SID,
			params.VerificationToken,
			params.OTP,
		)
		if !isCorrectOTP {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		utility.RespondWithJson(w, http.StatusOK, nil)
	} else if params.RequestType == 3 {
		isCorrectOTP, userDetails, err := cache.VerifyRegistrationOTP(
			apiTwilioConfig.TwilioCfg.Client,
			apiTwilioConfig.TwilioCfg.VERIFY_SERVICE_SID,
			params.VerificationToken,
			params.OTP,
		)
		if !isCorrectOTP {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// checking if the user already exist or not
		userIdExist, err := apiTwilioConfig.ApiCfg.DB.DoesUserExist(r.Context(), userDetails.Phonenumber)
		if err == nil && userIdExist != uuid.Nil {
			utility.RespondWithError(w, http.StatusBadRequest, "User already exist.")
			return
		}

		// registering user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		user, err := apiTwilioConfig.ApiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			Email:       userDetails.Email,
			PhoneNumber: userDetails.Phonenumber,
			Password:    string(hashedPassword),
		})
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utility.RespondWithJson(w, http.StatusCreated, user)
	}
}

func (twilioConfig *TwilioConfig) HandleResendOTP(w http.ResponseWriter, r *http.Request) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := ResendOTPRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Resend OTP request details")
		return
	}

	// resending the otp
	newVerificationToken, err := cache.ResendOTP(twilioConfig.Client, twilioConfig.VERIFY_SERVICE_SID, params.VerificationToken, params.RequestType)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	utility.RespondWithJson(w, http.StatusOK, SendOTPResponse{
		VerificationToken: newVerificationToken,
	})
}

func (apiConfig *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := LoginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid login credentials")
		return
	}

	// checking if the user exist or not
	userExist, err := apiConfig.DB.GetUser(r.Context(), params.Phonenumber)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, "User does not exist")
		return
	}

	// checking if the password is correct or not
	err = bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(params.Password))
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Password")
		return
	}

	// creating token claims for access token

	// checking if user is seller
	var userTokenClaims UserTokenClaims
	isUserASeller, err := apiConfig.DB.IsUserASeller(r.Context(), userExist.ID)
	if err == nil {
		userTokenClaims = UserTokenClaims{
			UserID:   userExist.ID,
			SellerID: isUserASeller,
		}

		// checking if seller owns a store
		sellerStoreID, err := apiConfig.DB.GetStoreID(r.Context(), isUserASeller)
		if err == nil {
			userTokenClaims.StoreID = sellerStoreID
		}
	}

	// creating access token
	accessToken, err := MakeJWT(userTokenClaims, apiConfig.JwtSecret, time.Hour)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// generating refresh token
	refreshToken, err := GenerateRefereshToken()
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = apiConfig.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userExist.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// sending access token as response for logged in user
	utility.RespondWithJson(w, http.StatusOK, ResponseUser{
		Email:       userExist.Email,
		Phonenumber: userExist.PhoneNumber,
		AccessToken: accessToken,
		CreatedAt:   userExist.CreatedAt,
		UpdatedAt:   userExist.UpdatedAt,
	})
}

func MakeJWT(tokenClaims UserTokenClaims, tokenSecret string, expiresIn time.Duration) (string, error) {
	// creating the signing key to be used for signing the token
	signingKey := []byte(tokenSecret)

	// creating token claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "http://localhost:8080",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprintf("%s,%s,%s", tokenClaims.UserID, tokenClaims.SellerID, tokenClaims.StoreID),
	}

	// signing the access token with the signing key
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedAccessToken, err := accessToken.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func GenerateRefereshToken() (string, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(refreshToken), nil
}
