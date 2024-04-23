package accesslog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/smiletrl/gateway/pkg/constant"
	"github.com/smiletrl/gateway/pkg/logger"
)

type Provider interface {
	Middleware() echo.MiddlewareFunc
}

type provider struct {
	logger logger.Provider
}

func NewProvider(logger logger.Provider) Provider {
	return provider{logger}
}

func (p provider) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			ctx := c.Request().Context()

			// request body
			reqBody, err := io.ReadAll(req.Body)
			if err == nil {
				// now set request body back to request
				req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
				c.SetRequest(req)
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			// response body
			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &BodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			// request id is for global debug/track, get from header firstly
			reqID := c.Request().Header.Get(constant.HttpRequestIDHeader)
			if reqID == "" {
				reqID = uuid.New().String()
			}

			// set the request id into http request context. We won't pass
			// echo context through the whole http life cycle, but the http
			// request context will be passed through life cycle.
			// We will use the request id when we want to track request.
			ctx = context.WithValue(ctx, constant.RequestID, reqID)
			r := c.Request().WithContext(ctx)
			c.SetRequest(r)

			// execute other middleware.
			nextErr := next(c)
			if nextHttpErr, ok := nextErr.(*echo.HTTPError); ok {
				res.Status = nextHttpErr.Code
			} else if nextErr != nil {
				res.Status = http.StatusInternalServerError
			}

			// log this request
			if res.Status != 200 || req.URL.Path != "/health" {

				var (
					requestHeader, responseHeader string
				)

				var sh, _ = time.LoadLocation("Asia/Shanghai")
				var layout = "2006-01-02 15:04:05"
				stop := time.Now()

				// request header
				requestHeaderBytes, err := json.Marshal(req.Header)
				if err != nil {
					log.Printf("error on request hader read: %v", err)
				} else {
					requestHeader = string(requestHeaderBytes)
				}

				// response header
				responseHeaderBytes, err := json.Marshal(c.Response().Header())
				if err != nil {
					log.Printf("error on response hader read: %v", err)
				} else {
					responseHeader = string(responseHeaderBytes)
				}

				contextFields := extractContextFields(c.Request().Context())

				// for file upload, request body is binary data and big. Avoid logging such big data
				fileParam, _ := c.FormFile("file")
				if fileParam != nil {
					reqBody = []byte{}
				}

				// set fixed(enough) length to avoid heap allocation in append
				logs := make([]interface{}, 0, 16)
				logs = append(logs, []interface{}{
					"http.method", req.Method,
					"http.uri", req.RequestURI,
					"http.status", strconv.Itoa(res.Status),
					"http.host", req.Host,
					"request_header", requestHeader,
					"request_body", string(reqBody),
					"response_header", responseHeader,
					"response_body", resBody.String(),
					"user_agent", req.UserAgent(),
					"remote_ip", c.RealIP(),
					"latency", stop.Sub(start).String(),
					"time", start.In(sh).Format(layout),
				}...)
				logs = append(logs, contextFields...)
				p.logger.Infow("http request", logs...)
			}

			if nextErr != nil {
				c.Error(nextErr)
			}
			return nil
		}
	}
}

func extractContextFields(c context.Context) []interface{} {
	// constants.RequestBody -- if this is image request, the length could be too large
	// @todo There could be other keys, such as user id, staff id, etc
	keys := []string{constant.RequestID}
	keyVals := make([]interface{}, 0, len(keys)*2)

	for _, key := range keys {
		keyVals = append(keyVals, key)
		val := c.Value(key)
		if val != nil {
			// @todo set the value length
			keyVals = append(keyVals, val)
		} else {
			keyVals = append(keyVals, "none")
		}
	}
	return keyVals
}
