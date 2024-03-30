package voucher

import (
	"fmt"
	"log/slog"
	"regexp"
)

// AccessCodeLength is the length of a voucher code "12345-67890"
const AccessCodeLength = 10

// AccessCode is a struct that represents a voucher code
type AccessCode struct {
	firstSet  []int
	secondSet []int
}

//todo: may delete this function as I dont think it is necessary

// String returns a string representation of a voucher code as it appears on the Unifi controller
func (v AccessCode) String() string {
	return fmt.Sprintf("%d%d%d%d%d-%d%d%d%d%d", v.firstSet[0], v.firstSet[1], v.firstSet[2], v.firstSet[3], v.firstSet[4], v.secondSet[0], v.secondSet[1], v.secondSet[2], v.secondSet[3], v.secondSet[4])
}

// NewAccessCodeFromString creates a new AccessCode struct from a string
func NewAccessCodeFromString(voucherCode string) (AccessCode, error) {
	re := regexp.MustCompile(`\d{10}`)

	if !re.MatchString(voucherCode) {
		slog.Error("voucher code is not 10 digits")
		return AccessCode{}, fmt.Errorf("invalid voucher code")
	}
	a, b := convertStringToIntArray(voucherCode)

	return AccessCode{
		firstSet:  a,
		secondSet: b,
	}, nil
}
