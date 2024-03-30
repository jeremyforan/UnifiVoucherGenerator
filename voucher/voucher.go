package voucher

import (
	"github.com/satori/go.uuid"
)

//todo: I think this should tweak the names a little bit.

// Voucher is a struct that holds the information needed to create a new voucher.
type Voucher struct {
	Id   string
	data Data
}

// NewDefaultVoucher creates a new Voucher struct to be used to create a new voucher.
func NewDefaultVoucher() Voucher {
	id := uuid.NewV4().String()
	return Voucher{
		Id: id,
		data: Data{
			Note:             id,
			Quota:            1,
			NumberOfVouchers: 1,
			ExpireNumber:     24,
			ExpireUnit:       int(vHours),
		},
	}
}

func NewSingleUseVoucher() Voucher {
	return Voucher{
		Id: uuid.NewV4().String(),
		data: Data{
			Quota: int(vSingleUse),
		},
	}
}

func NewMultiUseVoucher(quota int) Voucher {
	return Voucher{
		Id: uuid.NewV4().String(),
		data: Data{
			Quota: quota,
		},
	}
}

func NewUnlimitedUseVoucher() Voucher {
	return Voucher{
		Id: uuid.NewV4().String(),
		data: Data{
			Quota: int(vUnlimited),
		},
	}
}
