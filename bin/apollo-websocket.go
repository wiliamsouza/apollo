package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/globocom/config"
	"github.com/gorilla/mux"

	"github.com/wiliamsouza/apollo/ws"
)

const version = "0.0.1"

func main() {
	configFile := flag.String("config", "/etc/apollo.conf", "Apollo webserver configuration file")
	gVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *gVersion {
		fmt.Printf("apollo-websocket version %s\n", version)
		return
	}
	err := config.ReadAndWatchConfigFile(*configFile)
	if err != nil {
		msg := `Could not find apollo config file. Searched on %s. For an example conf check /etc/apollo.conf file.\n %s`
		log.Panicf(msg, *configFile, err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/ws/web", ws.Web)
	r.HandleFunc("/ws/runner", ws.Runner)
	http.Handle("/", r)

	bind, err := config.GetString("websocket:bind")
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("Error apollo-websocket ListenAndServe: ", err)
	}
}
