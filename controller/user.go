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

// 用户相关的一些路由处理

package controller

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqljson"
	"github.com/tidwall/gjson"
	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/predicate"
	"github.com/vicanso/elite/ent/schema"
	"github.com/vicanso/elite/ent/user"
	"github.com/vicanso/elite/ent/userlogin"
	"github.com/vicanso/elite/middleware"
	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

type (
	userCtrl struct{}

	// userInfoResp 用户信息响应
	userInfoResp struct {
		Date string `json:"date,omitempty"`
		service.UserSessionInfo
	}

	// userListResp 用户列表响应
	userListResp struct {
		Users []*ent.User `json:"users,omitempty"`
		Count int         `json:"count,omitempty"`
	}
	// userRoleListResp 用户角色列表响应
	userRoleListResp struct {
		UserRoles []*schema.UserRoleInfo `json:"userRoles,omitempty"`
	}
	// userLoginListResp 用户登录列表响应
	userLoginListResp struct {
		UserLogins []*ent.UserLogin `json:"userLogins,omitempty"`
		Count      int              `json:"count,omitempty"`
	}

	// userListParams 用户查询参数
	userListParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		Role    string `json:"role,omitempty" validate:"omitempty,xUserRole"`
		Group   string `json:"group,omitempty" validate:"omitempty,xUserGroup"`
		Status  string `json:"status,omitempty" validate:"omitempty,xStatus"`
	}

	// userLoginListParams 用户登录查询
	userLoginListParams struct {
		listParams

		Begin   time.Time `json:"begin,omitempty"`
		End     time.Time `json:"end,omitempty"`
		Account string    `json:"account,omitempty" validate:"omitempty,xUserAccount"`
	}

	// userRegisterLoginParams 注册与登录参数
	userRegisterLoginParams struct {
		// 账户
		Account string `json:"account,omitempty" validate:"required,xUserAccount"`
		// 密码，密码为sha256后的加密串
		Password string `json:"password,omitempty" validate:"required,xUserPassword"`
	}

	// userUpdateMeParams 用户信息更新参数
	userUpdateMeParams struct {
		Name        string `json:"name,omitempty" validate:"omitempty,xUserName"`
		Email       string `json:"email,omitempty" validate:"omitempty,xUserEmail"`
		Password    string `json:"password,omitempty" validate:"omitempty,xUserPassword"`
		NewPassword string `json:"newPassword,omitempty" validate:"omitempty,xUserPassword"`
	}
)

var (
	// session配置信息
	sessionConfig config.SessionConfig
)

const (
	errUserCategory = "user"
)

var (
	errLoginTokenNil = &hes.Error{
		Message:    "登录令牌不能为空",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errAccountOrPasswordInvalid = &hes.Error{
		Message:    "账户或者密码错误",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errOldPasswordWrong = &hes.Error{
		Message:    "旧密码错误，请重新输入",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errUserStatusInvalid = &hes.Error{
		Message:    "该账户不允许登录",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errUserAccountExists = &hes.Error{
		Message:    "该账户已注册",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
)

func init() {
	sessionConfig = config.GetSessionConfig()
	prefix := "/users"
	g := router.NewGroup(prefix, loadUserSession)
	noneSessionGroup := router.NewGroup(prefix)

	ctrl := userCtrl{}

	// 获取用户列表
	g.GET(
		"/v1",
		shouldBeAdmin,
		ctrl.list,
	)

	// 获取用户信息
	g.GET(
		"/v1/{id}",
		shouldBeAdmin,
		ctrl.findByID,
	)

	// 获取登录token
	g.GET(
		"/v1/me/login",
		shouldBeAnonymous,
		ctrl.getLoginToken,
	)

	// 获取用户信息
	g.GET(
		"/v1/me",
		ctrl.me,
	)

	// 用户注册
	g.POST(
		"/v1/me",
		middleware.WaitFor(time.Second, true),
		newTracker(cs.ActionRegister),
		captchaValidate,
		// 限制相同IP在60秒之内只能调用5次
		newIPLimit(5, 60*time.Second, cs.ActionRegister),
		shouldBeAnonymous,
		ctrl.register,
	)

	// 用户登录
	g.POST(
		"/v1/me/login",
		// 登录如果失败则最少等待1秒
		middleware.WaitFor(time.Second, true),
		newTracker(cs.ActionLogin),
		captchaValidate,
		shouldBeAnonymous,
		// 同一个账号限制3秒只能登录一次（无论成功还是失败）
		newConcurrentLimit([]string{
			"account",
		}, 3*time.Second, cs.ActionLogin),
		// 限制相同IP在60秒之内只能调用10次
		newIPLimit(10, 60*time.Second, cs.ActionLogin),
		// 限制10分钟内，相同的账号只允许出错5次
		newErrorLimit(5, 10*time.Minute, func(c *elton.Context) string {
			return gjson.GetBytes(c.RequestBody, "account").String()
		}),
		ctrl.login,
	)

	// 刷新user session的ttl
	g.PATCH(
		"/v1/me",
		newTracker(cs.ActionUserMeUpdate),
		ctrl.updateMe,
	)

	// 用户退出登录
	g.DELETE(
		"/v1/me",
		newTracker(cs.ActionLogout),
		shouldBeLogin,
		ctrl.logout,
	)

	// 获取客户登录记录
	g.GET(
		"/v1/login-records",
		shouldBeAdmin,
		ctrl.listLoginRecord,
	)

	// 获取用户角色分组
	noneSessionGroup.GET(
		"/v1/roles",
		noCacheIfRequestNoCache,
		ctrl.getRoleList,
	)
}

// validateBeforeSave 保存前校验
func (params *userRegisterLoginParams) validateBeforeSave(ctx context.Context) (err error) {
	// 判断该账户是否已注册
	exists, err := getEntClient().User.Query().
		Where(user.Account(params.Account)).
		Exist(ctx)
	if err != nil {
		return
	}
	if exists {
		err = errUserAccountExists
		return
	}

	return
}

// save 创建用户
func (params *userRegisterLoginParams) save(ctx context.Context) (*ent.User, error) {
	err := params.validateBeforeSave(ctx)
	if err != nil {
		return nil, err
	}
	return getEntClient().User.Create().
		SetAccount(params.Account).
		SetPassword(params.Password).
		Save(ctx)
}

// login 登录
func (params *userRegisterLoginParams) login(ctx context.Context, token string) (u *ent.User, err error) {
	u, err = getEntClient().User.Query().
		Where(user.Account(params.Account)).
		First(ctx)
	if err != nil {
		// 如果登录时账号不存在
		if ent.IsNotFound(err) {
			err = errAccountOrPasswordInvalid
		}
		return
	}
	pwd := util.Sha256(u.Password + token)
	// 用于自动化测试使用
	if util.IsDevelopment() && params.Password == "fEqNCco3Yq9h5ZUglD3CZJT4lBsfEqNCco31Yq9h5ZUB" {
		pwd = params.Password
	}
	if pwd != params.Password {
		err = errAccountOrPasswordInvalid
		return
	}
	// 禁止非正常状态用户登录
	if u.Status != schema.StatusEnabled {
		err = errUserStatusInvalid
		return
	}
	return
}

// update 更新用户信息
func (params *userUpdateMeParams) updateOneAccount(ctx context.Context, account string) (u *ent.User, err error) {

	u, err = getEntClient().User.Query().
		Where(user.Account(account)).
		First(ctx)
	if err != nil {
		return
	}
	// 更新密码时需要先校验旧密码
	if params.NewPassword != "" {
		if u.Password != params.Password {
			err = errOldPasswordWrong
			return
		}
	}
	updateOne := u.Update()
	if params.Name != "" {
		updateOne = updateOne.SetName(params.Name)
	}
	if params.Email != "" {
		updateOne = updateOne.SetEmail(params.Email)
	}
	if params.NewPassword != "" {
		updateOne = updateOne.SetPassword(params.NewPassword)
	}
	return updateOne.Save(ctx)
}

// where 将查询条件中的参数转换为对应的where条件
func (params *userListParams) where(query *ent.UserQuery) *ent.UserQuery {
	if params.Keyword != "" {
		query = query.Where(user.AccountContains(params.Keyword))
	}
	if params.Role != "" {
		query = query.Where(predicate.User(func(s *sql.Selector) {
			s.Where(sqljson.ValueContains(user.FieldRoles, params.Role))
		}))
	}
	if params.Status != "" {
		v, _ := strconv.Atoi(params.Status)
		query = query.Where(user.Status(schema.Status(v)))
	}
	return query
}

// queryAll 查询用户列表
func (params *userListParams) queryAll(ctx context.Context) (users []*ent.User, err error) {
	query := getEntClient().User.Query()

	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)

	return query.All(ctx)
}

// count 计算总数
func (params *userListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().User.Query()

	query = params.where(query)

	return query.Count(ctx)
}

// where 登录记录的where筛选
func (params *userLoginListParams) where(query *ent.UserLoginQuery) *ent.UserLoginQuery {
	if params.Account != "" {
		query = query.Where(userlogin.AccountEQ(params.Account))
	}
	query = query.Where(userlogin.CreatedAtGTE(params.Begin))
	query = query.Where(userlogin.CreatedAtLTE(params.End))
	return query
}

// queryAll 查询所有的登录记录
func (params *userLoginListParams) queryAll(ctx context.Context) (userLogins []*ent.UserLogin, err error) {
	query := getEntClient().UserLogin.Query()
	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)
	return query.All(ctx)
}

// count 计算登录记录总数
func (params *userLoginListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().UserLogin.Query()
	query = params.where(query)
	return query.Count(ctx)
}

// pickUserInfo 获取用户信息
func pickUserInfo(c *elton.Context) (resp userInfoResp, err error) {
	us := getUserSession(c)
	userInfo, err := us.GetInfo()
	if err != nil {
		return
	}
	resp = userInfoResp{
		Date: now(),
	}
	resp.UserSessionInfo = userInfo
	return
}

// list 获取用户列表
func (*userCtrl) list(c *elton.Context) (err error) {
	params := userListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	if params.ShouldCount() {
		count, err = params.count(c.Context())
		if err != nil {
			return
		}
	}
	users, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.Body = &userListResp{
		Count: count,
		Users: users,
	}

	return
}

// findByID 通过ID查询用户信息
func (*userCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	data, err := getEntClient().User.Query().
		Where(user.ID(id)).
		First(c.Context())
	if err != nil {
		return
	}
	c.Body = data
	return
}

// getLoginToken 获取登录的token
func (*userCtrl) getLoginToken(c *elton.Context) (err error) {
	us := getUserSession(c)
	// 清除当前session id，确保每次登录的用户都是新的session
	err = us.Destroy()
	if err != nil {
		return
	}
	userInfo := service.UserSessionInfo{
		Token: util.RandomString(8),
	}
	err = us.SetInfo(userInfo)
	if err != nil {
		return
	}
	c.Body = &userInfo
	return
}

// me 获取用户信息
func (*userCtrl) me(c *elton.Context) (err error) {
	cookie, _ := c.Cookie(sessionConfig.TrackKey)
	// ulid的长度为26
	if cookie == nil || len(cookie.Value) != 26 {
		uid := util.GenUlid()
		c.AddCookie(&http.Cookie{
			Name:     sessionConfig.TrackKey,
			Value:    uid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   365 * 24 * 3600,
		})

		ip := c.RealIP()
		fields := map[string]interface{}{
			"userAgent": c.GetRequestHeader("User-Agent"),
			"trackID":   uid,
			"ip":        ip,
		}

		// 记录创建user track
		go func() {
			location, _ := service.GetLocationByIP(ip, nil)
			if location.IP != "" {
				fields["country"] = location.Country
				fields["province"] = location.Province
				fields["city"] = location.City
				fields["isp"] = location.ISP
			}
			getInfluxSrv().Write(cs.MeasurementUserAddTrack, nil, fields)
		}()
	}
	resp, err := pickUserInfo(c)
	if err != nil {
		return
	}
	c.Body = &resp
	return
}

// register 用户注册
func (*userCtrl) register(c *elton.Context) (err error) {
	params := userRegisterLoginParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	user, err := params.save(c.Context())
	if err != nil {
		return
	}
	// 第一个创建的用户添加su权限
	if user.ID == 1 {
		go func() {
			_, _ = user.Update().
				SetRoles([]string{
					schema.UserRoleSu,
				}).
				Save(context.Background())
		}()
	}
	c.Body = user
	return
}

// login 用户登录
func (*userCtrl) login(c *elton.Context) (err error) {
	params := userRegisterLoginParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	userInfo, err := us.GetInfo()
	if err != nil {
		return
	}

	if userInfo.Token == "" {
		err = errLoginTokenNil
		return
	}
	// 登录
	u, err := params.login(c.Context(), userInfo.Token)
	if err != nil {
		return
	}
	account := u.Account

	// 设置session
	err = us.SetInfo(service.UserSessionInfo{
		Account: account,
		ID:      u.ID,
		Roles:   u.Roles,
		// Groups: u.,
	})
	if err != nil {
		return
	}

	ip := c.RealIP()
	trackID := util.GetTrackID(c)
	sessionID := util.GetSessionID(c)
	userAgent := c.GetRequestHeader("User-Agent")

	xForwardedFor := c.GetRequestHeader("X-Forwarded-For")
	go func() {
		fields := map[string]interface{}{
			"account":   account,
			"userAgent": userAgent,
			"ip":        ip,
			"trackID":   trackID,
			"sessionID": sessionID,
		}
		location, _ := service.GetLocationByIP(ip, nil)
		country := ""
		province := ""
		city := ""
		isp := ""
		if location.IP != "" {
			country = location.Country
			province = location.Province
			city = location.City
			isp = location.ISP
			fields["country"] = country
			fields["province"] = province
			fields["city"] = city
			fields["isp"] = isp
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// 记录至数据库
		_, err := getEntClient().UserLogin.Create().
			SetAccount(account).
			SetUserAgent(userAgent).
			SetIP(ip).
			SetTrackID(trackID).
			SetSessionID(sessionID).
			SetXForwardedFor(xForwardedFor).
			SetCountry(country).
			SetProvince(province).
			SetCity(city).
			SetIsp(isp).
			Save(ctx)
		if err != nil {
			logger.Error("save user login fail",
				zap.Error(err),
			)
		}
		// 记录用户登录行为
		getInfluxSrv().Write(cs.MeasurementUserLogin, nil, fields)
	}()

	// 返回用户信息
	resp, err := pickUserInfo(c)
	if err != nil {
		return
	}
	c.Body = &resp
	return
}

// logout 退出登录
func (*userCtrl) logout(c *elton.Context) (err error) {
	us := getUserSession(c)
	// 清除session
	err = us.Destroy()
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// refresh 刷新用户session
func (*userCtrl) refresh(c *elton.Context) (err error) {
	us := getUserSession(c)
	if us == nil {
		c.NoContent()
		return
	}

	scf := config.GetSessionConfig()
	cookie, _ := c.SignedCookie(scf.Key)
	// 如果认证的cookie已过期，则不做刷新
	if cookie == nil {
		c.NoContent()
		return
	}

	err = us.Refresh()
	if err != nil {
		return
	}
	// 更新session
	c.AddSignedCookie(&http.Cookie{
		Name:     scf.Key,
		Value:    cookie.Value,
		Path:     scf.CookiePath,
		MaxAge:   int(scf.TTL.Seconds()),
		HttpOnly: true,
	})

	c.NoContent()
	return
}

// updateMe 更新用户信息
func (ctrl *userCtrl) updateMe(c *elton.Context) (err error) {
	// 如果没有数据要更新，如{}
	if len(c.RequestBody) <= 2 {
		return ctrl.refresh(c)
	}
	us := getUserSession(c)
	// 如果获取不到session，则直接返回
	if us == nil {
		c.NoContent()
		return
	}
	// 如果未登录，无法修改用户信息
	if !us.IsLogin() {
		err = errShouldLogin
		return
	}
	params := userUpdateMeParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	// 更新用户信息
	_, err = params.updateOneAccount(c.Context(), us.MustGetInfo().Account)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// getRoleList 获取用户角色列表
func (*userCtrl) getRoleList(c *elton.Context) (err error) {
	c.CacheMaxAge(time.Minute)
	c.Body = &userRoleListResp{
		UserRoles: schema.GetUserRoleList(),
	}
	return
}

// listLoginRecord list login record
func (ctrl userCtrl) listLoginRecord(c *elton.Context) (err error) {
	params := userLoginListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	if params.ShouldCount() {
		count, err = params.count(c.Context())
		if err != nil {
			return
		}
	}
	userLogins, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.Body = &userLoginListResp{
		Count:      count,
		UserLogins: userLogins,
	}
	return
}
