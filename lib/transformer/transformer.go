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

package transformer

import (
	"errors"
	"github.com/SENERGY-Platform/metadata-migration/lib/config"
	"strings"
)

var KnownTransformers = map[string]Transformer{}

type Transformer func(sourceConf config.Config, token string, resource string, orig interface{}) (changed interface{}, err error)

type TransformerList []Transformer

func (this *TransformerList) Apply(conf config.Config, token string, resource string, orig interface{}) (changed interface{}, err error) {
	changed = orig
	for _, transformer := range *this {
		changed, err = transformer(conf, token, resource, changed)
		if err != nil {
			return
		}
	}
	return
}

func Use(names string) (result TransformerList, err error) {
	namesList := strings.Split(names, ",")
	for _, name := range namesList {
		name = strings.TrimSpace(name)
		if name != "" {
			transformer, ok := KnownTransformers[name]
			if !ok {
				return result, errors.New("unknown transformer " + name)
			}
			result = append(result, transformer)
		}
	}
	return result, nil
}
