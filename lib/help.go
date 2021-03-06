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
	"strings"
)

func init() {
	Registry.Register([]string{"help"}, func(library *Lib, args []string) error {
		return library.Help(args)
	})
}

func (this *Lib) Help([]string) error {
	fmt.Println("ask Ingo Rößner")
	fmt.Println("some commands may accept parameters such as a list of ids")
	fmt.Println("acceptable commands:")
	for _, path := range Registry.GetPaths() {
		fmt.Println("    ", strings.Join(path, " "))
	}
	return nil
}
