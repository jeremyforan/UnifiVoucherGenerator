package voucher

import (
	"fmt"

	"github.com/google/uuid"
)

// convertStringToIntArray helper function to convert a string of digits to two arrays of
// integers. This assumes the string has already been validated as a proper voucher code.
func convertStringToIntArray(s string) ([]int, []int) {
	buffer := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return []int{}, []int{}
		}
		buffer[i] = int(s[i] - '0')
	}
	return buffer[:5], buffer[5:]
}

// blankVoucher helper function to create a blank voucher with a new UUID
func blankVoucher() Voucher {
	id := uuid.NewString()

	d := blankVoucherData()

	d.Note = id

	return Voucher{
		Id:        id,
		published: false,
		data:      d,
	}
}

// blankVoucherData helper function to create a blank voucher data struct. This sets the default
// expiry to 24 hours.
func blankVoucherData() Data {
	return Data{
		Note:             "",
		Quota:            0,
		NumberOfVouchers: 1,
		ExpireNumber:     fmt.Sprintf("%d", defaultExpireHours),
		ExpireUnit:       int(Hours),
		Cmd:              createVoucher,
	}
}
