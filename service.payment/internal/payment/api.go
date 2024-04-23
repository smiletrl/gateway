package payment

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/smiletrl/gateway/service.payment/internal/payment/model"

	errorpkg "github.com/smiletrl/gateway/pkg/error"
)

func RegisterHandlers(g *echo.Group, svc Service) {
	res := &resource{svc}
	payG := g.Group("/payment")
	payG.POST("", res.create)
}

type resource struct {
	svc Service
}

func (r *resource) create(c echo.Context) error {
	req := new(model.CreatePaymentRequest)
	if err := c.Bind(req); err != nil {
		return errorpkg.BadRequest(c, err)
	}
	ctx := c.Request().Context()
	if err := r.svc.Create(ctx, req); err != nil {
		return errorpkg.Abort(c, err)
	}
	return c.JSON(http.StatusOK, model.OKResponse{Data: "ok"})
}
