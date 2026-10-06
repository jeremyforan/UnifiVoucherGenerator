package UnifiVoucherGenerator

import (
	"testing"
)

const sampleVoucherList = `{"meta":{"rc":"ok"},"data":[
  {"duration":1440,"qos_overwrite":false,"note":"abc","code":"1234567890","for_hotspot":true,
   "create_time":1722972000,"quota":1,"site_id":"site1","_id":"66b1","admin_name":"admin","used":0,
   "status":"VALID_ONE","status_expires":0},
  {"duration":60,"qos_overwrite":true,"note":"def","code":"0987654321","for_hotspot":true,
   "create_time":1722972001,"quota":0,"site_id":"site1","_id":"66b2","admin_name":"admin","used":3,
   "status":"VALID_MULTI","status_expires":0}
]}`

func Test_processLoginResponse(t *testing.T) {
	got, err := processLoginResponse(`{"meta":{"rc":"ok"},"data":[]}`)
	if err != nil {
		t.Fatal(err)
	}
	if got.Meta.Rc != "ok" {
		t.Errorf("rc = %q", got.Meta.Rc)
	}

	if _, err := processLoginResponse(`not json`); err == nil {
		t.Errorf("expected error for invalid json")
	}
}

func Test_processNewVoucherRequestResponse(t *testing.T) {
	got, err := processNewVoucherRequestResponse(`{"meta":{"rc":"ok"},"data":[{"create_time":1722972000}]}`)
	if err != nil {
		t.Fatal(err)
	}
	if !got.successful() || len(got.Data) != 1 || got.Data[0].CreateTime != 1722972000 {
		t.Errorf("unexpected response %+v", got)
	}

	if _, err := processNewVoucherRequestResponse(`{`); err == nil {
		t.Errorf("expected error for invalid json")
	}
}

func Test_processVoucherListResponse(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got, err := processVoucherListResponse(sampleVoucherList)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Fatalf("len = %d, want 2", len(got))
		}
		if got[0].Note != "abc" || got[0].Code != "1234567890" || got[0].CreateTime != 1722972000 {
			t.Errorf("first voucher decoded incorrectly: %+v", got[0])
		}
		if got[1].Quota != 0 || got[1].Used != 3 || !got[1].QosOverwrite {
			t.Errorf("second voucher decoded incorrectly: %+v", got[1])
		}
	})

	t.Run("error envelope is an error", func(t *testing.T) {
		_, err := processVoucherListResponse(`{"meta":{"rc":"error","msg":"api.err.LoginRequired"},"data":[]}`)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		if _, err := processVoucherListResponse(`{`); err == nil {
			t.Errorf("expected error")
		}
	})
}

func Test_responseMessage(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{`{"meta":{"rc":"error","msg":"api.err.Invalid"}}`, "api.err.Invalid"},
		{`{"meta":{"rc":"error"}}`, "error"},
		{`<html>`, "<html>"},
		{``, ``},
	}
	for _, tt := range tests {
		if got := responseMessage(tt.body); got != tt.want {
			t.Errorf("responseMessage(%q) = %q, want %q", tt.body, got, tt.want)
		}
	}
}
