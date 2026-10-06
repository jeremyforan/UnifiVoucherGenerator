package voucher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Voucher is a struct that holds the information needed to create a new voucher.
type Voucher struct {
	// Id is the identifier sent to the controller in the voucher's note field. It is used
	// to find the voucher again after it has been created.
	Id        string
	published bool
	data      Data
	// AC is the access code assigned by the controller. It is populated by Client.AddVoucher.
	AC AccessCode
}

// NewDefaultVoucher creates a new Voucher struct to be used to create a new voucher.
// The default voucher is single use and expires 24 hours after first use.
func NewDefaultVoucher() *Voucher {
	id := uuid.NewString()
	return &Voucher{
		Id:        id,
		published: false,
		data: Data{
			Note:             id,
			Quota:            int(vSingleUse),
			NumberOfVouchers: 1,
			ExpireNumber:     fmt.Sprintf("%d", defaultExpireHours),
			ExpireUnit:       int(Hours),
			Cmd:              createVoucher,
		},
	}
}

// NewSingleUseVoucher creates a new Single Use Voucher.
func NewSingleUseVoucher() *Voucher {
	v := blankVoucher()

	v.data.Quota = int(vSingleUse)

	return &v
}

// NewMultiUseVoucher creates a new Multi Use Voucher that can be redeemed quota times.
// A quota of 0 produces an unlimited use voucher, matching the controller's semantics.
func NewMultiUseVoucher(quota int) *Voucher {
	v := blankVoucher()

	v.data.Quota = quota

	return &v
}

// NewUnlimitedUseVoucher creates a new Unlimited use Voucher.
func NewUnlimitedUseVoucher() *Voucher {
	v := blankVoucher()

	v.data.Quota = int(vUnlimited)

	return &v
}

// String Stringer interface implementation. It is safe to call on a nil Voucher.
func (v *Voucher) String() string {
	if v == nil {
		return ""
	}
	return v.data.String()
}

// HttpPayload returns the Voucher struct as a bytes.Reader to be used in as a http request body.
func (v *Voucher) HttpPayload() *bytes.Reader {
	marshalled, err := json.Marshal(v.data)
	if err != nil {
		slog.Error("Error marshaling to JSON", "error", err)
		return nil
	}

	return bytes.NewReader(marshalled)
}

// AccessCode returns the AccessCode for the voucher. This is the 10-digit code that Guest can use to access the network.
// It is the zero AccessCode until the voucher has been added successfully.
func (v *Voucher) AccessCode() AccessCode {
	return v.AC
}

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

// SetDataLimitMB sets the `Data Limit` in MB. If not set, the default is unlimited.
func (v *Voucher) SetDataLimitMB(limit int) {
	v.data.Bytes = fmt.Sprintf("%d", limit)
}

// SetExpire sets the `Expire Number` and `Expire Unit` for the voucher.
func (v *Voucher) SetExpire(expiration int, unit ExpireUnit) {
	v.data.ExpireNumber = fmt.Sprintf("%d", expiration)
	v.data.ExpireUnit = int(unit)
}

// SetId sets the identifier for the voucher. It is sent to the controller as the voucher's
// `Note` and is used to look the voucher up after creation, so both are kept in sync.
func (v *Voucher) SetId(id string) {
	v.Id = id
	v.data.Note = id
}

// TODO: The publishing of vouchers doesnt incorporate the number of vouchers to create. The code that
// would need to be reworked before this can be implemented.
// SetAmountOfVouchers sets the number of vouchers to create.
//func (v *Voucher) SetAmountOfVouchers(amount int) {
//	v.data.NumberOfVouchers = amount
//}
