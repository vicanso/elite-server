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
package cs

const (
	// ActionLogin login
	ActionLogin = "login"
	// ActionRegister register
	ActionRegister = "register"
	// ActionLogout logout
	ActionLogout = "logout"

	// ActionUserInfoUpdate update user info
	ActionUserInfoUpdate = "updateUserInfo"
	// ActionUserMeUpdate update my info
	ActionUserMeUpdate = "updateUserMe"
	// ActionAddUserTracker add user tracker
	ActionAddUserTracker = "addUserTracker"

	// ActionConfigurationAdd add configuration
	ActionConfigurationAdd = "addConfiguration"
	// ActionConfigurationUpdate update configuration
	ActionConfigurationUpdate = "updateConfiguration"

	// ActionAdminCleanSession clean session
	ActionAdminCleanSession = "cleanSession"
)

// 小说相关的操作
const (
	// ActionNovelUpdate update novel
	ActionNovelUpdate = "updateNovel"
	// ActionNovelChaptersUpdate update novel chapters
	ActionNovelChaptersUpdate = "updateNovelChapters"
	// ActionNovelChapterUpdate update novel chapter
	ActionNovelChapterUpdate = "updateNovelChapter"
)

// 客户端的相关操作
const (
	ActionContinueReading    = "continueReading"
	ActionFetchMoreNovel     = "fetchMoreNovel"
	ActionNovelDetail        = "novelDetail"
	ActionChapterList        = "chapterList"
	ActionChapterDetail      = "chapterDetail"
	ActionAddToFavorite      = "addToFavorite"
	ActionRemoveFromFavorite = "removeFromFavorite"
	ActionFetchCategoryNovel = "fetchCategoryNovel"
)
