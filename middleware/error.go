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

package middleware

import (
	"bytes"
	"net/http"
	"time"

	warner "github.com/vicanso/count-warner"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/session"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

// New Error handler
func NewError() elton.Handler {
	// 如果有基于influxdb的统计监控，建议使用influxdb的告警
	// 如果某个IP大量出错，则可能是该IP攻击
	errorWarner := warner.NewWarner(5*time.Minute, 30)
	errorWarner.On(func(ip string, _ int) {
		service.AlarmError("too many errors, ip:" + ip)
	})
	// 定时清理有效数据
	go func() {
		for range time.NewTicker(5 * time.Minute).C {
			errorWarner.ClearExpired()
		}
	}()

	return func(c *elton.Context) error {
		err := c.Next()
		if err == nil {
			return nil
		}
		uri := c.Request.RequestURI
		he, ok := err.(*hes.Error)
		if !ok {
			// 如果不是以http error的形式返回的error则为非主动抛出错误
			he = hes.NewWithError(err)
			he.StatusCode = http.StatusInternalServerError
			he.Exception = true
		} else {
			// 避免修改了原有的error对象
			he = he.Clone()
		}
		if he.StatusCode == 0 {
			he.StatusCode = http.StatusInternalServerError
		}
		if he.Extra == nil {
			he.Extra = make(map[string]interface{})
		}
		account := ""
		tid := util.GetTrackID(c)
		us := session.NewUserSession(c)
		if us != nil && us.IsLogin() {
			account = us.MustGetInfo().Account
		}

		log.Default().Info().
			Str("catgory", "httpError").
			Bool("exception", he.Exception).
			Str("method", c.Request.Method).
			Str("route", c.Route).
			Str("uri", uri).
			Str("error", he.Error()).
			Msg("")

		ip := c.RealIP()
		// 出错则按IP + 1
		errorWarner.Inc(ip, 1)
		sid := util.GetSessionID(c)

		he.Extra["route"] = c.Route
		// 记录用户相关信息
		fields := map[string]interface{}{
			cs.FieldStatus:    he.StatusCode,
			cs.FieldError:     he.Error(),
			cs.FieldURI:       uri,
			cs.FieldException: he.Exception,
			cs.FieldIP:        ip,
			cs.FieldSID:       sid,
			cs.FieldTID:       tid,
		}
		if account != "" {
			fields[cs.FieldAccount] = account
		}
		tags := map[string]string{
			cs.TagMethod: c.Request.Method,
			cs.TagRoute:  c.Route,
		}
		if he.Category != "" {
			tags[cs.TagCategory] = he.Category
		}

		helper.GetInfluxDB().Write(cs.MeasurementHTTPError, tags, fields)
		c.StatusCode = he.StatusCode
		c.BodyBuffer = bytes.NewBuffer(he.ToJSON())
		return nil
	}
}
