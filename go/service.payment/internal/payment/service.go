package payment

import (
	"context"
	"fmt"

	"github.com/smiletrl/gateway/service.payment/internal/payment/model"
)

// Service is for main business logic
type Service interface {
	// create a new payment.
	Create(ctx context.Context, req *model.CreatePaymentRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, req *model.CreatePaymentRequest) error {
	// validate request
	if err := req.Validate(); err != nil {
		return fmt.Errorf("error validating: %w", err)
	}

	// create transaction firstly
	transactionID, err := s.repo.CreateTransaction(ctx, req)
	if err != nil {
		return fmt.Errorf("error creating transaction at Create: %w", err)
	}

	// acquirer action
	isApproved := s.isAcquirerApproved(req.Card)

	var newStatus model.TransactionStatus
	if isApproved {
		newStatus = model.TransactionStatusApproved
	} else {
		newStatus = model.TransactionStatusDenied
	}

	// update transaction status as per acquirer's action
	if err := s.repo.UpdateTransactionStatus(ctx, transactionID, newStatus); err != nil {
		return fmt.Errorf("error updating transaction status at Create: %w", err)
	}
	return nil
}
