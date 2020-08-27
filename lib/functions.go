/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"encoding/json"
	"io"
	"strconv"
)

func init() {
	Registry.Register([]string{"functions"}, func(library *Lib, args []string) error {
		return library.Functions(args)
	})
}

type FunctionsIdListResult struct {
	Functions []IdWrapper `json:"functions"`
}

func (this *Lib) Functions(ids []string) error {
	var idProducer IdProducer
	if len(ids) > 0 {
		idProducer = rawIdsProducer(ids)
	} else {
		idProducer = func(token string) (ids chan string) {
			return listResourceIds(token, this.sourceConfig.SourceSemanticUrl, func(limit int, offset int) string {
				return "/functions?limit" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
			}, func(reader io.Reader) (result []IdWrapper, err error) {
				temp := FunctionsIdListResult{}
				err = json.NewDecoder(reader).Decode(&temp)
				return temp.Functions, err
			})
		}
	}
	return this.MigrateWithIdsProducer(
		this.sourceConfig.DeviceManagerUrl,
		this.targetConfig.DeviceManagerUrl,
		"functions",
		idProducer)
}
