package unifi

import (
	"fmt"
	"regexp"
	"strconv"
)

const VoucherCodeLength = 10

type VoucherCode struct {
	firstSet  []int
	secondSet []int
}

func NewVoucherCode() VoucherCode {
	return VoucherCode{
		firstSet:  []int{0, 0, 0, 0, 0},
		secondSet: []int{0, 0, 0, 0, 0},
	}
}

func (v VoucherCode) String() string {
	return fmt.Sprintf("%d%d%d%d%d-%d%d%d%d%d", v.firstSet[0], v.firstSet[1], v.firstSet[2], v.firstSet[3], v.firstSet[4], v.secondSet[0], v.secondSet[1], v.secondSet[2], v.secondSet[3], v.secondSet[4])
}

func NewVoucherFromString(voucherCode string) (VoucherCode, error) {
	// use regex to validate code. it must be 10 digits and only 10 digits.
	if len(voucherCode) != VoucherCodeLength {
		return VoucherCode{}, fmt.Errorf("invalid voucher code")
	}

	re := regexp.MustCompile(`\d{10}`)
	if !re.MatchString(voucherCode) {
		return VoucherCode{}, fmt.Errorf("invalid voucher code")
	}

	buffer := make([]int, VoucherCodeLength)

	for i := 0; i < VoucherCodeLength; i++ {
		c, err := strconv.Atoi(string(voucherCode[i]))
		if err != nil {
			return VoucherCode{}, fmt.Errorf("invalid voucher code")
		}
		buffer[i] = c
	}

	return VoucherCode{
		firstSet:  buffer[:5],
		secondSet: buffer[5:],
	}, nil
}
