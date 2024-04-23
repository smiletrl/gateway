package payment

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	errors "github.com/smiletrl/gateway/pkg/error"
	"github.com/stretchr/testify/assert"

	"github.com/smiletrl/gateway/pkg/logger"
	"github.com/smiletrl/gateway/pkg/test"
)

var tests = []test.APITestCase{
	{
		Name:    Case1,
		Method:  "POST",
		URL:     "/payment",
		Context: nil,
		Body: `
		{
			"card": "5555555555554444",
			"expiry_date": "2023-12-23",
			"cvv": "123",
			"amount": "18.89",
			"currency": "CNY",
			"merchant_id": "12333"
		}
		`,
		ExpectStatus: http.StatusOK,
		ExpectResponse: `
		{
			"data": "ok"
		}`,
	},
	{
		Name:    Case2,
		Method:  "POST",
		URL:     "/payment",
		Context: nil,
		Body: `
		{
			"card": "4111111111111111",
			"expiry_date": "2023-12-23",
			"cvv": "123",
			"amount": "18.89",
			"currency": "CNY",
			"merchant_id": "12333"
		}
		`,
		ExpectStatus: http.StatusOK,
		ExpectResponse: `
		{
			"data": "ok"
		}`,
	},
	{
		Name:    Case3,
		Method:  "POST",
		URL:     "/payment",
		Context: nil,
		Body: `
		{
			"card": "40000000",
			"expiry_date": "2023-12-23",
			"cvv": "123",
			"amount": "18.89",
			"currency": "CNY",
			"merchant_id": "12333"
		}
		`,
		ExpectStatus: http.StatusBadRequest,
		ExpectResponse: `
		{
			"code": "error",
			"message": "card number is invalid"
		}`,
	},
}

func TestAPI(t *testing.T) {
	e := echo.New()

	logger := logger.NewProvider()
	e.Use(errors.Middleware(logger))

	// ctx, cancel := test.NewContext()
	// defer cancel()

	echo.New()
	g := e.Group("")

	repo := NewRepository()
	svc := NewService(repo)

	RegisterHandlers(g, svc)

	repoD, ok := repo.(*repository)
	assert.True(t, ok)

	verify := &VerifyTest{repo: repoD}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actualResponse := test.Endpoint(t, e, tc)
			verify.Switch(tc, actualResponse, t)

			// delete all transactions generated in this test case
			repoD.mockDB.Range(func(key, value interface{}) bool {
				repoD.mockDB.Delete(key)
				return true
			})
		})
	}
}
