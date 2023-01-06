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
	"database/sql"
	"fmt"
	"github.com/SENERGY-Platform/metadata-migration/lib/security"
	"github.com/SENERGY-Platform/models/go/models"
	_ "github.com/lib/pq"
	"log"
	"os/exec"
)

func init() {
	Registry.Register([]string{"device-data"}, func(library *Lib, args []string) error {
		return library.DeviceData(args)
	})
}

func (this *Lib) DeviceData(ids []string) error {
	var idProducer IdProducer
	if len(ids) > 0 {
		idProducer = rawIdsProducer(ids)
	} else {
		idProducer = permSearchIdsProducer(this.sourceConfig.SourceListUrl, "devices")
	}
	this.VerboseLog("start to migrate device-data")
	sourceToken, err := security.GetOpenidPasswordToken(
		this.sourceConfig.AuthUrl,
		this.sourceConfig.AuthClient,
		this.sourceConfig.AuthClientSecret,
		this.sourceConfig.SenergyUser,
		this.sourceConfig.Password)

	if err != nil {
		return err
	}

	idChannel := idProducer(sourceToken.JwtToken())

	sourcePsqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", this.sourceConfig.PostgresHost,
		this.sourceConfig.PostgresPort, this.sourceConfig.PostgresUser, this.sourceConfig.PostgresPw, this.sourceConfig.PostgresDb)
	log.Println("Connecting to source PSQL...", sourcePsqlconn)
	// open database
	sourceDb, err := sql.Open("postgres", sourcePsqlconn)

	targetPsqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", this.targetConfig.PostgresHost,
		this.targetConfig.PostgresPort, this.targetConfig.PostgresUser, this.targetConfig.PostgresPw, this.targetConfig.PostgresDb)
	log.Println("Connecting to target PSQL...", targetPsqlconn)
	// open database
	targetDb, err := sql.Open("postgres", targetPsqlconn)
	if err != nil {
		return err
	}
	err = targetDb.Close() // This was just a config check
	if err != nil {
		return err
	}

	sourcePsql := "psql " + fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		this.sourceConfig.PostgresUser, this.sourceConfig.PostgresPw, this.sourceConfig.PostgresHost, this.sourceConfig.PostgresPort, this.sourceConfig.PostgresDb)
	targetPsql := "psql " + fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		this.targetConfig.PostgresUser, this.targetConfig.PostgresPw, this.targetConfig.PostgresHost, this.targetConfig.PostgresPort, this.targetConfig.PostgresDb)

	for id := range idChannel {
		this.VerboseLog(id)
		shortDeviceId, err := models.ShortenId(id)
		if err != nil {
			return err
		}
		query := "SELECT table_name FROM information_schema.tables WHERE table_name like 'device:" + shortDeviceId + "%';"
		res, err := sourceDb.Query(query)
		if err != nil {
			return err
		}
		for res.Next() {
			var table []byte
			err = res.Scan(&table)
			if err != nil {
				return err
			}
			this.VerboseLog(string(table))
			tableS := "\\\"" + string(table) + "\\\""
			cmd := sourcePsql + " -c \"\\COPY (SELECT * FROM " + tableS + ") TO stdout DELIMITER ',' CSV\" | " + targetPsql + " -c \"\\COPY " + tableS + " FROM stdin CSV\""
			ex := exec.Command("bash", "-c", cmd)
			output, err := ex.Output()
			if err != nil {
				return err
			}
			this.VerboseLog(string(output))
		}
	}

	this.VerboseLog("finished migrating device-data")
	return nil
}
