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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/metadata-migration/lib/security"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

func init() {
	Registry.Register([]string{"characteristics"}, func(library *Lib, args []string) error {
		return library.Characteristics(args)
	})
}

func (this *Lib) Characteristics(ids []string) error {
	this.VerboseLog("start to migrate characteristics")
	sourceToken, err := security.GetOpenidPasswordToken(
		this.sourceConfig.AuthUrl,
		this.sourceConfig.AuthClient,
		this.sourceConfig.AuthClientSecret,
		this.sourceConfig.SenergyUser,
		this.sourceConfig.Password)

	if err != nil {
		return err
	}

	targetToken, err := security.GetOpenidPasswordToken(
		this.targetConfig.AuthUrl,
		this.targetConfig.AuthClient,
		this.targetConfig.AuthClientSecret,
		this.targetConfig.SenergyUser,
		this.targetConfig.Password)

	if err != nil {
		return err
	}

	idChannel := listChracteristicIds(sourceToken.JwtToken(), this.sourceConfig.SourceListUrl)

	for id := range idChannel {
		this.VerboseLog(id)
		var temp interface{}
		err, _ = getResource(sourceToken.JwtToken(), this.sourceConfig.SourceSemanticUrl+"/characteristics", id.Characteristic, &temp)
		if err != nil {
			return err
		}
		transformed, err := this.transformer.Apply(this.sourceConfig, sourceToken.JwtToken(), "characteristics", temp)
		if err != nil {
			return err
		}
		err, code := setResource(targetToken.JwtToken(), this.targetConfig.DeviceManagerUrl+"/characteristics", id.Characteristic, transformed)
		if err != nil {
			this.VerboseLog(code, err)
			return err
		}
	}
	this.VerboseLog("finished to migrate characteristics")
	return nil
}

type CharacteristicId struct {
	Characteristic string `json:"id"`
	Concept        string `json:"concept_id"`
}

func listChracteristicIds(token string, endpoint string) (ids chan CharacteristicId) {
	ids = make(chan CharacteristicId, BATCH_SIZE)
	go func() {
		defer close(ids)
		limit := BATCH_SIZE
		offset := 0
		temp := []CharacteristicId{}
		for len(temp) == limit || offset == 0 {
			temp := []CharacteristicId{}
			req, err := http.NewRequest("GET", endpoint+"/jwt/list/characteristics/r/"+strconv.Itoa(limit)+"/"+strconv.Itoa(offset)+"/name/asc", nil)
			if err != nil {
				log.Println("ERROR:", err)
				debug.PrintStack()
				return
			}
			req.Header.Set("Authorization", token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("ERROR:", err)
				debug.PrintStack()
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 300 {
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				err = errors.New(buf.String())
				log.Println("ERROR: unable to get resource", endpoint, err)
				debug.PrintStack()
				return
			}
			err = json.NewDecoder(resp.Body).Decode(&temp)
			if err != nil {
				log.Println("ERROR:", err)
				debug.PrintStack()
				return
			}
			offset = offset + limit
			for _, id := range temp {
				ids <- id
			}
		}
	}()
	return ids
}
