package voucher

import (
	"fmt"
	"log/slog"
	"strings"
)

//todo: these should be explained in the documentation. Called out by the value they represent.

// Data is a struct that holds the parameters that can be selected online.
type Data struct {
	Quota            int    `json:"quota"`
	Note             string `json:"note"`
	NumberOfVouchers int    `json:"n"`
	ExpireNumber     int    `json:"expire_number"`
	ExpireUnit       int    `json:"expire_unit"`
	Cmd              string `json:"cmd"`
	Up               int    `json:"up"`
	Down             int    `json:"down"`
	Bytes            int    `json:"bytes"`
}

// String returns the NewVoucherRequestPayload struct as a string.
func (v *Data) String() string {
	return fmt.Sprintf(`{"quota":%d,"note":"%s","n":%d,"expire_number":%d,"expire_unit":%d,"cmd":"%s"}`, v.Quota, v.Note, v.NumberOfVouchers, v.ExpireNumber, v.ExpireUnit, v.Cmd)
}

// HttpPayload returns the NewVoucherRequestPayload struct as a strings.Reader to be used in
// as a http request body.
func (v *Data) HttpPayload() *strings.Reader {
	return strings.NewReader(v.String())
}

// LogGroup returns the NewVoucherRequestPayload struct as a slog.Value for logging.
func (v *Data) LogGroup() slog.Value {
	return slog.GroupValue(
		slog.Int("quota", v.Quota),
		slog.String("note", v.Note),
		slog.Int("amount", v.NumberOfVouchers),
		slog.Int("expire_number", v.ExpireNumber),
		slog.Int("expire_unit", v.ExpireUnit),
		slog.String("cmd", v.Cmd),
	)
}
