// Copyright 2018, Andrew C. Young
// License: MIT

package iot_test

import (
	"fmt"
	"io/ioutil"

	"github.com/vaelen/iot"
)

func ExampleThing() {
	id := &iot.ID{
		DeviceID:  "deviceName",
		Registry:  "my-registry",
		Location:  "asia-east1",
		ProjectID: "my-project",
	}

	credentials, err := iot.LoadCredentials("rsa_cert.pem", "rsa_private.pem")
	if err != nil {
		panic("Couldn't load credentials")
	}

	tmpDir, err := ioutil.TempDir("", "queue-")
	if err != nil {
		panic("Couldn't create temp directory")
	}

	options := iot.DefaultOptions(id, credentials)
	options.DebugLogger = func(a ...interface{}) { fmt.Println(a...) }
	options.InfoLogger = func(a ...interface{}) { fmt.Println(a...) }
	options.ErrorLogger = func(a ...interface{}) { fmt.Println(a...) }
	options.QueueDirectory = tmpDir
	options.ConfigHandler = func(thing iot.Thing, config []byte) {
		// Do something here to process the updated config and create an updated state string
		state := []byte("ok")
		thing.PublishState(state)
	}

	thing := iot.New(options)

	err = thing.Connect("ssl://mqtt.googleapis.com:443")
	if err != nil {
		panic("Couldn't connect to server")
	}
	defer thing.Disconnect()

	// This publishes to /events
	thing.PublishEvent([]byte("Top level telemetry event"))
	// This publishes to /events/a
	thing.PublishEvent([]byte("Sub folder telemetry event"), "a")
	// This publishes to /events/a/b
	thing.PublishEvent([]byte("Sub folder telemetry event"), "a", "b")
}