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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/metadata-migration/lib/config"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
)

func init() {
	var protocolInteractionCache = map[string]string{}

	var getProtocolInteraction = func(conf config.Config, token string, protocolId string) (string, error) {
		if interaction, ok := protocolInteractionCache[protocolId]; ok {
			return interaction, nil
		}
		endpoint := conf.DeviceManagerUrl + "/protocols/" + url.PathEscape(protocolId)
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			debug.PrintStack()
			return "", err
		}
		req.Header.Set("Authorization", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			debug.PrintStack()
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			err = errors.New(buf.String())
			log.Println("ERROR: unable to get resource", endpoint, err)
			debug.PrintStack()
			return "", err
		}
		result := struct {
			Id          string `json:"id"`
			Interaction string `json:"interaction"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			debug.PrintStack()
			return "", err
		}
		protocolInteractionCache[protocolId] = result.Interaction
		return result.Interaction, nil
	}

	var getProtocolInteractionOfService = func(conf config.Config, token string, service map[string]interface{}) (string, error) {
		protocolId, ok := service["protocol_id"]
		if !ok {
			//no protocolId found
			return "", nil
		}
		protocolIdString, ok := protocolId.(string)
		if !ok {
			return "", errors.New("addserviceinteraction: unable to interpret protocol_id as string")
		}
		return getProtocolInteraction(conf, token, protocolIdString)
	}

	var getInteractionHints = func(conf config.Config, token string, service map[string]interface{}) (protocolInteraction string, inputsCount int, outputsCount int, isZwayGetLevel bool, err error) {
		protocolInteraction, err = getProtocolInteractionOfService(conf, token, service)
		if err != nil {
			return
		}

		inputs, ok := service["inputs"]
		if ok {
			inputsList, ok := inputs.([]interface{})
			if ok {
				inputsCount = len(inputsList)
			}
		}

		outputs, ok := service["outputs"]
		if ok {
			outputsList, ok := outputs.([]interface{})
			if ok {
				outputsCount = len(outputsList)
			}
		}

		localId, ok := service["local_id"]
		if ok {
			localIdStr, ok := localId.(string)
			if ok && strings.Contains(localIdStr, "get_level") {
				isZwayGetLevel = true
			}
		}
		return
	}

	var addserviceinteraction = func(conf config.Config, token string, service interface{}) (interface{}, error) {
		serviceMap, ok := service.(map[string]interface{})
		if !ok {
			return service, errors.New("addserviceinteraction: unable to interpret service as map")
		}

		serviceInteractionInterface, ok := serviceMap["interaction"]
		if !ok {
			serviceInteractionStr, ok := serviceInteractionInterface.(string)
			if ok && serviceInteractionStr != "" {
				return serviceMap, nil //service interaction is already set -> no change
			}
		}

		protocolInteraction, inputCount, outputCount, isZwayGetLevel, err := getInteractionHints(conf, token, serviceMap)
		if err != nil {
			return service, err
		}

		serviceMap["interaction"] = protocolInteraction //default

		if inputCount > 0 {
			serviceMap["interaction"] = "request"
			return serviceMap, nil
		}

		if outputCount == 0 {
			serviceMap["interaction"] = "request"
			return serviceMap, nil
		}

		if isZwayGetLevel {
			serviceMap["interaction"] = "event+request"
			return serviceMap, nil
		}

		return serviceMap, nil
	}

	KnownTransformers["addserviceinteraction"] = func(sourceConf config.Config, token string, resource string, orig interface{}) (changed interface{}, err error) {
		if resource != "device-types" {
			return orig, nil
		}

		dtMap, ok := orig.(map[string]interface{})
		if !ok {
			return orig, errors.New("addserviceinteraction: unable to interpret orig as map")
		}
		services, ok := dtMap["services"]
		if !ok {
			//no services found -> no changes
			return orig, nil
		}
		servicesList, ok := services.([]interface{})
		if !ok {
			return orig, errors.New("addserviceinteraction: unable to interpret services as array")
		}

		for index, service := range servicesList {
			servicesList[index], err = addserviceinteraction(sourceConf, token, service)
			if err != nil {
				return orig, err
			}
		}

		dtMap["services"] = servicesList

		return dtMap, nil
	}
}
