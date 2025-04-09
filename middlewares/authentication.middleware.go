package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/controllers"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
)

type authedHandler func(http.ResponseWriter, *http.Request, controllers.ID, string)

func ValidateJWT(handler authedHandler, tokenSecret string, db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracting JWT token from authentication header
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) < 2 {
			utility.RespondWithError(w, http.StatusUnauthorized, "Access Token Malformed.")
			return
		}

		// initializing and empty struct to parse JWT claims
		claimsStruct := jwt.RegisteredClaims{}

		// parsing the token
		token, parseError := jwt.ParseWithClaims(authHeader[1], &claimsStruct, func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})

		tokenSubject, err := token.Claims.GetSubject()
		if err != nil {
			utility.RespondWithError(w, http.StatusUnauthorized, "Access Token Malformed.")
			return
		}

		// extracting user_id, seller_id and store_id from subject
		subjects := strings.Split(tokenSubject, ",")
		if len(subjects) == 0 {
			utility.RespondWithError(w, http.StatusUnauthorized, "Access Token Malformed.")
			return
		}

		// user id
		userId, err := uuid.Parse(subjects[0])
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, "Unable to parse user id in middleware.")
			return
		}

		// seller id
		sellerId := subjects[1]

		// store id
		storeId, _ := uuid.Parse(subjects[2])

		// creating the IDs struct
		IDs := controllers.ID{
			UserID:   userId,
			SellerID: sellerId,
			StoreID:  storeId,
		}

		// checking if the token is expired or not
		if parseError != nil {
			// checking if the access token is expired or not
			tokenExpiresAt, err := token.Claims.GetExpirationTime()
			if err != nil {
				utility.RespondWithError(w, http.StatusUnauthorized, "Access Token Malformed.")
				return
			}

			// if access token is expired then we will check if the refresh token is expired or not
			// if refresh token is not expired then we will create a new access token and continue
			// if refresh token is expired then we will ask user to login again
			if time.Now().After(tokenExpiresAt.Time) {
				refreshTokenExpirationAt, err := db.GetRefreshToken(r.Context(), userId)
				if err != nil {
					utility.RespondWithError(w, http.StatusNotFound, "Refresh Token Not Found.")
					return
				}
				if time.Now().After(refreshTokenExpirationAt) {
					utility.RespondWithError(w, http.StatusUnauthorized, "Please login again")
					return
				} else {
					newTokenClaims := controllers.UserTokenClaims{
						UserID:   userId,
						SellerID: sellerId,
						StoreID:  storeId,
					}
					newAccessToken, err := controllers.MakeJWT(newTokenClaims, tokenSecret, time.Hour)
					if err != nil {
						utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
						return
					}
					handler(w, r, IDs, newAccessToken)
				}
			}

			utility.RespondWithError(w, http.StatusUnauthorized, parseError.Error())
			return
		}
		handler(w, r, IDs, "")
	}
}
