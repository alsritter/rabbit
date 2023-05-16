package amap

import (
	"context"
	"fmt"

	"alsritter.icu/rabbit-template/internal/pkg/httpclient"
)

type AMap struct {
	apiDomain string
	key       string
	client    *httpclient.Client
}

func New(key string, client *httpclient.Client) *AMap {
	return &AMap{
		apiDomain: "https://restapi.amap.com/v3",
		key:       key,
		client:    client,
	}
}

const (
	directionPath = "direction/driving" // 驾车路径规划
	districtPath  = "config/district"   // 行政区域查询
)

// https://lbs.amap.com/api/webservice/guide/api/direction/#driving
func (a *AMap) GetDirection(ctx context.Context, origin, destination string) (result *DirectionResp, err error) {
	var api = fmt.Sprintf("%s/%s", a.apiDomain, directionPath)
	var req = a.client.NewRequest(httpclient.Get, api)
	req.AddParam("key", a.key)
	req.AddParam("origin", origin)
	req.AddParam("destination", destination)
	req.AddParam("extensions", "base")

	var rsp = req.Exec(ctx)
	if err := rsp.Error(); err != nil {
		return nil, err
	}

	if err = rsp.Unmarshal(&result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("get direction failed, status: %s, info: %s", result.Status, result.Info)
	}

	return result, nil
}

// https://lbs.amap.com/api/webservice/guide/api/district
func (a *AMap) GetDistrict(ctx context.Context, keywords string, subdistrict int32) (result *DistrictResp, err error) {
	var api = fmt.Sprintf("%s/%s", a.apiDomain, districtPath)
	var req = a.client.NewRequest(httpclient.Get, api)
	req.AddParam("key", a.key)
	req.AddParam("keywords", keywords)
	req.AddParam("subdistrict", fmt.Sprintf("%d", subdistrict))

	var rsp = req.Exec(ctx)
	if err := rsp.Error(); err != nil {
		return nil, err
	}

	if err = rsp.Unmarshal(&result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("get district failed, status: %s, info: %s", result.Status, result.Info)
	}

	return result, nil
}
