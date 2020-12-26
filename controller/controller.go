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
	"net/http"
	"regexp"
	"strconv"

	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/ent/schema"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/elite/middleware"
	"github.com/vicanso/elite/novel"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

type listParams = helper.EntListParams

var (
	errShouldLogin = &hes.Error{
		Message:    "请先登录",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errLoginAlready = &hes.Error{
		Message:    "已是登录状态，请先退出登录",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errForbidden = &hes.Error{
		StatusCode: http.StatusForbidden,
		Message:    "禁止使用该功能",
		Category:   errUserCategory,
	}
)

var (
	logger       = log.Default()
	getEntClient = helper.EntGetClient
	now          = util.NowString

	getUserSession = service.NewUserSession
	// 加载用户session
	loadUserSession = elton.Compose(sessionInterceptor, middleware.NewSession())
	// 判断用户是否登录
	shouldBeLogin = checkLogin
	// 判断用户是否未登录
	shouldBeAnonymous = checkAnonymous
	// 判断用户是否admin权限
	shouldBeAdmin = newCheckRolesMiddleware([]string{
		schema.UserRoleSu,
		schema.UserRoleAdmin,
	})
	// shouldBeSu 判断用户是否su权限
	shouldBeSu = newCheckRolesMiddleware([]string{
		schema.UserRoleSu,
	})

	// 创建新的并发控制中间件
	newConcurrentLimit = middleware.NewConcurrentLimit
	// 创建IP限制中间件
	newIPLimit = middleware.NewIPLimit
	// 创建出错限制中间件
	newErrorLimit = middleware.NewErrorLimit
	// noCacheIfRequestNoCache 请求参数指定no cache，则设置no-cache
	noCacheIfRequestNoCache = middleware.NewNoCacheWithCondition("cacheControl", "no-cache")

	// 图形验证码校验
	captchaValidate = newMagicalCaptchaValidate()
	// 获取influx service
	getInfluxSrv = helper.GetInfluxSrv
	// 文件服务
	fileSrv = &service.FileSrv{}
	// 小说服务
	novelSrv = &novel.Srv{}
	// 图片服务
	imageSrv = &service.ImageSrv{}
	// prof服务
	profSrv = &service.ProfSrv{}
)

func newMagicalCaptchaValidate() elton.Handler {
	magicValue := ""
	if !util.IsProduction() {
		magicValue = "0145"
	}
	return middleware.ValidateCaptcha(magicValue)
}

// isLogin 判断是否登录状态
func isLogin(c *elton.Context) bool {
	us := service.NewUserSession(c)
	return us.IsLogin()
}

// checkLogin 校验是否登录中间件
func checkLogin(c *elton.Context) (err error) {
	if !isLogin(c) {
		err = errShouldLogin
		return
	}
	return c.Next()
}

// checkAnonymous 判断是匿名状态
func checkAnonymous(c *elton.Context) (err error) {
	if isLogin(c) {
		err = errLoginAlready
		return
	}
	return c.Next()
}

// newCheckRolesMiddleware 创建用户角色校验中间件
func newCheckRolesMiddleware(validRoles []string) elton.Handler {
	return func(c *elton.Context) (err error) {
		if !isLogin(c) {
			err = errShouldLogin
			return
		}
		us := service.NewUserSession(c)
		userInfo, err := us.GetInfo()
		if err != nil {
			return
		}
		valid := util.ContainsAny(validRoles, userInfo.Roles)
		if valid {
			return c.Next()
		}
		err = errForbidden
		return
	}
}

// newTracker 初始化用户行为跟踪中间件
func newTracker(action string) elton.Handler {
	return M.NewTracker(M.TrackerConfig{
		Mask: regexp.MustCompile(`(?i)password`),
		OnTrack: func(info *M.TrackerInfo, c *elton.Context) {
			account := ""
			us := service.NewUserSession(c)
			if us != nil && us.IsLogin() {
				account = us.MustGetInfo().Account
			}
			ip := c.RealIP()
			sid := util.GetSessionID(c)
			zapFields := make([]zap.Field, 0, 10)
			zapFields = append(
				zapFields,
				zap.String("action", action),
				zap.String("account", account),
				zap.String("ip", ip),
				zap.String("sid", sid),
				zap.Int("result", info.Result),
			)
			fields := map[string]interface{}{
				"account": account,
				"ip":      ip,
				"sid":     sid,
			}
			if len(info.Query) != 0 {
				zapFields = append(zapFields, zap.Any("query", info.Query))
				fields["query"] = info.Query
			}
			if len(info.Params) != 0 {
				zapFields = append(zapFields, zap.Any("params", info.Params))
				fields["params"] = info.Params
			}
			if len(info.Form) != 0 {
				zapFields = append(zapFields, zap.Any("form", info.Form))
				fields["form"] = info.Form
			}
			if info.Err != nil {
				zapFields = append(zapFields, zap.Error(info.Err))
			}
			logger.Info("tracker", zapFields...)
			getInfluxSrv().Write(cs.MeasurementUserTracker, map[string]string{
				"action": action,
				"result": strconv.Itoa(info.Result),
			}, fields)
		},
	})
}

// getIDFromParams get id form context params
func getIDFromParams(c *elton.Context) (id int, err error) {
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		he := hes.Wrap(err)
		he.Category = "parseInt"
		err = he
		return
	}
	return
}

// sessionInterceptor session的拦截
func sessionInterceptor(c *elton.Context) error {
	message, ok := service.GetSessionInterceptorMessage()
	// 如果有配置拦截信息，则以出错返回
	if ok {
		he := hes.New(message)
		he.Category = "sessionInterceptor"
		return he
	}
	return c.Next()
}
