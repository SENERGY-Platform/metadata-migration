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
	"fmt"
	"github.com/SENERGY-Platform/metadata-migration/lib/security"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func init() {
	Registry.Register([]string{"import"}, func(library *Lib, args []string) error {
		if len(args) == 0 {
			return errors.New("missing import/export file location")
		}
		orFilter := []string{""}
		if len(args) > 1 {
			orFilter = args[1:]
		}
		return library.Import(args[0], orFilter)
	})
}

func (this *Lib) Import(fileLocation string, orFilter []string) error {
	export, err := this.loadImportFile(fileLocation)
	if err != nil {
		return err
	}
	for _, filter := range orFilter {
		err = this.importFileWithFilter(export, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Lib) loadImportFile(fileLocation string) (result map[string]interface{}, err error) {
	var file io.Reader
	if _, err = url.ParseRequestURI(fileLocation); err != nil {
		//is not a url -> local file
		file, err = os.Open(fileLocation)
		if err != nil {
			err = fmt.Errorf("unable to open import/export file: %w", err)
			return result, err
		}
	} else {
		//is url -> online file
		resp, err := http.Get(fileLocation)
		if err != nil {
			err = fmt.Errorf("unable to load import/export file: %w", err)
			return result, err
		}
		file = resp.Body
		defer resp.Body.Close()
	}
	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		err = fmt.Errorf("unable to interpret import/export file: %w", err)
		return result, err
	}
	return result, err
}

func (this *Lib) importFileWithFilter(export map[string]interface{}, filter string) (err error) {
	token, err := security.GetOpenidPasswordToken(
		this.targetConfig.AuthUrl,
		this.targetConfig.AuthClient,
		this.targetConfig.AuthClientSecret,
		this.targetConfig.SenergyUser,
		this.targetConfig.Password)
	if err != nil {
		return fmt.Errorf("login access denied: %w", err)
	}

	for path, msg := range export {
		if strings.Contains(path, filter) {
			this.VerboseLog("export", path)
			b := new(bytes.Buffer)
			err = json.NewEncoder(b).Encode(msg)
			if err != nil {
				return fmt.Errorf("unable to send resource to target: %w", err)
			}
			req, err := http.NewRequest("PUT", this.targetConfig.DeviceManagerUrl+path, b)
			if err != nil {
				return fmt.Errorf("unable to send resource to target: %w", err)
			}
			req.Header.Set("Authorization", token.JwtToken())
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("unable to send resource to target: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 300 {
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				err = errors.New(buf.String())
				return fmt.Errorf("unable to send resource to target: %w", err)
			}
		}
	}
	return nil
}
