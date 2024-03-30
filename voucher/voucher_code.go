package voucher

import (
	"fmt"
	"log/slog"
	"regexp"
)

// VoucherCodeLength is the length of a voucher code "12345-67890"
const VoucherCodeLength = 10

// Code is a struct that represents a voucher code
type Code struct {
	firstSet  []int
	secondSet []int
}

//todo: may delete this function as I dont think it is necessary

func NewVoucherCode() Code {
	return Code{
		firstSet:  []int{0, 0, 0, 0, 0},
		secondSet: []int{0, 0, 0, 0, 0},
	}
}

// String returns a string representation of a voucher code as it appears on the Unifi controller
func (v Code) String() string {
	return fmt.Sprintf("%d%d%d%d%d-%d%d%d%d%d", v.firstSet[0], v.firstSet[1], v.firstSet[2], v.firstSet[3], v.firstSet[4], v.secondSet[0], v.secondSet[1], v.secondSet[2], v.secondSet[3], v.secondSet[4])
}

// NewVoucherFromString creates a new Code struct from a string
func NewVoucherFromString(voucherCode string) (Code, error) {
	re := regexp.MustCompile(`\d{10}`)

	if !re.MatchString(voucherCode) {
		slog.Error("voucher code is not 10 digits")
		return Code{}, fmt.Errorf("invalid voucher code")
	}
	a, b := convertStringToIntArray(voucherCode)

	return Code{
		firstSet:  a,
		secondSet: b,
	}, nil
}
