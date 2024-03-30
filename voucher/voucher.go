package voucher

import (
	"bytes"
	"encoding/json"
	"github.com/satori/go.uuid"
	"log/slog"
)

//todo: I think this should tweak the names a little bit.

// Voucher is a struct that holds the information needed to create a new voucher.
type Voucher struct {
	Id        string
	published bool
	data      Data
	AC        AccessCode
}

// NewDefaultVoucher creates a new Voucher struct to be used to create a new voucher.
func NewDefaultVoucher() *Voucher {
	id := uuid.NewV4().String()
	return &Voucher{
		Id:        id,
		published: false,
		data: Data{
			Note:             id,
			Quota:            1,
			NumberOfVouchers: 1,
			ExpireNumber:     24,
			ExpireUnit:       int(vHours),
			Cmd:              createVoucher,
		},
	}
}

// NewSingleUseVoucher creates a new Single Use Voucher.
func NewSingleUseVoucher() *Voucher {
	v := blankVoucher()
	v.data = Data{
		Note:  v.Id,
		Quota: int(vSingleUse),
		Cmd:   createVoucher,
	}
	return &v
}

// NewMultiUseVoucher creates a new Multi Use Voucher.
func NewMultiUseVoucher(quota int) *Voucher {
	v := blankVoucher()
	v.data = Data{
		Note:  v.Id,
		Quota: quota,
		Cmd:   createVoucher,
	}
	return &v
}

// NewUnlimitedUseVoucher creates a new Unlimited use Voucher.
func NewUnlimitedUseVoucher() *Voucher {
	v := blankVoucher()
	v.data = Data{
		Note:  v.Id,
		Quota: int(vUnlimited),
		Cmd:   createVoucher,
	}
	return &v
}

func (v *Voucher) String() string {
	return v.data.String()
}

func (v *Voucher) HttpPayload() *bytes.Reader {
	marshalled, err := json.Marshal(v.data)
	if err != nil {
		slog.Error("Error marshaling to JSON", "error", err)
		return nil
	}

	return bytes.NewReader(marshalled)
}

// AccessCode returns the AccessCode for the voucher. This is the 10-digit code that Guest can use to access the network.
func (v *Voucher) AccessCode() AccessCode {
	return v.AC
}

// todo: move and maybe rename this function.

// PublishedSuccesfully sets the voucher as published.
func (v *Voucher) PublishedSuccesfully() {
	v.published = true
}

// Published returns true if the voucher has been published and available on the UniFi Network Application.
func (v *Voucher) Published() bool {
	return v.published
}

// SetDownloadLimitMbps sets the `Download Limit` in Mbps. If not set, the default is unlimited.
func (v *Voucher) SetDownloadLimitMbps(limit int) {
	v.data.Down = limit * vMbps
}

// SetUploadLimitMbps sets the `Upload Limit` in Mbps. If not set, the default is unlimited.
func (v *Voucher) SetUploadLimitMbps(limit int) {
	v.data.Up = limit * vMbps
}

// SetBandwidthLimitMB sets the `Data Limit` in MB. If not set, the default is unlimited.
func (v *Voucher) SetBandwidthLimitMB(limit int) {
	v.data.Bytes = limit
}

// SetExpireInHours sets the `Expire Number` and `Expire Unit` for the voucher.
func (v *Voucher) SetExpire(expiration int, unit ExpirationUnit) {
	v.data.ExpireNumber = expiration
	v.data.ExpireUnit = int(unit)
}
