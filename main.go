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
	"context"
	"flag"
	"github.com/SENERGY-Platform/metadata-migration/lib"
	"github.com/SENERGY-Platform/metadata-migration/lib/config"
	"github.com/SENERGY-Platform/metadata-migration/lib/export"
	"log"
	"sync"
)

func main() {
	sourceLocation := flag.String("source", "source.json", "source configuration file")
	targetLocation := flag.String("target", "target.json", "target configuration file")
	exportTarget := flag.String("export", "", "if set the target will be ignored and the metadata will be exported to the given file")
	quiet := flag.Bool("quiet", false, "quiet log")
	flag.Parse()

	source, err := config.Load(*sourceLocation)
	if err != nil {
		log.Fatal("ERROR: unable to load source config", err)
	}

	var target config.Config
	if exportTarget != nil && *exportTarget != "" {
		wg := sync.WaitGroup{}
		defer wg.Wait()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		target, err = export.Target(*exportTarget, ctx, &wg)
		if err != nil {
			log.Fatal("ERROR: unable to start export target", err)
		}
	} else {
		target, err = config.Load(*targetLocation)
		if err != nil {
			log.Fatal("ERROR: unable to load target config", err)
		}
	}

	err = lib.New(!*quiet, source, target).Run(flag.Args())
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}
