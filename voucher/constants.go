package voucher

// ExpireUnit is the unit of time that the voucher will expire in
type ExpireUnit int

// UsageType is the type of voucher that will be created. There is single, multi and unlimited use.
type UsageType int

const (
	vSingleUse UsageType = 1
	vUnlimited UsageType = 0
	//	vMultiUse is set in constructor

	Minutes ExpireUnit = 1
	Hours   ExpireUnit = 60
	Days    ExpireUnit = 1440

	createVoucher string = "create-voucher"

	vMbps = 1000
)
