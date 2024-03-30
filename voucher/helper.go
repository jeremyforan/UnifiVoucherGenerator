package voucher

import (
	uuid "github.com/satori/go.uuid"
	"log/slog"
	"strconv"
)

// convertStringToIntArray helper function to convert a string to an array of integers. This assumes the
// string has already been validated as a proper voucher.
func convertStringToIntArray(s string) ([]int, []int) {
	buffer := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		c, err := strconv.Atoi(string(s[i]))
		if err != nil {
			slog.Error("can not convert char to int", "char", string(s[i]), "error", err, "string", s)
			return []int{}, []int{}
		}
		buffer[i] = c
	}
	return buffer[:5], buffer[5:]
}

// blankVoucher helper function to create a blank voucher with a new UUID
func blankVoucher() Voucher {
	id := uuid.NewV4().String()

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
		ExpireNumber:     "25",
		ExpireUnit:       int(Hours),
		Cmd:              createVoucher,
	}
}
