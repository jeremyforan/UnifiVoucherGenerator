package UnifiVoucherGenerator

//todo: is there a proper way to do this?
//todo: do I need to consider dashboard versions?

type ExpirationUnit int

const (
	unifiApiLogin          = "/api/login"
	unifiApiSelf           = "/api/self"
	unifiApiLoginReferer   = "/manage/account/login"
	unifiApiCreateVoucher  = "/api/s/default/cmd/hotspot"
	unifiApiVouchers       = "/api/s/default/stat/voucher"
	unifiApiVoucherReferer = "/manage/default/hotspot/vouchers"

	Minutes ExpirationUnit = 1
	Hours   ExpirationUnit = 60
	Days    ExpirationUnit = 1440
)
