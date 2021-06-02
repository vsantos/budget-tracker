package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	mongodbDatabase          = "budget-tracker"
	mongodbUserCollection    = "users"
	mongodbCardsCollection   = "cards"
	mongodbBalanceCollection = "balance"
	mongodbSpendsCollection  = "spends"
)

// User struct defines a user
// swagger:model
type User struct {
	// swagger:ignore
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// example: vsantos
	Login string `json:"login" bson:"login"`
	// example: Victor
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	// example: Santos
	Lastname string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	// example: vsantos.py@gmail.com
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	// example: myplaintextpassword
	SaltedPassword string `json:"password,omitempty" bson:"password,omitempty"`
	// swagger:ignore
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// JWTUser defines a user to generate JWT tokens
// swagger:model
type JWTUser struct {
	// ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// example: vsantos
	Login string `json:"login" bson:"login"`
	// example: myplaintextpassword
	Password string `json:"password" bson:"password"`
}

// JWTResponse returns as HTTP response the user details (to be used along with the generated JWT token)
// swagger:model
type JWTResponse struct {
	Type         string        `json:"type"`
	RefreshToken string        `json:"refresh"`
	AccessToken  string        `json:"token"`
	Details      SanitizedUser `json:"details,omitempty"`
}

// SanitizedUser defines a sanited user to GET purposes
// swagger:model
type SanitizedUser struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Login     string             `json:"login" bson:"login"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
}

// CreditCard defines a user credit card
// swagger:model
type CreditCard struct {
	// swagger:ignore
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// example: 5f4e76699c362be701856be6
	OwnerID primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	// example: My Platinum Card
	Alias string `json:"alias" bson:"alias"`
	// example: VISA
	Network string `json:"network" bson:"network"`
	// example: #ffffff
	Color string `json:"color" bson:"color"`
	// example: 1234
	LastDigits int32 `json:"last_digits" bson:"last_digits"`
	// swagger:ignore
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// Income defines an user outcome for a certain month
type Income struct {
	GrossIncome float64 `json:"gross" bson:"gross"`
	NetIncome   float64 `json:"net" bson:"net"`
}

// Outcome defines an user outcome for a certain month
type Outcome struct {
	FixedOutcome   float64 `json:"fixed" bson:"fixed"`
	DynamicOutcome float64 `json:"dynamic" bson:"dynamic"`
}

// swagger:model
// PaymentMethod defines which payment method was used for a certain spend
type PaymentMethod struct {
	Credit      CreditCard `json:"credit,omitempty" bson:"credit,omitempty"`
	Debit       bool       `json:"debit,omitempty" bson:"debit,omitempty"`
	PaymentSlip bool       `json:"payment_slip,omitempty" bson:"payment_slip,omitempty"`
}

// swagger:model
// Spend defines a user spend to be added to Balance
type Spend struct {
	// swagger:ignore
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	// example: fixed
	Type string `json:"type" bson:"type"`
	// example: guitar lessons
	Description string `json:"description" bson:"description"`
	// example: 12.90
	Cost float64 `json:"cost" bson:"cost"`
	// example: debit: true
	PaymentMethod PaymentMethod `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	// example: "categories": ["personal development"]
	Categories []string `json:"category,omitempty" bson:"category,omitempty"`
	// swagger:ignore
	CreatedAt primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// swagger:model
// Balance defines an user balance
type Balance struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerID         primitive.ObjectID `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Income          Income             `json:"income,omitempty" bson:"income,omitempty"`
	Outcome         Outcome            `json:"outcome" bson:"outcome"`
	SpendableAmount float64            `json:"spendable_amount" bson:"spendable_amount"`
	Historic        []Spend            `json:"historic" bson:"historic"`
	Currency        string             `json:"currency" bson:"currency"`
	Month           int64              `json:"month" bson:"month"`
	Year            int64              `json:"year" bson:"year"`
	CreatedAt       primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt       primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
