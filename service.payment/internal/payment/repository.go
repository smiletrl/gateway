package payment

import (
	"context"
	"fmt"
	"sync"

	"github.com/smiletrl/gateway/service.payment/internal/payment/model"
)

// Repository is for db storage

//go:generate mockery --name=Repository
type Repository interface {
	// create a new transaction, and store it in db.
	CreateTransaction(ctx context.Context, req *model.CreatePaymentRequest) (transactionID string, err error)

	// update existing transaction's status
	UpdateTransactionStatus(ctx context.Context, transactionID string, transactionStatus model.TransactionStatus) error
}

type repository struct {
	// a mock db in map, with key being transaction id, and value being its transaction detail.
	// use sync.Map instead of regular map to handle concurrent requests.
	mockDB sync.Map
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateTransaction(ctx context.Context, req *model.CreatePaymentRequest) (transactionID string, err error) {
	transactionID, err = r.generateTransactionID()
	if err != nil {
		return "", fmt.Errorf("error generating transaction id: %w", err)
	}

	// save this new transaction into storage
	r.mockDB.Store(transactionID, model.Transaction{
		ID:      transactionID,
		Status:  model.TransactionStatusPending,
		Payment: req.Payment,
	})

	return transactionID, nil
}

func (r *repository) UpdateTransactionStatus(ctx context.Context, transactionID string, transactionStatus model.TransactionStatus) error {
	// transaction status validate
	if !transactionStatus.IsValid() {
		return fmt.Errorf("unknown transaction status: %s", transactionStatus)
	}

	transaction, ok := r.mockDB.Load(transactionID)
	if !ok {
		return fmt.Errorf("unknown transaction id: %s", transactionID)
	}

	// update transaction status with new status and save it into storage
	trans := transaction.(model.Transaction)
	trans.Status = transactionStatus
	r.mockDB.Store(transactionID, trans)

	return nil
}
