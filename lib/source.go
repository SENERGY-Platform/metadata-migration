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
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
)

const BATCH_SIZE = 100

func getResource(token string, endpoint string, id string, result interface{}) (err error, code int) {
	req, err := http.NewRequest("GET", endpoint+"/"+url.PathEscape(id), nil)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = errors.New(buf.String())
		log.Println("ERROR: unable to get resource", endpoint, id, err)
		debug.PrintStack()
		return err, resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		debug.PrintStack()
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

type IdProducer func(token string) (ids chan string)

func rawIdsProducer(ids []string) IdProducer {
	return func(token string) (result chan string) {
		idChannel := make(chan string, len(ids))
		for _, id := range ids {
			idChannel <- id
		}
		close(idChannel)
		return idChannel
	}
}

func permSearchIdsProducer(endpoint string, resource string) IdProducer {
	return func(token string) (ids chan string) {
		return listResourceIdsFromPermSearch(token, endpoint, resource)
	}
}

type IdWrapper struct {
	Id string `json:"id"`
}

func listResourceIdsFromPermSearch(token string, endpoint string, resource string) (ids chan string) {
	return listResourceIds(token, endpoint, func(limit int, offset int) string {
		return "/jwt/list/" + resource + "/r/" + strconv.Itoa(limit) + "/" + strconv.Itoa(offset) + "/name/asc"
	}, func(reader io.Reader) (result []IdWrapper, err error) {
		err = json.NewDecoder(reader).Decode(&result)
		return
	})
}

type PathProducer func(limit int, offset int) string
type IdListParser func(reader io.Reader) ([]IdWrapper, error)

func listResourceIds(token string, endpoint string, pathProducer PathProducer, listParser IdListParser) (ids chan string) {
	ids = make(chan string, BATCH_SIZE)
	go func() {
		defer close(ids)
		limit := BATCH_SIZE
		offset := 0
		temp := []IdWrapper{}
		for len(temp) == limit || offset == 0 {
			temp = []IdWrapper{}
			path := pathProducer(limit, offset)
			if path != "" {
				req, err := http.NewRequest("GET", endpoint+path, nil)
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
					log.Println("ERROR: unable to get resource ids", endpoint, err)
					debug.PrintStack()
					return
				}
				temp, err = listParser(resp.Body)
				if err != nil {
					log.Println("ERROR:", err)
					debug.PrintStack()
					return
				}
				offset = offset + limit
				for _, id := range temp {
					ids <- id.Id
				}
			}
		}
	}()
	return ids
}
