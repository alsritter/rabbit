package amap

import (
	"context"
	"reflect"
	"testing"

	"alsritter.icu/rabbit-template/internal/pkg/httpclient"

	"go.opentelemetry.io/otel"
)

func TestAMap_GetDirection(t *testing.T) {
	type fields struct {
		apiDomain string
		key       string
	}
	type args struct {
		ctx         context.Context
		origin      string
		destination string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult *DirectionResp
		wantErr    bool
	}{
		{
			name: "",
			fields: fields{
				apiDomain: "https://restapi.amap.com/v3",
				key:       "xxxxxxxxxxxxxxxxxxxxxxxx",
			},
			args: args{
				ctx:         context.Background(),
				origin:      "116.481028,39.989643",
				destination: "116.434446,39.90816",
			},
			wantResult: &DirectionResp{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AMap{
				apiDomain: tt.fields.apiDomain,
				key:       tt.fields.key,
				client:    httpclient.New(otel.Tracer("rabbit-local")),
			}
			gotResult, err := a.GetDirection(tt.args.ctx, tt.args.origin, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("AMap.GetDirection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("AMap.GetDirection() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestAMap_GetDistrict(t *testing.T) {
	type fields struct {
		apiDomain string
		key       string
	}
	type args struct {
		ctx         context.Context
		keywords    string
		subdistrict int32
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult *DistrictResp
		wantErr    bool
	}{
		{
			name: "",
			fields: fields{
				apiDomain: "https://restapi.amap.com/v3",
				key:       "xxxxxxxxxxxxxxxxxxxxxxxx",
			},
			args: args{
				ctx:         context.Background(),
				keywords:    "440000",
				subdistrict: 1,
			},
			wantResult: &DistrictResp{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AMap{
				apiDomain: tt.fields.apiDomain,
				key:       tt.fields.key,
				client:    httpclient.New(otel.Tracer("rabbit-local")),
			}
			gotResult, err := a.GetDistrict(tt.args.ctx, tt.args.keywords, tt.args.subdistrict)
			if (err != nil) != tt.wantErr {
				t.Errorf("AMap.GetDistrict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("AMap.GetDistrict() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
