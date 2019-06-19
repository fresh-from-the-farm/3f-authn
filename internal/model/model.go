package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"time"
)

const (
	Customer        = "Customer"
	Employee        = "Employee"
	AmericanExpress = "American Express"
	DinersClub      = "Diners Club"
	Discover        = "Discover"
	JCB             = "JCB"
	MasterCard      = "MasterCard"
	Visa            = "Visa"
	UnknownCard     = "Unknown"
)

type AccountOperationResponse struct {
	Account Account  `json:"account,omitempty"`
	Errors  []string `json:"errors,omitempty"`
	Status  string   `json:"status,omitempty"`
	Message string   `json:"message,omitempty"`
}

type AccessToken struct {
	AccessToken string
}

type Account struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName      string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName       string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Relation       []string           `json:"relation,omitempty" bson:"relation,omitempty"`
	Username       string             `json:"userName,omitempty" bson:"userName,omitempty"`
	Password       string             `json:"password,omitempty" bson:"password,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty"`
	PaymentMethods []Card             `json:"paymentMethods,omitempty" bson:"paymentMethods,omitempty"`
	Addresses      []Address          `json:"addresses,omitempty" bson:"addresses,omitempty"`
	Active         bool               `json:"isActive" bson:"isActive"`
	EmailVerified  bool               `json:"emailVerified,omitempty" bson:"emailVerified,omitempty"`
	Created        *time.Time         `json:"created,omitempty" bson:"created,omitempty"`
	Updated        *time.Time         `json:"updated,omitempty" bson:"updated,omitempty"`
}

type Address struct {
	StreetAddress      []string `json:"streetAddress,omitempty" bson:"streetAddress,omitempty"`
	City               string   `json:"city,omitempty" bson:"city,omitempty"`
	State              string   `json:"state,omitempty" bson:"state,omitempty"`
	ZipCode            string   `json:"zipCode,omitempty" bson:"zipCode,omitempty"`
	Country            string   `json:"country,omitempty" bson:"country,omitempty"`
	AccessInstructions []string `json:"accessInstructions,omitempty" bson:"accessInstructions,omitempty"`
	AccessCode         string   `json:"accessCode,omitempty" bson:"accessCode,omitempty"`
}

type Card struct {
	Id                string `json:"id" bson:"id"`
	Name              string `json:"name,omitempty" bson:"name,omitempty""`
	Type              string `json:"type" bson:"type"`
	ExpMonth          int    `json:"exp_month" bson:"exp_month"`
	ExpYear           int    `json:"exp_year" bson:"exp_year"`
	Last4             string `json:"last4" bson:"last4"`
	Fingerprint       string `json:"fingerprint" bson:"fingerprint"`
	Country           string `json:"country,omitempty" bson:"country,omitempty"`
	Address1          string `json:"address_line1,omitempty" bson:"address_line1,omitempty"`
	Address2          string `json:"address_line2,omitempty" bson:"address_line2,omitempty"`
	AddressCountry    string `json:"address_country,omitempty" bson:"address_country,omitempty"`
	AddressState      string `json:"address_state,omitempty" bson:"address_state,omitempty"`
	AddressZip        string `json:"address_zip,omitempty" bson:"address_zip,omitempty"`
	AddressCity       string `json:"address_city" bson:"address_city"`
	AddressLine1Check string `json:"address_line1_check,omitempty" bson:"address_line1_check,omitempty"`
	AddressZipCheck   string `json:"address_zip_check,omitempty" bson:"address_zip_check,omitempty"`
	CVCCheck          string `json:"cvc_check,omitempty" bson:"cvc_check,omitempty"`
}

func IsLuhnValid(card string) (bool, error) {

	var sum = 0
	var cardNumber = strings.Split(card, "")

	// iterate through the cardNumber in reverse order
	for i, even := len(cardNumber)-1, false; i >= 0; i, even = i-1, !even {

		// convert the digit to an integer
		digit, err := strconv.Atoi(cardNumber[i])
		if err != nil {
			return false, err
		}

		// we multiply every other digit by 2, adding the product to the sum.
		// note: if the product is double cardNumber (i.e. 14) we add the two cardNumber
		//       to the sum (14 -> 1+4 = 5). A simple shortcut is to subtract 9
		//       from a double digit product (14 -> 14 - 9 = 5).
		switch {
		case even && digit > 4:
			sum += (digit * 2) - 9
		case even:
			sum += digit * 2
		case !even:
			sum += digit
		}
	}

	// if the sum is divisible by 10, it passes the check
	return sum%10 == 0, nil
}

func GetCardType(card string) string {

	switch card[0:1] {
	case "4":
		return Visa
	case "2", "1":
		switch card[0:4] {
		case "2131", "1800":
			return JCB
		}
	case "6":
		switch card[0:4] {
		case "6011":
			return Discover
		}
	case "5":
		switch card[0:2] {
		case "51", "52", "53", "54", "55":
			return MasterCard
		}
	case "3":
		switch card[0:2] {
		case "34", "37":
			return AmericanExpress
		case "36":
			return DinersClub
		case "30":
			switch card[0:3] {
			case "300", "301", "302", "303", "304", "305":
				return DinersClub
			}
		default:
			return JCB
		}
	}

	return UnknownCard
}
