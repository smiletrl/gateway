package payment

import (
	"strconv"
)

// helper methods for service
func (s *service) isAcquirerApproved(card string) bool {
	rs := []rune(card)
	// get last rune char
	lastDigit := string(rs[len(rs)-1])
	last, _ := strconv.Atoi(lastDigit)
	// approve if last digit of card is even
	return last%2 == 0
}
