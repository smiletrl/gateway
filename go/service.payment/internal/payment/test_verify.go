package payment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smiletrl/gateway/pkg/test"
	"github.com/smiletrl/gateway/service.payment/internal/payment/model"
)

const (
	Case1 = "create a new payment -- status approved"
	Case2 = "create a new payment -- status denied"
	Case3 = "invalid card"
)

// Verify result more than http response result
type VerifyTest struct {
	repo *repository
}

func (v *VerifyTest) Switch(tc test.APITestCase, actualResponse string, t *testing.T) {
	switch tc.Name {
	case Case1:
		v.VerifyCase1(tc, actualResponse, t)
	case Case2:
		v.VerifyCase2(tc, actualResponse, t)
	}
}

func (v *VerifyTest) VerifyCase1(tc test.APITestCase, actualResponse string, t *testing.T) {
	v.repo.mockDB.Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		assert.True(t, ok)

		expectedTrans := model.Transaction{
			ID:     keyStr,
			Status: model.TransactionStatusApproved,
			Payment: model.Payment{
				Card:       "5555555555554444",
				ExpiryDate: "2023-12-23",
				Cvv:        "123",
				Amount:     "18.89",
				Currency:   "CNY",
				MerchantID: "12333",
			},
		}
		trans, ok := value.(model.Transaction)
		assert.True(t, ok)
		assert.Equal(t, expectedTrans, trans)
		return true
	})
}

func (v *VerifyTest) VerifyCase2(tc test.APITestCase, actualResponse string, t *testing.T) {
	v.repo.mockDB.Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		assert.True(t, ok)

		expectedTrans := model.Transaction{
			ID:     keyStr,
			Status: model.TransactionStatusDenied,
			Payment: model.Payment{
				Card:       "4111111111111111",
				ExpiryDate: "2023-12-23",
				Cvv:        "123",
				Amount:     "18.89",
				Currency:   "CNY",
				MerchantID: "12333",
			},
		}
		trans, ok := value.(model.Transaction)
		assert.True(t, ok)
		assert.Equal(t, expectedTrans, trans)
		return true
	})
}
