/**
 * Copyright (c) 2017 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/nats-io/go-nats"
	"log"
	"os"
	"runtime"

	"github.com/cenkalti/backoff"
	"github.com/mainflux/mainflux-mongodb-writer/db"
)

const (
	help string = `
Usage: mainflux-influxdb [options]
Options:
	-m, --mongo MongoDB address
	-p, --mport	MongoDB port
	-d, --db	MongoDB database name
	-n, --nats	NATS address
	-q, --nport	NATS Port
	-h, --help	Prints this message end exits`
)

type (
	Opts struct {
		NatsHost string
		NatsPort string

		MongoHost     string
		MongoPort     string
		MongoDatabase string

		Help bool
	}

	NatsMsg struct {
		Channel     string `json:"channel"`
		Publisher   string `json:"publisher"`
		Protocol    string `json:"protocol"`
		ContentType string `json:"content_type"`
		Payload     []byte `json:"payload"`
	}
)

var (
	NatsConn *nats.Conn
	opts     Opts
)

func tryMongoInit() error {
	var err error

	log.Print("Connecting to MongoDB... ")
	err = db.InitMongo(opts.MongoHost, opts.MongoPort, opts.MongoDatabase)
	return err
}

func tryNatsConnect() error {
	var err error

	log.Print("Connecting to NATS... ")
	NatsConn, err = nats.Connect("nats://" + opts.NatsHost + ":" + opts.NatsPort)
	return err
}

func writerHandler(nm *nats.Msg) {
	fmt.Printf("Received a message: %s\n", string(nm.Data))

	m := NatsMsg{}
	if len(nm.Data) > 0 {
		if err := json.Unmarshal(nm.Data, &m); err != nil {
			println("Can not decode MQTT msg")
			return
		}
	}

	println("Calling writeMessage()")
	fmt.Println(m.Publisher, m.Protocol, m.Channel, m.Payload)
	writeMessage(m)
}

func main() {
	flag.StringVar(&opts.MongoHost, "m", "localhost", "MongoDB address.")
	flag.StringVar(&opts.MongoPort, "p", "27017", "MongoDB port.")
	flag.StringVar(&opts.MongoDatabase, "d", "mainflux", "MongoDB database name.")
	flag.StringVar(&opts.NatsHost, "n", "localhost", "NATS broker address.")
	flag.StringVar(&opts.NatsPort, "q", "4222", "NATS broker port.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	// MongoDb
	// Connect to MongoDB
	if err := backoff.Retry(tryMongoInit, backoff.NewExponentialBackOff()); err != nil {
		log.Fatalf("MongoDd: Can't connect: %v\n", err)
	} else {
		log.Println("OK")
	}

	// Connect to NATS broker
	if err := backoff.Retry(tryNatsConnect, backoff.NewExponentialBackOff()); err != nil {
		log.Fatalf("NATS: Can't connect: %v\n", err)
	} else {
		log.Println("OK")
	}

	// Print banner
	color.Cyan(banner)

	// Subscribe to NATS
	NatsConn.Subscribe("msg.*", writerHandler)

	// Prevent program to exit
	runtime.Goexit()
}

var banner = `
┌┬┐┌─┐┌┐┌┌─┐┌─┐┌┬┐┌┐    ┬ ┬┬─┐┬┌┬┐┌─┐┬─┐
││││ │││││ ┬│ │ ││├┴┐───│││├┬┘│ │ ├┤ ├┬┘
┴ ┴└─┘┘└┘└─┘└─┘─┴┘└─┘   └┴┘┴└─┴ ┴ └─┘┴└─

      == Industrial IoT System ==
     Made with <3 by Mainflux Team
[w] http://mainflux.io
[t] @mainflux
`
