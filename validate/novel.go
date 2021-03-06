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

package validate

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func init() {
	AddAlias("xNovelName", "min=1,max=30")
	AddAlias("xNovelAuthor", "min=1,max=20")
	AddAlias("xNovelID", "number")
	AddAlias("xNovelStatus", "number,min=1")
	AddAlias("xNovelSummary", "min=1,max=1000")
	AddAlias("xNovelCategory", "min=1,max=5")
	AddAlias("xNovelCoverWidth", "number")
	AddAlias("xNovelCoverHeight", "number")
	AddAlias("xNovelCoverQuality", "number")
	Add("xNovelCoverType", newIsInString([]string{
		"jpg",
		"webp",
		"png",
	}))
	AddAlias("xNovelChapterTitle", "min=1,max=1000")
	AddAlias("xNovelChapterContent", "min=1,max=50000")

	Add("xNovelIDS", func(fl validator.FieldLevel) bool {
		value, ok := toString(fl)
		if !ok {
			return false
		}
		//  长度不能超过
		if len(value) > 1000 {
			return false
		}
		for _, v := range strings.Split(value, ",") {
			_, e := strconv.Atoi(v)
			if e != nil {
				return false
			}
		}
		return true
	})
}
