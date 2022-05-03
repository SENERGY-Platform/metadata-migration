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

import "time"

func init() {
	Registry.Register([]string{"all"}, func(library *Lib, args []string) error {
		return library.All(args)
	})
	Registry.Register([]string{"iot"}, func(library *Lib, args []string) error {
		return library.Iot(args)
	})
	Registry.Register([]string{"iot-metadata"}, func(library *Lib, args []string) error {
		return library.IotMetadata(args)
	})
}

func (this *Lib) All([]string) (err error) {
	err = this.Iot([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.ProcessModels([]string{})
	if err != nil {
		return err
	}
	return nil
}

func (this *Lib) Iot([]string) (err error) {
	err = this.IotMetadata([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.Devices([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.DeviceGroups([]string{})
	if err != nil {
		return err
	}
	return nil
}

func (this *Lib) IotMetadata([]string) (err error) {
	err = this.Characteristics([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.Concepts([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.Functions([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.Aspects([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.DeviceClasses([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.Protocols([]string{})
	if err != nil {
		return err
	}
	this.VerboseLog("wait 10s for cqrs")
	time.Sleep(10 * time.Second)
	err = this.DeviceTypes([]string{})
	if err != nil {
		return err
	}
	return nil
}
