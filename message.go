/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cisco/senml"
	"github.com/mainflux/mainflux-mongodb-writer/db"
	"github.com/mainflux/mainflux-mongodb-writer/models"
)

// writeMessage function
// Writtes message into DB.
func writeMessage(nm NatsMsg) error {

	Db := db.MgoDb{}
	Db.Init()
	defer Db.Close()

	var s senml.SenML
	var err error
	if s, err = senml.Decode(nm.Payload, senml.JSON); err != nil {
		println("ERROR")
		return err
	}

	// Normalize (i.e. resolve) SenMLRecord
	sn := senml.Normalize(s)

	// Timestamp
	t := time.Now().UTC().Format(time.RFC3339)
	for _, r := range sn.Records {

		m := models.Message{}

		// Copy SenMLRecord struct to Message struct
		b, err := json.Marshal(r)
		if err != nil {
			log.Print(err)
			return err
		}
		if err := json.Unmarshal(b, &m); err != nil {
			log.Print(err)
			return err
		}

		// Fill-in Mainflux stuff
		m.Channel = nm.Channel
		m.Publisher = nm.Publisher
		m.Protocol = nm.Protocol
		m.ContentType = nm.ContentType
		m.Timestamp = t

		// Insert message in DB
		if err := Db.C("messages").Insert(m); err != nil {
			log.Print(err)
			return err
		}
	}

	fmt.Println("Msg written")
	return nil
}
