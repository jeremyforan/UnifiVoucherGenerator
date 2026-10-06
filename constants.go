package UnifiVoucherGenerator

// Paths on the UniFi Network Application. These are the legacy controller
// paths; the site is currently fixed to "default".
const (
	unifiApiLogin          = "/api/login"
	unifiApiLoginReferer   = "/manage/account/login"
	unifiApiCreateVoucher  = "/api/s/default/cmd/hotspot"
	unifiApiVouchers       = "/api/s/default/stat/voucher"
	unifiApiVoucherReferer = "/manage/default/hotspot/vouchers"
)
