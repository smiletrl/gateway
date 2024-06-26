package test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/smiletrl/gateway/pkg/constant"
)

// APITestCase represents the data needed to describe an API test case.
type APITestCase struct {
	Name    string
	Method  string
	URL     string
	Context map[constant.ContextString]interface{}
	Body    string
	// ignore this (dynamic) field in result comparement, something like auto increased id generated by db.
	// since it's dynamic, we can't predefine the result and use it for comparement.
	IgnoreField string
	// execute this query before call api request, something like prepare data for this test case
	Prequery string
	// after this test case, execute certain (db) query, something like reset db.
	AfterQuery     string
	ExpectStatus   int
	ExpectResponse string
}

// Endpoint tests an HTTP endpoint using the given APITestCase spec.
func Endpoint(t *testing.T, e *echo.Echo, tc APITestCase) (ActualResponse string) {
	var res *httptest.ResponseRecorder
	t.Run(tc.Name, func(t *testing.T) {
		ctx := context.Background()
		if tc.Context != nil {
			for key, val := range tc.Context {
				ctx = context.WithValue(ctx, key, val)
			}
		}
		req, _ := http.NewRequestWithContext(ctx, tc.Method, tc.URL, bytes.NewBufferString(tc.Body))

		res = httptest.NewRecorder()
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		e.ServeHTTP(res, req)

		require.Equal(t, tc.ExpectStatus, res.Code, "status mismatch")
		pattern := strings.Trim(tc.ExpectResponse, "*")

		// If returned json string is too long, and you think this api is correct, print the result here
		// for debug purpose
		// prettyJson := pretty.PrettyOptions(res.Body.Bytes(), &pretty.Options{Indent: "    "})
		// fmt.Printf("API test res: (\n\n%s\n)\n", prettyJson)

		// compare string firstly
		if pattern != tc.ExpectResponse {
			require.Contains(t, res.Body.String(), pattern, "response string mismatch")
		} else {
			// compare json
			// get "id" , from string, and remove this part. This is because id is dynamic value.
			body := res.Body.String()
			if tc.IgnoreField != "" {
				tc.ExpectResponse = RemoveDynamicField(tc.ExpectResponse, tc.IgnoreField)
				body = RemoveDynamicField(body, tc.IgnoreField)
			}
			require.JSONEq(t, tc.ExpectResponse, body, "response json mismatch")
		}
	})

	// return true api result. Due to ignored field, the expect response will be different than the real response
	return res.Body.String()
}

func RemoveDynamicField(s, fields string) string {
	fieldsArray := strings.Split(fields, ",")
	for _, field := range fieldsArray {
		s = removeDynamicField(s, field)
	}

	return s
}

// This func is temporary fix
func removeDynamicField(s, field string) string {
	// string s might look like `{"id": 1232, "name": "jack"}`, and we need to return
	// `{"name": "jack"}`
	idIndex := strings.Index(s, fmt.Sprintf(`"%s"`, field))
	if idIndex == -1 {
		return s
	}
	subString := s[idIndex:]
	commaIndex := strings.Index(subString, ",")
	if commaIndex == -1 {
		// if comma is not found, search `}`
		closeIndex := strings.Index(subString, `}`)
		if closeIndex == -1 {
			return s
		}
		idString := subString[:closeIndex]
		return strings.Replace(s, idString, "", 1)
	}

	idString := subString[:commaIndex+1]

	newIDString := strings.Replace(s, idString, "", 1)

	// repeat until no field is found
	return removeDynamicField(newIDString, field)
}
