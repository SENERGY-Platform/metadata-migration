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

package export

import (
	"context"
	"encoding/json"
	"github.com/SENERGY-Platform/metadata-migration/lib/config"
	"github.com/SENERGY-Platform/metadata-migration/lib/security"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
)

func Target(location string, ctx context.Context, wg *sync.WaitGroup) (conf config.Config, err error) {
	export := newExport(location)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		export.Flush()
	}()
	conf.AuthUrl = auth(ctx, wg)
	conf.DeviceManagerUrl = deviceManager(export, ctx, wg)
	return conf, nil
}

func deviceManager(export Export, ctx context.Context, wg *sync.WaitGroup) string {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		temp, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}
		export.Add(request.URL.Path, temp)
	}))
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		server.Close()
	}()
	return server.URL
}

type Export interface {
	Add(path string, body []byte)
	Flush()
}

type JsonFileExport struct {
	location string
	values   map[string]interface{}
}

func (this *JsonFileExport) Add(path string, body []byte) {
	var bodyObj interface{}
	err := json.Unmarshal(body, &bodyObj)
	if err != nil {
		log.Fatal(err)
	}
	this.values[path] = bodyObj
}

func (this *JsonFileExport) Flush() {
	file, err := json.MarshalIndent(this.values, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(this.location, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func newExport(location string) (export Export) {
	return &JsonFileExport{
		location: location,
		values:   map[string]interface{}{},
	}
}

func auth(ctx context.Context, wg *sync.WaitGroup) (url string) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		json.NewEncoder(writer).Encode(security.OpenidToken{
			AccessToken: "export",
		})
	}))
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		server.Close()
	}()
	return server.URL
}
