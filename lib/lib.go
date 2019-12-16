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
)

type Lib struct {
	sourceConfig config.Config
	targetConfig config.Config
	verbose      bool
}

func New(verbose bool, source config.Config, target config.Config) *Lib {
	return &Lib{
		sourceConfig: source,
		targetConfig: target,
		verbose:      verbose,
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

func (this *Lib) Migrate(semanticSource bool, listResource string, resource string, ids []string) error {
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

	var idChannel chan string
	if len(ids) > 0 {
		idChannel = make(chan string, len(ids))
		for _, id := range ids {
			idChannel <- id
		}
		close(idChannel)
	} else {
		idChannel = listResourceIds(sourceToken.JwtToken(), this.sourceConfig.SourceListUrl, listResource)
	}

	for id := range idChannel {
		this.VerboseLog(id)
		var temp interface{}
		if semanticSource {
			err, _ = getResource(sourceToken.JwtToken(), this.sourceConfig.SourceSemanticUrl+"/"+resource, id, &temp)
		} else {
			err, _ = getResource(sourceToken.JwtToken(), this.sourceConfig.DeviceManagerUrl+"/"+resource, id, &temp)
		}
		if err != nil {
			return err
		}
		err, code := setResource(targetToken.JwtToken(), this.targetConfig.DeviceManagerUrl+"/"+resource, id, temp)
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
