package constant

// logger
const (
	Logger = "logger"
)

type ContextString string

// request
const (
	// request id
	RequestID ContextString = "request_id"

	// request bdy
	RequestBody         string = "request_body"
	HttpRequestIDHeader string = "x-request-id"
)
