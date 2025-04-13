// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID          uuid.UUID
	Name        string
	Email       string
	Phonenumber string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Cart struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Quantity  int32
}

type Category struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Characteristic struct {
	ID          uuid.UUID
	Description json.RawMessage
	ProductID   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type City struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	State     uuid.UUID
}

type Customer struct {
	ID              uuid.UUID
	DeliveryAddress string
	Pincode         string
	City            uuid.UUID
	State           uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Order struct {
	ID         uuid.UUID
	TotalValue float64
	SellerID   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrdersUsersProduct struct {
	OrderID   uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID          uuid.UUID
	Name        string
	Description json.RawMessage
	Price       float64
	ImageUrls   json.RawMessage
	StockAmount int32
	StoreID     uuid.UUID
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RefreshToken struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Return struct {
	ID         uuid.UUID
	SellerID   string
	TotalValue float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ReturnsUsersProduct struct {
	ReturnID  uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Review struct {
	ID          uuid.UUID
	Description string
	UserID      uuid.UUID
	ProductID   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SavedItem struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Seller struct {
	ID                    string
	UserID                uuid.UUID
	GstNumber             string
	PanNumber             string
	PickupAddress         string
	BankAccountHolderName string
	BankAccountNumber     string
	IfscCode              string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type State struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store struct {
	ID        uuid.UUID
	Name      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID          uuid.UUID
	Email       string
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Firstname   string
	Lastname    string
}
