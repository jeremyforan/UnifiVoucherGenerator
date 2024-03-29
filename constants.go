package unifi

const (
	unifiApiBaseUrl = "https://unifi.jeremyforan.com"

	unifiApiLogin          = unifiApiBaseUrl + "/api/login"
	unifiApiLoginReferer   = unifiApiBaseUrl + "/manage/account/login"
	unifiApiCreateVoucher  = unifiApiBaseUrl + "/api/s/default/cmd/hotspot"
	unifiApiVouchers       = unifiApiBaseUrl + "/api/s/default/stat/voucher"
	unifiApiVoucherReferer = unifiApiBaseUrl + "/manage/default/hotspot/vouchers"
)
