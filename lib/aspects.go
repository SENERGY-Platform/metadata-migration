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
)

func init() {
	Registry.Register([]string{"aspects"}, func(library *Lib, args []string) error {
		return library.Aspects(args)
	})
}

func (this *Lib) Aspects(ids []string) error {
	var idProducer IdProducer
	if len(ids) > 0 {
		idProducer = rawIdsProducer(ids)
	} else {
		idProducer = func(token string) (ids chan string) {
			return listResourceIds(token, this.sourceConfig.SourceSemanticUrl, func(limit int, offset int) string {
				//semantic /aspects does not use limit offset
				if offset == 0 {
					return "/aspects"
				} else {
					return ""
				}
			}, func(reader io.Reader) (result []IdWrapper, err error) {
				err = json.NewDecoder(reader).Decode(&result)
				return
			})
		}
	}
	return this.MigrateWithIdsProducer(
		this.sourceConfig.DeviceManagerUrl,
		this.targetConfig.DeviceManagerUrl,
		"aspects",
		idProducer)
}
