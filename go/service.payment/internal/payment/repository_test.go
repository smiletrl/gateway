package payment

import (
	"context"
	"sync"
	"testing"

	"github.com/smiletrl/gateway/service.payment/internal/payment/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTransactionID(t *testing.T) {
	repo := &repository{}

	// test our transaction id generation can work for high concurrent scenarios.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transID, err := repo.generateTransactionID()
			assert.NoError(t, err)
			assert.NotEmpty(t, transID)
		}()
	}
	wg.Wait()
}

func TestCreateTransaction(t *testing.T) {
	repo := &repository{}
	cases := []struct {
		name string
		req  model.CreatePaymentRequest
	}{
		{
			name: "case 1",
			req: model.CreatePaymentRequest{
				Payment: model.Payment{
					Card:       "12",
					ExpiryDate: "2022-12-12",
					Cvv:        "234",
					Amount:     "12.89",
					Currency:   "CNY",
					MerchantID: "xx",
				},
			},
		},
		{
			name: "case 2",
			req: model.CreatePaymentRequest{
				Payment: model.Payment{
					Card:       "15",
					ExpiryDate: "2024-12-12",
					Cvv:        "234",
					Amount:     "12.89",
					Currency:   "CNY",
					MerchantID: "yy",
				},
			},
		},
	}
	ctx := context.Background()
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			transID, err := repo.CreateTransaction(ctx, &ca.req)
			assert.NoError(t, err)
			assert.NotEmpty(t, transID)

			// verify new transaction exists at storage now
			transaction, ok := repo.mockDB.Load(transID)
			assert.True(t, ok)
			trans := transaction.(model.Transaction)
			expectedTrans := model.Transaction{
				ID:      transID,
				Status:  model.TransactionStatusPending,
				Payment: ca.req.Payment,
			}
			assert.Equal(t, expectedTrans, trans)
		})
	}
}

func TestUpdateTransactionStatus(t *testing.T) {
	repo := &repository{}
	// create one transaction firstly
	req := &model.CreatePaymentRequest{
		Payment: model.Payment{
			Card:       "12",
			ExpiryDate: "2022-12-12",
			Cvv:        "234",
			Amount:     "12.89",
			Currency:   "CNY",
			MerchantID: "xx",
		},
	}
	ctx := context.Background()
	validTransID, err := repo.CreateTransaction(ctx, req)
	assert.NoError(t, err)

	cases := []struct {
		name        string
		transID     string
		transStatus model.TransactionStatus
		hasError    bool
	}{
		{
			name:        "case 1",
			transID:     validTransID,
			transStatus: model.TransactionStatusApproved,
			hasError:    false,
		},
		{
			name:        "case 1",
			transID:     validTransID,
			transStatus: model.TransactionStatusDenied,
			hasError:    false,
		},
		{
			name:        "case 3",
			transID:     validTransID,
			transStatus: model.TransactionStatus("Unknown status"),
			hasError:    true,
		},
	}
	for _, ca := range cases {
		t.Run(ca.name, func(t *testing.T) {
			err := repo.UpdateTransactionStatus(ctx, ca.transID, ca.transStatus)
			// if case has error, just return
			if ca.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// if no error, verify status has changed at storage successfully
			transaction, ok := repo.mockDB.Load(ca.transID)
			assert.True(t, ok)
			trans := transaction.(model.Transaction)
			assert.Equal(t, ca.transStatus, trans.Status)
		})
	}
}
