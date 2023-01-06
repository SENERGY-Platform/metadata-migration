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
	"errors"
	"github.com/SENERGY-Platform/metadata-migration/lib/security"
)

func init() {
	Registry.Register([]string{"dashboards"}, func(library *Lib, args []string) error {
		return library.Dashboards(args)
	})
}

func (this *Lib) Dashboards(ids []string) error {
	this.VerboseLog("Start migrating dashboards")
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

	dashboards := []map[string]interface{}{}

	err, _ = getResource(sourceToken.JwtToken(), this.sourceConfig.DashboardUrl, "dashboards", &dashboards)
	if err != nil {
		return err
	}

	for _, dashboard := range dashboards {
		id, ok := dashboard["id"].(string)
		if !ok {
			return errors.New("unexpected dashboard result")
		}
		this.VerboseLog(id)
		if len(ids) == 0 || has(ids, id) {
			err, _ = setResource(targetToken.JwtToken(), this.targetConfig.DashboardUrl, "dashboard", dashboard)
			if err != nil {
				return err
			}
		}
	}

	this.VerboseLog("Finished migrating dashboards")
	return nil
}

func has[E comparable](slice []E, elem E) bool {
	for i := range slice {
		if elem == slice[i] {
			return true
		}
	}
	return false
}
