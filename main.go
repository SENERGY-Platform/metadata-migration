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

package main

import (
	"flag"
	"github.com/SENERGY-Platform/metadata-migration/lib"
	"github.com/SENERGY-Platform/metadata-migration/lib/config"
	"log"
)

func main() {
	sourceLocation := flag.String("source", "source.json", "source configuration file")
	targetLocation := flag.String("target", "target.json", "target configuration file")
	quiet := flag.Bool("quiet", false, "quiet log")
	flag.Parse()

	source, err := config.Load(*sourceLocation)
	if err != nil {
		log.Fatal("ERROR: unable to load source config", err)
	}

	target, err := config.Load(*targetLocation)
	if err != nil {
		log.Fatal("ERROR: unable to load target config", err)
	}

	err = lib.New(!*quiet, source, target).Run(flag.Args())
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}
