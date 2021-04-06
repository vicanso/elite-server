// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package location

import (
	"context"

	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/request"
	"github.com/vicanso/go-axios"
)

var ins = newInstance()

func newInstance() *axios.Instance {
	locationConfig := config.GetLocationConfig()
	service := "location"
	return request.NewHTTP(service, locationConfig.BaseURL, locationConfig.Timeout)
}

// 相关的URL
const (
	locationURL = "/ip-locations/json/:ip"
)

// Location location
type Location struct {
	IP string `json:"ip"`
	// Country 国家
	Country string `json:"country"`
	// Province 省
	Province string `json:"province"`
	// City 市
	City string `json:"city"`
	// ISP 网络接入商
	ISP string `json:"isp"`
}

// GetByIP get location by ip
func GetByIP(ctx context.Context, ip string) (lo Location, err error) {
	conf := &axios.Config{
		URL: locationURL,
		Params: map[string]string{
			"ip": ip,
		},
		Context: ctx,
	}

	lo = Location{}
	err = ins.EnhanceRequest(&lo, conf)
	if err != nil {
		return
	}
	return
}
