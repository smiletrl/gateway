package errors

import (
	errCore "errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"

	"github.com/smiletrl/gateway/pkg/constant"
	"github.com/smiletrl/gateway/pkg/logger"
)

// Error represents business error
type Error struct {
	// Code is for business error code, such as `invalid_username`, 'password_not_match'
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`

	// http status code
	Status int `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

// New is to create a new error
func New(message string, code ...string) error {
	if len(code) == 1 {
		return &Error{
			Code:    code[0],
			Message: message,
		}
	}
	// This is the default error
	return &Error{
		Code:    "error",
		Message: message,
	}
}

func BadRequestMessage(message string, code ...string) error {
	defaultCode := "error"
	if len(code) == 1 {
		defaultCode = code[0]
	}

	return &Error{
		Code:    defaultCode,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func BadRequestMessagef(message string, args ...interface{}) error {
	msg := fmt.Sprintf(message, args...)
	return &Error{
		Code:    "error",
		Message: msg,
		Status:  http.StatusBadRequest,
	}
}

func ForbiddenMessage(message string, code ...string) error {
	defaultCode := "error"
	if len(code) == 1 {
		defaultCode = code[0]
	}

	return &Error{
		Code:    defaultCode,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

// Response is error response
type Response struct {
	*Error
}

// Abort means error out
func Abort(c echo.Context, err error) error {
	return abort(c, err, http.StatusInternalServerError)
}

// BadRequest means bad request `400`
func BadRequest(c echo.Context, err error) error {
	return abort(c, err, http.StatusBadRequest)
}

func Unauthorized(c echo.Context, err error) error {
	return abort(c, err, http.StatusUnauthorized)
}

func Forbidden(c echo.Context, err error) error {
	return abort(c, err, http.StatusForbidden)
}

func PreconditionFailed(c echo.Context, err error) error {
	return abort(c, err, http.StatusPreconditionFailed)
}

func abort(c echo.Context, err error, status int) error {
	logger := c.Get(constant.Logger).(logger.Provider)

	req := c.Request()
	uri := req.RequestURI

	// defer - status could come from passed error's status value.
	defer func() {
		url := fmt.Sprintf("%s %s status: %d", req.Method, req.Host+uri, status)
		logger.Errorf("http request abort error URL: %s error: %v", url, err)
	}()

	// This is to deal with go above 1.13 error handling
	var customErr *Error
	if errCore.As(err, &customErr) {
		if customErr.Status != 0 {
			status = customErr.Status
		}
		return c.JSON(status, Response{
			Error: customErr,
		})
	}

	// Don't show internal error to frontend directly
	return c.JSON(status, Response{
		Error: &Error{
			Message: "System error, please contact technical support",
		},
	})
}

func Recover(r interface{}, logger logger.Provider) {
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}
	stack := make([]byte, 4<<10)
	length := runtime.Stack(stack, true)
	msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
	logger.Errorf(msg)
}
