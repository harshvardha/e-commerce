package controllers

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/harshvardha/e-commerce/internal/database"
	"github.com/harshvardha/e-commerce/utility"
	"github.com/twilio/twilio-go"
)

type ApiConfig struct {
	DB        *database.Queries
	JwtSecret string
}

type TwilioConfig struct {
	VERIFY_SERVICE_SID string
	Client             *twilio.RestClient
	DataValidator      utility.Validator
}

type ApiTwilioConfig struct {
	ApiCfg    ApiConfig
	TwilioCfg TwilioConfig
}

type SendOtpRequest struct {
	Email       string `json:"email"`
	Phonenumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type VerifyOtpRequest struct {
	VerificationToken string `json:"verification_token"`
	OTP               string `json:"otp"`
	RequestType       int    `json:"request_type"`
}

type SendOTPResponse struct {
	VerificationToken string `json:"verification_token"`
}

type ResendOTPRequest struct {
	VerificationToken string `json:"verification_token"`
	RequestType       int    `json:"request_type"`
}

type LoginRequest struct {
	Phonenumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserTokenClaims struct {
	UserID   uuid.UUID
	SellerID string
	StoreID  uuid.UUID
}

type ResponseUser struct {
	Email       string    `json:"email"`
	Phonenumber string    `json:"phone_number"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ID struct {
	UserID   uuid.UUID
	SellerID string
	StoreID  uuid.UUID
	AdminID  uuid.UUID
}

type UpdateUserInfo struct {
	Email       string `json:"email"`
	Phonenumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type CreateSellerRequest struct {
	GstNumber             string `json:"gst_number"`
	PanNumber             string `json:"pan_number"`
	PickupAddress         string `json:"pickup_address"`
	BankAccountHolderName string `json:"bank_account_holder_name"`
	BankAccountNumber     string `json:"bank_account_number"`
	IFSCCode              string `json:"ifsc_code"`
}

type UpdateSellerTaxAndAddress struct {
	GstNumber     string `json:"gst_number"`
	PanNumber     string `json:"pan_number"`
	PickupAddress string `json:"pickup_address"`
}

type UpdateSellerBankDetails struct {
	BankAccountHolderName string `json:"bank_account_holder_name"`
	BankAccountNumber     string `json:"bank_account_number"`
	IFSCCode              string `json:"ifsc_code"`
}

type CreateSellerResponse struct {
	ID                    string    `json:"id"`
	GstNumber             string    `json:"gst_number"`
	PanNumber             string    `json:"pan_number"`
	PickupAddress         string    `json:"pickup_address"`
	BankAccountHolderName string    `json:"bank_account_holder_name"`
	BankAccountNumber     string    `json:"bank_account_number"`
	IFSCCode              string    `json:"ifsc_code"`
	AccessToken           string    `json:"access_token"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type SellerContactInfoResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Phonenumber string    `json:"phone_number"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SellerTaxAndAddressInfoResponse struct {
	ID            string    `json:"id"`
	GstNumber     string    `json:"gst_number"`
	PanNumber     string    `json:"pan_number"`
	PickupAddress string    `json:"pickup_address"`
	AccessToken   string    `json:"access_token"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SellerBankDetailsResponse struct {
	ID                    string    `json:"id"`
	BankAccountHolderName string    `json:"bank_account_holder_name"`
	BankAccountNumber     string    `json:"bank_account_number"`
	IFSCCode              string    `json:"ifsc_code"`
	AccessToken           string    `json:"access_token"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type CreateOrUpdateStoreRequest struct {
	Name string `json:"name"`
}

type CreateStoreResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SellerID    string    `json:"seller_id"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateStoreResponse struct {
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmptyResponse struct {
	AccessToken string `json:"access_token"`
}

type GetStoreInformationResponse struct {
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminResponse struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phonenumber string    `json:"phonenumber"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateAdmin struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateAdminPasswordOrPhonenumber struct {
	Phonenumber          string `json:"phonenumber"`
	Password             string `json:"password"`
	OTPVerificationToken string `json:"otp_verification_token"`
	OTP                  string `json:"otp"`
}

type UpdateAdminPhonenumberResponse struct {
	Phonenumber string    `json:"phonenumber"`
	AccessToken string    `json:"access_token"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrUpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrUpdateProductRequest struct {
	Name            string            `json:"name"`
	Description     map[string]string `json:"description"`
	Characteristics map[string]string `json:"characteristics"`
	Price           float64           `json:"price"`
	ImageUrls       map[string]string `json:"image_urls"`
	StockAmount     int32             `json:"stock_amount"`
	CategoryID      uuid.UUID         `json:"category_id"`
}

type ProductResponse struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     json.RawMessage `json:"description"`
	Characteristics json.RawMessage `json:"characteristics"`
	Price           float64         `json:"price"`
	ImageUrls       json.RawMessage `json:"image_urls"`
	StockAmount     int32           `json:"stock_amount"`
	StoreID         uuid.UUID       `json:"store_id"`
	CategoryID      uuid.UUID       `json:"category_id"`
	AccessToken     string          `json:"access_token"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ProductsByCategoryOrStoreID struct {
	Products    []database.Product `json:"products"`
	AccessToken string             `json:"access_token"`
}

type CustomerResponse struct {
	ID              uuid.UUID `json:"id"`
	DeliveryAddress string    `json:"delivery_address"`
	Pincode         string    `json:"pincode"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	AccessToken     string    `json:"access_token"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpdateCustomerAddress struct {
	DeliveryAddress string    `json:"delivery_address"`
	Pincode         string    `json:"pincode"`
	City            uuid.UUID `json:"city"`
	State           uuid.UUID `json:"state"`
}

type AddOrUpdateStateRequest struct {
	Name string `json:"name"`
}

type StateResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AddOrUpdateCityRequest struct {
	Name  string    `json:"name"`
	State uuid.UUID `json:"state"`
}

type CityResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	State       string    `json:"state"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrUpdateReviewRequest struct {
	Description string `json:"description"`
}

type ReviewResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AddProductToCartRequest struct {
	Quantity int32 `json:"quantity"`
}

type OrderDetailsResponse struct {
	OrderDetails []database.GetOrderDetailsRow `json:"order_details"`
	AccessToken  string                        `json:"access_token"`
}
