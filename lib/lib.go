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
	"fmt"
	"github.com/SENERGY-Platform/metadata-migration/lib/config"
	"github.com/SENERGY-Platform/metadata-migration/lib/security"
	"github.com/SENERGY-Platform/metadata-migration/lib/transformer"
)

type Lib struct {
	sourceConfig config.Config
	targetConfig config.Config
	verbose      bool
	transformer  transformer.TransformerList
}

func New(verbose bool, source config.Config, target config.Config, transformer transformer.TransformerList) *Lib {
	return &Lib{
		sourceConfig: source,
		targetConfig: target,
		verbose:      verbose,
		transformer:  transformer,
	}
}

func (this *Lib) Run(args []string) (err error) {
	if len(args) == 0 {
		args = []string{"help"}
	}
	cmd, rest, err := Registry.Get(args)
	if err == CommandNotFoundError {
		fmt.Println(err)
		this.Help(nil)
		return nil
	}
	if err != nil {
		return err
	}
	return cmd(this, rest)
}

func (this *Lib) MigrateDeviceManager(listResource string, resource string, ids []string) error {
	return this.MigrateFromSourceToTarget(
		this.sourceConfig.DeviceManagerUrl,
		this.targetConfig.DeviceManagerUrl,
		listResource,
		resource,
		ids)
}

func (this *Lib) MigrateFromSourceToTarget(sourceUrl string, targetUrl string, listResource string, resource string, ids []string) error {
	var idProducer IdProducer
	if len(ids) > 0 {
		idProducer = rawIdsProducer(ids)
	} else {
		idProducer = permSearchIdsProducer(this.sourceConfig.SourceListUrl, listResource)
	}

	return this.MigrateWithIdsProducer(
		sourceUrl,
		targetUrl,
		resource,
		idProducer)
}

func (this *Lib) MigrateWithIdsProducer(sourceUrl string, targetUrl string, resource string, idProducer IdProducer) error {
	this.VerboseLog("start to migrate", resource)
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

	idChannel := idProducer(sourceToken.JwtToken())

	for id := range idChannel {
		this.VerboseLog(id)
		var temp interface{}
		err, _ = getResource(sourceToken.JwtToken(), sourceUrl+"/"+resource, id, &temp)
		if err != nil {
			return err
		}
		transformed, err := this.transformer.Apply(this.sourceConfig, sourceToken.JwtToken(), resource, temp)
		if err != nil {
			return err
		}
		err, code := setResource(targetToken.JwtToken(), targetUrl+"/"+resource, id, transformed)
		if err != nil {
			this.VerboseLog(code, err)
			return err
		}
	}
	this.VerboseLog("finished to migrate", resource)
	return nil
}

func (this *Lib) VerboseLog(msg ...interface{}) {
	if this.verbose {
		fmt.Println(msg...)
	}
}
