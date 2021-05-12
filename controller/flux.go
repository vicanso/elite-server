// Copyright 2021 tree xie
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

// flux查询influxdb相关数据

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
)

type fluxCtrl struct{}

// 参数相关定义
type (
	// fluxListParams flux查询参数
	fluxListParams struct {
		Measurement string    `json:"measurement"`
		Begin       time.Time `json:"begin" validate:"required"`
		End         time.Time `json:"end" validate:"required"`
		Account     string    `json:"account" validate:"omitempty,xUserAccount"`
		Limit       string    `json:"limit" validate:"required,xLargerLimit"`
		Exception   string    `json:"exception" validate:"omitempty,xBoolean"`
		// 用户行为类型筛选
		Action      string `json:"action" validate:"omitempty,xTag"`
		Result      string `json:"result" validate:"omitempty,xTag"`
		Category    string `json:"category" validate:"omitempty,xTag"`
		ErrCategory string `json:"errCategory" validate:"omitempty,xTag"`
		Route       string `json:"route" validate:"omitempty,xTag"`
		Service     string `json:"service" validate:"omitempty,xTag"`
	}
	// flux tags/fields查询参数
	fluxListTagOrFieldParams struct {
		Measurement string `json:"measurement" validate:"required,xMeasurement"`
	}
	// fluxListTagValuesParams flux tag values查询参数
	fluxListTagValuesParams struct {
		Measurement string `json:"measurement" validate:"required,xMeasurement"`
		Tag         string `json:"tag" validate:"required,xTag"`
	}
)

// 响应相关定义

func init() {
	sessionConfig = config.GetSessionConfig()
	g := router.NewGroup("/fluxes", loadUserSession)

	ctrl := fluxCtrl{}
	// 查询用户tracker
	g.GET(
		"/v1/trackers",
		shouldBeAdmin,
		ctrl.listTracker,
	)
	// 查询http出错
	g.GET(
		"/v1/http-errors",
		shouldBeAdmin,
		ctrl.listHTTPError,
	)
	// 获取用户action
	g.GET(
		"/v1/actions",
		shouldBeAdmin,
		ctrl.listAction,
	)
	// 获取request相关调用统计
	g.GET(
		"/v1/requests",
		shouldBeAdmin,
		ctrl.listRequest,
	)

	// 获取tag
	// 不校验登录状态，无敏感信息
	g.GET(
		"/v1/tags/{measurement}",
		ctrl.listTag,
	)
	// 获取tag的取值列表
	// 不校验登录状态，无敏感信息
	g.GET(
		"/v1/tag-values/{measurement}/{tag}",
		ctrl.listTagValue,
	)
	// 获取field
	g.GET(
		"/v1/fields/{measurement}",
		ctrl.ListField,
	)
}

// Query get flux query string
func (params *fluxListParams) Query() string {
	start := util.FormatTime(params.Begin.UTC())
	stop := util.FormatTime(params.End.UTC())
	query := fmt.Sprintf(`|> range(start: %s, stop: %s)
|> filter(fn: (r) => r["_measurement"] == "%s")
`,
		start,
		stop,
		params.Measurement,
	)
	addStrQuery := func(key, value string) {
		query += fmt.Sprintf(`|> filter(fn: (r) => r.%s == "%s")
`, key, value)
	}
	addQuery := func(key string, value interface{}) {
		query += fmt.Sprintf(`|> filter(fn: (r) => r.%s == %s)
`, key, value)
	}

	// TODO 根据measurement判断是tag还是field

	// tag 的筛选
	// 用户行为类型
	if params.Action != "" {
		addStrQuery(cs.TagAction, params.Action)
	}
	// 结果
	if params.Result != "" {
		addStrQuery(cs.TagResult, params.Result)
	}
	// 分类
	if params.Category != "" {
		addStrQuery(cs.TagCategory, params.Category)
	}
	// service
	if params.Service != "" {
		addStrQuery(cs.TagService, params.Service)
	}
	// route
	if params.Route != "" {
		addStrQuery("route", params.Route)
	}

	// 筛选完成之后执行pivot
	query += fmt.Sprintf(`|> sort(columns:["_time"], desc: true)
|> limit(n:%s)
|> pivot(
	rowKey:["_time"],
	columnKey: ["_field"],
	valueColumn: "_value"
)
`, params.Limit)

	// field 的筛选
	// 账号
	if params.Account != "" {
		addStrQuery(cs.FieldAccount, params.Account)
	}
	// 异常
	if params.Exception != "" {
		addQuery(cs.FieldException, params.Exception)
	}
	// 出错类型
	if params.ErrCategory != "" {
		addStrQuery(cs.FieldErrCategory, params.ErrCategory)
	}

	return query
}

func (params *fluxListParams) Do(ctx context.Context) (items []map[string]interface{}, err error) {
	items, err = GetInfluxSrv().Query(ctx, params.Query())
	if err != nil {
		return
	}
	// 清除不需要字段
	for _, item := range items {
		delete(item, "_measurement")
		delete(item, "_start")
		delete(item, "_stop")
		delete(item, "table")
	}
	return
}

// listTag returns the tags of measurement
func (ctrl fluxCtrl) listTag(c *elton.Context) (err error) {
	params := fluxListTagOrFieldParams{}
	err = validate.Do(&params, c.Params.ToMap())
	if err != nil {
		return
	}
	tags, err := GetInfluxSrv().ListTag(c.Context(), params.Measurement)
	if err != nil {
		return
	}
	c.CacheMaxAge(time.Minute)
	c.Body = map[string][]string{
		"tags": tags,
	}
	return
}

// ListField return the fields of measurement
func (ctrl fluxCtrl) ListField(c *elton.Context) (err error) {
	params := fluxListTagOrFieldParams{}
	err = validate.Do(&params, c.Params.ToMap())
	if err != nil {
		return
	}
	fields, err := GetInfluxSrv().ListField(c.Context(), params.Measurement, "-30d")
	if err != nil {
		return
	}
	c.CacheMaxAge(time.Minute)
	c.Body = map[string][]string{
		"fields": fields,
	}
	return
}

// listValue get the values of tag
func (ctrl fluxCtrl) listTagValue(c *elton.Context) (err error) {
	params := fluxListTagValuesParams{}
	err = validate.Do(&params, c.Params.ToMap())
	if err != nil {
		return
	}
	values, err := GetInfluxSrv().ListTagValue(c.Context(), params.Measurement, params.Tag)
	if err != nil {
		return
	}
	c.CacheMaxAge(time.Minute)
	c.Body = map[string][]string{
		"values": values,
	}
	return
}

func (ctrl fluxCtrl) list(c *elton.Context, measurement, responseKey string) (err error) {
	params := fluxListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	params.Measurement = measurement
	result, err := params.Do(c.Context())
	if err != nil {
		return
	}

	fromBucket := fmt.Sprintf(`from(bucket: "%s")
`, GetInfluxDB().GetBucket())
	c.Body = map[string]interface{}{
		responseKey: result,
		"count":     len(result),
		"flux":      fromBucket + params.Query(),
	}
	return
}

// listHTTPError list http error
func (ctrl fluxCtrl) listHTTPError(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementHTTPError, "httpErrors")
}

// listTracker list user tracker
func (ctrl fluxCtrl) listTracker(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementUserTracker, "trackers")
}

// listAction list user action
func (ctrl fluxCtrl) listAction(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementUserAction, "actions")
}

// listRequest list request
func (ctrl fluxCtrl) listRequest(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementHTTPRequest, "requests")
}
