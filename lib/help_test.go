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
	"strings"
)

func ExampleHelp() {
	command := "help"
	commandArgs := strings.Split(command, " ")
	fmt.Println(New(config.Config{}, config.Config{}).Run(commandArgs))

	//output:
	//ask Ingo Rößner
	//<nil>
}

func ExampleHelpWithIgnoredArgs() {
	command := "help something unknown"
	commandArgs := strings.Split(command, " ")
	fmt.Println(New(config.Config{}, config.Config{}).Run(commandArgs))

	//output:
	//ask Ingo Rößner
	//<nil>
}

func ExampleBadHelp() {
	command := "bad help"
	commandArgs := strings.Split(command, " ")
	fmt.Println(New(config.Config{}, config.Config{}).Run(commandArgs))

	//output:
	//'bad help' command not found
}
