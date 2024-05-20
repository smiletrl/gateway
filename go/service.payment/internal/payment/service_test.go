package payment

import (
	"context"
	"fmt"
	"testing"

	"github.com/smiletrl/gateway/service.payment/internal/payment/mocks"
	"github.com/smiletrl/gateway/service.payment/internal/payment/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsAcquirerApproved(t *testing.T) {
	svc := &service{}
	cases := []struct {
		name       string
		card       string
		isApproved bool
	}{
		{
			name:       "case even",
			card:       "18",
			isApproved: true,
		},
		{
			name:       "case odd",
			card:       "19",
			isApproved: false,
		},
	}
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			isApproved := svc.isAcquirerApproved(ca.card)
			assert.Equal(t, ca.isApproved, isApproved)
		})
	}
}

func TestCreate(t *testing.T) {
	// mock repository
	repo := mocks.NewRepository(t)
	repo.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, req *model.CreatePaymentRequest) (transactionID string, err error) {
			return "transaction-id", nil
		})
	repo.On("UpdateTransactionStatus", mock.Anything, mock.Anything, mock.Anything).
		Return(func(ctx context.Context, transactionID string, transactionStatus model.TransactionStatus) error {
			if transactionStatus == model.TransactionStatusApproved {
				return nil
			}
			// if acquirer denies, mock an error returned
			return fmt.Errorf("please approve!")
		})

	svc := NewService(repo)
	cases := []struct {
		name     string
		req      model.CreatePaymentRequest
		hasError bool
	}{
		{
			name: "case no error",
			req: model.CreatePaymentRequest{
				Payment: model.Payment{
					Card: "18",
				},
			},
			hasError: false,
		},
		{
			name: "case error out",
			req: model.CreatePaymentRequest{
				Payment: model.Payment{
					Card: "19",
				},
			},
			hasError: true,
		},
	}

	ctx := context.Background()

	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			err := svc.Create(ctx, &ca.req)
			if ca.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
