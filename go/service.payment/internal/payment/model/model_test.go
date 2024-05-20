package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionStatusIsvalid(t *testing.T) {
	cases := []struct {
		name              string
		transactionStatus TransactionStatus
		isValid           bool
	}{
		{
			name:              "case 1",
			transactionStatus: TransactionStatusApproved,
			isValid:           true,
		},
		{
			name:              "case 2",
			transactionStatus: TransactionStatusDenied,
			isValid:           true,
		},
		{
			name:              "case 3",
			transactionStatus: TransactionStatusPending,
			isValid:           true,
		},
		{
			name:              "case 4",
			transactionStatus: TransactionStatus("unknown"),
			isValid:           false,
		},
	}

	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			assert.Equal(t, ca.isValid, ca.transactionStatus.IsValid())
		})
	}
}

func TestIsCardValid(t *testing.T) {
	cases := []struct {
		name    string
		card    string
		isValid bool
	}{
		{
			name:    "case 1",
			card:    "4000000",
			isValid: false,
		},
		{
			name:    "case 2",
			card:    "5555555555554444",
			isValid: true,
		},
		{
			name:    "case 3",
			card:    "4111111111111111",
			isValid: true,
		},
	}

	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			pay := Payment{
				Card: ca.card,
			}
			assert.Equal(t, ca.isValid, pay.isCardValid())
		})
	}
}
