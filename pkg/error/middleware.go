package errors

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"

	"github.com/smiletrl/gateway/pkg/constant"
	"github.com/smiletrl/gateway/pkg/logger"
)

func Middleware(logger logger.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, true)
					msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					logger.Errorf(msg)

					c.Error(err)
				}
			}()

			// set logger.
			c.Set(constant.Logger, logger)
			return next(c)
		}
	}
}
