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

package util

import (
	"regexp"

	"github.com/vicanso/elite/config"
	"github.com/vicanso/elton"
)

var sessionConfig = config.GetSessionConfig()
var uuidReg = regexp.MustCompile(`uuid/(\S+)`)
var versionReg = regexp.MustCompile(`elite/(\S+)`)

// GetDeviceID 获取设备ID
func GetDeviceID(c *elton.Context) string {
	deviceID := ""
	arr := uuidReg.FindStringSubmatch(c.Request.UserAgent())
	if len(arr) == 2 {
		deviceID = arr[1]
	}
	if deviceID == "" {
		deviceID = GetTrackID(c)
	}
	return deviceID
}

// GetAppVersion 获取应用版本
func GetAppVersion(c *elton.Context) string {
	version := ""
	arr := versionReg.FindStringSubmatch(c.Request.UserAgent())
	if len(arr) == 2 {
		version = arr[1]
	}
	return version
}

// GetTrackID 获取track id
func GetTrackID(c *elton.Context) string {
	trackCookie := sessionConfig.TrackKey
	if trackCookie == "" {
		return ""
	}
	cookie, _ := c.Cookie(trackCookie)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

// GetSessionID 获取session id
func GetSessionID(c *elton.Context) string {
	cookie, _ := c.Cookie(sessionConfig.Key)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}
