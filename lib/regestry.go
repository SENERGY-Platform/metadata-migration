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
	"sort"
	"strings"
)

var Registry = &RegistryElement{}

type Command func(library *Lib, args []string) error

type RegistryElement struct {
	command Command
	sub     map[string]*RegistryElement
}

func (this *RegistryElement) Register(path []string, cmd Command) {
	if len(path) == 0 {
		this.command = cmd
	} else {
		next := path[0]
		rest := path[1:]
		if this.sub == nil {
			this.sub = map[string]*RegistryElement{}
		}
		reg, ok := this.sub[next]
		if !ok {
			reg = &RegistryElement{}
			this.sub[next] = reg
		}
		reg.Register(rest, cmd)
	}
}

var CommandNotFoundError = errors.New("command not found")

func (this *RegistryElement) Get(path []string) (cmd Command, rest []string, err error) {
	if this.command != nil {
		return this.command, path, nil
	}
	if this.sub == nil {
		return cmd, path, CommandNotFoundError
	}
	if len(path) == 0 {
		options := []string{}
		for option, _ := range this.sub {
			options = append(options, option)
		}
		return cmd, path, errors.New("incomplete command, use one of the following:\n\t" + strings.Join(options, "\n\t"))
	}
	next := path[0]
	rest = path[1:]
	reg, ok := this.sub[next]
	if !ok {
		return cmd, path, CommandNotFoundError
	}
	return reg.Get(rest)
}

func (this *RegistryElement) GetChild(path []string) (result *RegistryElement, err error) {
	if len(path) == 0 {
		return this, nil
	}
	sub, ok := this.sub[path[0]]
	if !ok {
		return result, errors.New("not found")
	}
	return sub.GetChild(path[1:])
}

func (this *RegistryElement) GetPaths() [][]string {
	paths := this.getPaths([]string{})
	sort.Slice(paths, func(i, j int) bool {
		a := strings.Join(paths[i], " ")
		b := strings.Join(paths[j], " ")
		return a < b
	})
	return paths
}

func (this *RegistryElement) getPaths(current []string) (result [][]string) {
	if this.command != nil {
		return [][]string{current}
	}
	for key, sub := range this.sub {
		next := append(current, key)
		result = append(result, sub.getPaths(next)...)
	}
	return result
}
