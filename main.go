package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/harshvardha/e-commerce/controllers"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/twilio/twilio-go"
)

func main() {
	// loading all the env variables from .env file
	godotenv.Load()

	// port env variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port env variable not set")
	}

	// jwt secret varaible
	jwtSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if jwtSecret == "" {
		log.Fatal("jwt_secret variable not set")
	}

	// db_url env variable
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbUrl env variable not set")
	}

	// twilio account sid variable
	TWILIO_ACCOUNT_SID := os.Getenv("TWILIO_ACCOUNT_SID")
	if TWILIO_ACCOUNT_SID == "" {
		log.Fatal("TWILIO_ACCOUNT_SID env variable not set")
	}

	// twilio auth token
	TWILIO_AUTH_TOKEN := os.Getenv("TWILIO_AUTH_TOKEN")
	if TWILIO_AUTH_TOKEN == "" {
		log.Fatal("TWILIO_AUTH_TOKEN env variable not set")
	}

	// twilio service sid
	VERIFY_SERVICE_SID := os.Getenv("VERIFY_SERVICE_SID")
	if VERIFY_SERVICE_SID == "" {
		log.Fatal("VERIFY_SERVICE_SID env variable not set")
	}

	// creating twilio rest client
	var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})

	// setting the data validator
	dataValidator := utility.Validator{
		Validate: validator.New(),
	}

	// registering custom password validator
	dataValidator.Validate.RegisterValidation("password", utility.CustomPasswordValidator)

	// setting twilio config
	twilioConfig := controllers.TwilioConfig{
		VERIFY_SERVICE_SID: VERIFY_SERVICE_SID,
		Client:             client,
		DataValidator:      dataValidator,
	}

	//creating database connection
	dbConnection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	db := database.New(dbConnection)

	// setting apiConfig struct
	apiCfg := controllers.ApiConfig{
		DB:        db,
		JwtSecret: jwtSecret,
	}

	// setting ApiTwilioConfig struct
	apiTwilioConfig := controllers.ApiTwilioConfig{
		ApiCfg:    apiCfg,
		TwilioCfg: twilioConfig,
	}

	// creating new server
	mux := http.NewServeMux()

	// healthz api endpoint to check whether the server is successfully setup or not
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; chatset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// api endpoints for user
	mux.HandleFunc("POST /api/auth/sendOTP", twilioConfig.HandleSendOTP)
	mux.HandleFunc("POST /api/auth/verifyOTP", apiTwilioConfig.HandleVerifyOTP)
	mux.HandleFunc("POST /api/auth/resendOTP", twilioConfig.HandleResendOTP)
	mux.HandleFunc("POST /api/auth/login", apiCfg.HandleLogin)

	// starting the server
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
