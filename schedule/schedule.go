// Copyright 2019 tree xie
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

package schedule

import (
	"github.com/robfig/cron/v3"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/elite/service"

	"go.uber.org/zap"
)

func init() {
	// go func() {
	// 	err := new(service.LongzuSrv).Sync()
	// 	if err != nil {
	// 		log.Default().Error("sync longzu fail",
	// 			zap.Error(err),
	// 		)
	// 	} else {
	// 		log.Default().Info("sync longzu success")
	// 	}
	// }()
	c := cron.New()
	_, _ = c.AddFunc("@every 5m", redisCheck)
	_, _ = c.AddFunc("@every 1m", configRefresh)
	_, _ = c.AddFunc("@every 10m", novelBasicInfoRefresh)
	_, _ = c.AddFunc("00 00 * * *", resetNovelSearchHotKeywords)
	_, _ = c.AddFunc("@every 2h", updateNovelUnfinished)
	c.Start()
}

func redisCheck() {
	err := helper.RedisPing()
	if err != nil {
		log.Default().Error("redis check fail",
			zap.Error(err),
		)
		service.AlarmError("redis check fail")
	}
}

func configRefresh() {
	configSrv := new(service.ConfigurationSrv)
	err := configSrv.Refresh()
	if err != nil {
		log.Default().Error("config refresh fail",
			zap.Error(err),
		)
		service.AlarmError("config refresh fail")
	}
}

func novelBasicInfoRefresh() {
	err := new(service.NovelSrv).RefreshAllBasicInfo()
	if err != nil {
		log.Default().Error("novel basic info refresh fail",
			zap.Error(err),
		)
	}
}

func resetNovelSearchHotKeywords() {
	_, err := helper.RedisGetClient().Del(cs.NovelSearchHotKeyWords).Result()
	if err != nil {
		log.Default().Error("reset novel search hot key words fail",
			zap.Error(err),
		)
	}
}

func updateNovelUnfinished() {
	err := new(service.NovelSrv).UpdateUnfinished()
	if err != nil {
		log.Default().Error("novel update unfinished fail",
			zap.Error(err),
		)
	}
}
