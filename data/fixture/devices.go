package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/device"
)

func main() {
	db.Connect()

	jsonFile, err := ioutil.ReadFile("devices.json")
	if err != nil {
		panic(err)
	}
	var devices []device.Device
	err = json.Unmarshal(jsonFile, &devices)
	if err != nil {
		panic(err)
	}
	for _, d := range devices {
		if err := db.Session.Device().Insert(&d); err != nil {
			fmt.Println(err)
		}
	}
}
