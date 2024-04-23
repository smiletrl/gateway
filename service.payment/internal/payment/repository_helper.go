package payment

import (
	"net/http"
	"strconv"
	"time"

	errorpkg "github.com/smiletrl/gateway/pkg/error"
)

// helper methods for repository

// generate a unique transaction id
func (r *repository) generateTransactionID() (transactionID string, err error) {
	// here we use timestamp to generate one unique id, and try it 5 times.
	// if we can't get one unique id for 5 times, we return 503 to frontend to tell user
	// service is unavailable right now and try again later
	retry := 0
	for retry < 5 {
		transactionID = strconv.FormatInt(time.Now().UnixNano(), 10)
		// if id exists already, try again
		if _, ok := r.mockDB.Load(transactionID); ok {
			retry++
		} else {
			return transactionID, nil
		}
	}

	customErr := errorpkg.Error{
		Message: "can not generate an unique transaction id at this moment, please try again later",
		Status:  http.StatusServiceUnavailable,
	}
	return "", &customErr
}
