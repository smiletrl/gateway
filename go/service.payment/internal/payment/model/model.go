package model

import (
	errorpkg "github.com/smiletrl/gateway/pkg/error"
)

// represent one transaction's status
type TransactionStatus string

func (t TransactionStatus) IsValid() bool {
	if t == TransactionStatusApproved || t == TransactionStatusDenied || t == TransactionStatusPending {
		return true
	}
	return false
}

const (
	// all available status
	TransactionStatusPending  TransactionStatus = "Pending"
	TransactionStatusApproved TransactionStatus = "Approved"
	TransactionStatusDenied   TransactionStatus = "Denied"
)

type Payment struct {
	Card       string      `json:"card"`
	ExpiryDate string      `json:"expiry_date"`
	Cvv        string      `json:"cvv"`
	Amount     MoneyString `json:"amount"`
	Currency   string      `json:"currency"`
	MerchantID string      `json:"merchant_id"`
}

func (p *Payment) Validate() error {
	if err := p.Amount.Validate(); err != nil {
		return errorpkg.BadRequestMessagef("amount is invalid")
	}

	if isCardValid := p.isCardValid(); !isCardValid {
		return errorpkg.BadRequestMessagef("card number is invalid")
	}
	// @todo validate other fields, like
	// -- ExpiryDate needs to be in string format `2023-12-12`
	// -- Cvv needs to be three numbers
	// -- Currency should be one of the available currencies in the world
	// -- Maybe MechaintID should exist in our system

	// @todo add test converage for above validation
	return nil
}

// source code from https://dev.to/claudbytes/build-a-credit-card-validator-using-go-5d2b
func (p *Payment) isCardValid() bool {
	// this function implements the luhn algorithm
	// it takes as argument a cardnumber of type string
	// and it returns a boolean (true or false) if the
	// card number is valid or not

	// initialise a variable to keep track of the total sum of digits
	total := 0
	// Initialize a flag to track whether the current digit is the second digit from the right.
	isSecondDigit := false

	// iterate through the card number digits in reverse order
	for i := len(p.Card) - 1; i >= 0; i-- {
		// conver the digit character to an integer
		digit := int(p.Card[i] - '0')

		if isSecondDigit {
			// double the digit for each second digit from the right
			digit *= 2
			if digit > 9 {
				// If doubling the digit results in a two-digit number,
				//subtract 9 to get the sum of digits.
				digit -= 9
			}
		}

		// Add the current digit to the total sum
		total += digit

		//Toggle the flag for the next iteration.
		isSecondDigit = !isSecondDigit
	}

	// return whether the total sum is divisible by 10
	// making it a valid luhn number
	return total%10 == 0
}

type Transaction struct {
	// transaction id is should be globally unique
	ID     string            `json:"id"`
	Status TransactionStatus `json:"status"`
	Payment
}

// MoneyString is for money amount string, like `19.89`, `78.00`, `67`.
type MoneyString string

func (m *MoneyString) Validate() error {
	// @todo, validate money string following pattern like `19.89`.
	// -- all numbers, and no more than two decimal points, and can be converted to int64 `1989`.
	return nil
}

type CreatePaymentRequest struct {
	Payment
}

func (c *CreatePaymentRequest) Validate() error {
	return c.Payment.Validate()
}

type OKResponse struct {
	Data string `json:"data"`
}
