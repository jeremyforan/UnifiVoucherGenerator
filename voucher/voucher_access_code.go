package voucher

import (
	"fmt"
	"regexp"
	"strings"
)

// AccessCodeLength is the number of digits in a voucher code "12345-67890"
const AccessCodeLength = 10

var accessCodePattern = regexp.MustCompile(`^\d{5}-?\d{5}$`)

// AccessCode is a struct that represents a voucher code
type AccessCode struct {
	firstSet  []int
	secondSet []int
}

// String returns a string representation of a voucher code as it appears on the Unifi
// controller, for example "12345-67890". The zero AccessCode returns an empty string.
func (v AccessCode) String() string {
	if len(v.firstSet) != 5 || len(v.secondSet) != 5 {
		return ""
	}
	var b strings.Builder
	for _, d := range v.firstSet {
		fmt.Fprintf(&b, "%d", d)
	}
	b.WriteByte('-')
	for _, d := range v.secondSet {
		fmt.Fprintf(&b, "%d", d)
	}
	return b.String()
}

// IsZero reports whether the access code has not been set.
func (v AccessCode) IsZero() bool {
	return len(v.firstSet) == 0 && len(v.secondSet) == 0
}

// NewAccessCodeFromString creates a new AccessCode struct from a string. The code must be
// ten digits, with or without the dash the controller displays after the fifth digit.
func NewAccessCodeFromString(voucherCode string) (AccessCode, error) {
	voucherCode = strings.TrimSpace(voucherCode)
	if !accessCodePattern.MatchString(voucherCode) {
		return AccessCode{}, fmt.Errorf("invalid voucher code: want %d digits", AccessCodeLength)
	}

	cleanVoucherCode := strings.ReplaceAll(voucherCode, "-", "")

	a, b := convertStringToIntArray(cleanVoucherCode)

	return AccessCode{
		firstSet:  a,
		secondSet: b,
	}, nil
}
