package UnifiVoucherGenerator

import (
	"reflect"
	"testing"
)

func Test_processLoginResponse(t *testing.T) {
	type args struct {
		body string
	}
	var tests []struct {
		name    string
		args    args
		want    *LoginResponse
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processLoginResponse(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("processLoginResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processLoginResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processNewVoucherRequestResponse(t *testing.T) {
	type args struct {
		body string
	}
	var tests []struct {
		name    string
		args    args
		want    RequestNewVoucherResponse
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processNewVoucherRequestResponse(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("processNewVoucherRequestResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processNewVoucherRequestResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_processResponse(t *testing.T) {
//	type args struct {
//		body string
//	}
//	type testCase[T any] struct {
//		name    string
//		args    args
//		want    *T
//		wantErr bool
//	}
//	var tests []testCase[ /* TODO: Insert concrete types here */ ]
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := processResponse(tt.args.body)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("processResponse() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("processResponse() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func Test_processVoucherListResponse(t *testing.T) {
	type args struct {
		body string
	}
	var tests []struct {
		name    string
		args    args
		want    UnifiVouchers
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processVoucherListResponse(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("processVoucherListResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processVoucherListResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
