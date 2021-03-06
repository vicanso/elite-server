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

package controller

import (
	"time"

	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elton"
)

type applicationCtrl struct{}

func init() {
	g := router.NewGroup("/applications")
	ctrl := applicationCtrl{}

	g.GET("/v1/setting", ctrl.getSetting)

}

func (*applicationCtrl) getSetting(c *elton.Context) (err error) {
	settings, err := configurationSrv.ListApplicationSetting(c.Context())
	if err != nil {
		return
	}
	if len(settings) == 0 {
		c.NoContent()
		return
	}
	setting, err := settings.First(util.GetAppVersion(c))
	if err != nil {
		return
	}
	if setting == nil {
		c.NoContent()
		return
	}
	c.PrivateCacheMaxAge(5 * time.Minute)
	c.Body = setting
	return
}
