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

package service

import (
	"context"
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/configuration"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/elite/request"
	"github.com/vicanso/elite/schema"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elton"
)

// ConfigurationSrv 配置的相关函数
type ConfigurationSrv struct{}

// 配置数据
type (
	// SessionInterceptorData session拦截的数据
	SessionInterceptorData struct {
		Message       string   `json:"message"`
		AllowAccounts []string `json:"allowAccounts"`
		AllowRoutes   []string `json:"allowRoutes"`
	}

	// CurrentValidConfiguration 当前有效配置
	CurrentValidConfiguration struct {
		UpdatedAt          time.Time               `json:"updatedAt"`
		MockTime           string                  `json:"mockTime"`
		IPBlockList        []string                `json:"ipBlockList"`
		SignedKeys         []string                `json:"signedKeys"`
		RouterConcurrency  map[string]uint32       `json:"routerConcurrency"`
		RouterMock         map[string]RouterConfig `json:"routerMock"`
		SessionInterceptor *SessionInterceptorData `json:"sessionInterceptor"`
	}
	// RequestLimitConfiguration HTTP请求实例并发限制
	RequestLimitConfiguration struct {
		Name string `json:"name"`
		Max  int    `json:"max"`
	}

	// ApplicationSetting 应用配置
	ApplicationSetting struct {
		Name              string `json:"name"`
		LatestVersion     string `json:"latestVersion"`
		ApplIcableVersion string `json:"applIcableVersion"`
		PrefetchSize      int    `json:"prefetchSize"`
	}
	ApplicationSettings []*ApplicationSetting
)

var (
	sessionSignedKeys = new(elton.RWMutexSignedKeys)
	// sessionInterceptorConfig session拦截的配置
	sessionInterceptorConfig = new(sync.Map)
)

// 配置刷新时间
var configurationRefreshedAt time.Time

const (
	sessionInterceptorKey = "sessionInterceptor"
)

func init() {
	sessionConfig := config.GetSessionConfig()
	// session中用于cookie的signed keys
	sessionSignedKeys.SetKeys(sessionConfig.Keys)
}

// GetSignedKeys 获取用于cookie加密的key列表
func GetSignedKeys() elton.SignedKeysGenerator {
	return sessionSignedKeys
}

// 获取首个匹配的设置
func (settings ApplicationSettings) First(currentVersion string) (setting *ApplicationSetting, err error) {
	if currentVersion == "" {
		currentVersion = "0.0.0"
	}
	ver, err := semver.Parse(currentVersion)
	if err != nil {
		return
	}
	// 对配置按期望版本号排序，新版本的排在前面
	versionReg := regexp.MustCompile(`\d+\.\d+\.\d+`)
	sort.Slice(settings, func(i, j int) bool {
		v1 := versionReg.FindString(settings[i].ApplIcableVersion)
		v2 := versionReg.FindString(settings[j].ApplIcableVersion)
		// 由于配置已按时间排序，如果版本配置相同，则返回true
		if v1 == v2 {
			return true
		}
		// 版本大的排在前面
		return v1 > v2
	})
	for _, item := range settings {
		expectedRange, e := semver.ParseRange(item.ApplIcableVersion)
		if e != nil {
			err = e
			return
		}
		if expectedRange(ver) {
			setting = item
			break
		}
	}
	return
}

// GetCurrentValidConfiguration 获取当前有效配置
func GetCurrentValidConfiguration() *CurrentValidConfiguration {
	interData, _ := GetSessionInterceptorData()
	result := &CurrentValidConfiguration{
		UpdatedAt:         configurationRefreshedAt,
		MockTime:          util.GetMockTime(),
		IPBlockList:       GetIPBlockList(),
		SignedKeys:        sessionSignedKeys.GetKeys(),
		RouterConcurrency: GetRouterConcurrency(),
		RouterMock:        GetRouterMockConfig(),
	}
	// 复制数据，避免对此数据修改
	if interData != nil {
		v := *interData
		result.SessionInterceptor = &v
	}
	return result
}

// GetSessionInterceptorMessage 获取session拦截的配置信息
func GetSessionInterceptorData() (*SessionInterceptorData, bool) {
	value, ok := sessionInterceptorConfig.Load(sessionInterceptorKey)
	if !ok {
		return nil, false
	}
	data, ok := value.(*SessionInterceptorData)
	if !ok {
		return nil, false
	}
	return data, true
}

// available 获取可用的配置
func (*ConfigurationSrv) available() ([]*ent.Configuration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	return helper.EntGetClient().Configuration.Query().
		Where(configuration.Status(schema.StatusEnabled)).
		Where(configuration.StartedAtLT(now)).
		Where(configuration.EndedAtGT(now)).
		Order(ent.Desc(configuration.FieldUpdatedAt)).
		All(ctx)
}

// ListApplicationSetting 获取应用配置
func (*ConfigurationSrv) ListApplicationSetting(ctx context.Context) (settings ApplicationSettings, err error) {
	now := time.Now()
	result, err := helper.EntGetClient().Configuration.Query().
		Where(configuration.CategoryEQ(configuration.CategoryApplicationSetting)).
		Where(configuration.Status(schema.StatusEnabled)).
		Where(configuration.StartedAtLT(now)).
		Where(configuration.EndedAtGT(now)).
		Order(ent.Desc(configuration.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return
	}
	settings = make(ApplicationSettings, len(result))
	maxPrefetchSize := 50
	for index, item := range result {
		setting := &ApplicationSetting{}
		err = json.Unmarshal([]byte(item.Data), setting)
		if err != nil {
			return
		}
		if setting.PrefetchSize > maxPrefetchSize {
			setting.PrefetchSize = maxPrefetchSize
		}
		setting.Name = item.Name
		settings[index] = setting
	}

	return
}

// Refresh 刷新配置
func (srv *ConfigurationSrv) Refresh() (err error) {
	configs, err := srv.available()
	if err != nil {
		return
	}
	configurationRefreshedAt = time.Now()
	var mockTimeConfig *ent.Configuration
	routerConcurrencyConfigs := make([]string, 0)
	routerConfigs := make([]string, 0)
	var signedKeys []string
	blockIPList := make([]string, 0)
	sessionInterceptorValue := ""

	requestLimitConfigs := make(map[string]int)
	for _, item := range configs {
		switch item.Category {
		case schema.ConfigurationCategoryMockTime:
			// 由于排序是按更新时间，因此取最新的记录
			if mockTimeConfig == nil {
				mockTimeConfig = item
			}
		case schema.ConfigurationCategoryBlockIP:
			blockIPList = append(blockIPList, item.Data)
		case schema.ConfigurationCategorySignedKey:
			// 按更新时间排序，因此如果已获取则不需要再更新
			if len(signedKeys) == 0 {
				signedKeys = strings.Split(item.Data, ",")
			}
		case schema.ConfigurationCategoryRouterConcurrency:
			routerConcurrencyConfigs = append(routerConcurrencyConfigs, item.Data)
		case schema.ConfigurationCategoryRouter:
			routerConfigs = append(routerConfigs, item.Data)
		case schema.ConfigurationCategorySessionInterceptor:
			// 按更新时间排序，因此如果已获取则不需要再更新
			if sessionInterceptorValue == "" {
				sessionInterceptorValue = item.Data
			}
		case schema.ConfigurationCategoryRequestConcurrency:
			c := RequestLimitConfiguration{}
			err := json.Unmarshal([]byte(item.Data), &c)
			if err != nil {
				log.Default().Error().
					Err(err).
					Msg("request limit config is invalid")
				AlarmError("request limit config is invalid:" + err.Error())
			}
			if c.Name != "" {
				requestLimitConfigs[c.Name] = c.Max
			}
		}
	}

	// 设置session interceptor的拦截信息
	if sessionInterceptorValue == "" {
		sessionInterceptorConfig.Delete(sessionInterceptorKey)
	} else {
		interData := &SessionInterceptorData{}
		err := json.Unmarshal([]byte(sessionInterceptorValue), interData)
		if err != nil {
			log.Default().Error().
				Err(err).
				Msg("session interceptor config is invalid")
			AlarmError("session interceptor config is invalid:" + err.Error())
		}
		sessionInterceptorConfig.Store(sessionInterceptorKey, interData)
	}

	// 如果未配置mock time，则设置为空
	if mockTimeConfig == nil {
		util.SetMockTime("")
	} else {
		util.SetMockTime(mockTimeConfig.Data)
	}

	// 如果数据库中未配置，则使用默认配置
	if len(signedKeys) == 0 {
		sessionConfig := config.GetSessionConfig()
		sessionSignedKeys.SetKeys(sessionConfig.Keys)
	} else {
		sessionSignedKeys.SetKeys(signedKeys)
	}

	// 更新router configs
	updateRouterMockConfigs(routerConfigs)

	// 重置IP拦截列表
	err = ResetIPBlocker(blockIPList)
	if err != nil {
		log.Default().Error().
			Err(err).
			Msg("reset ip blocker fail")
	}

	// 重置路由并发限制
	ResetRouterConcurrency(routerConcurrencyConfigs)

	// 更新HTTP请求实例并发限制
	request.UpdateConcurrencyLimit(requestLimitConfigs)

	return
}

func NewConfigurationSrv() *ConfigurationSrv {
	return &ConfigurationSrv{}
}
