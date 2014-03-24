package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsuru/config"

	"github.com/wiliamsouza/apollo/api"
	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/token"
	"github.com/wiliamsouza/apollo/ws"
)

const version = "0.0.1"

func main() {
	configFile := flag.String("config", "/etc/apollo/apollod.conf",
		"Apollo daemon configuration file")
	gVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *gVersion {
		fmt.Printf("apollod version %s\n", version)
		return
	}
	err := config.ReadAndWatchConfigFile(*configFile)
	if err != nil {
		msg := `Could not find apollod config file.`
		log.Panicf(msg, *configFile, err)
	}

	token.LoadKeys()

	db.Connect()

	go ws.Bridge.Run()

	r := mux.NewRouter()
	r.Handle("/tests/packages",
		api.AuthNHandleFunc(api.ListPackages)).Methods("GET")
	r.Handle("/tests/packages",
		api.AuthNHandleFunc(api.UploadPackage)).Methods("POST")
	r.Handle("/tests/packages/{filename}",
		api.AuthNHandleFunc(api.DetailPackage)).Methods("GET")
	r.Handle("/tests/packages/downloads/{filename}",
		api.AuthNHandleFunc(api.DownloadPackage)).Methods("GET")
	r.Handle("/users", api.CORSHandle(api.PreFlightHandleFunc(api.NewUser))).Methods("POST", "OPTIONS")
	r.Handle("/users/{email}", api.AuthNHandleFunc(api.DetailUser)).Methods("GET")
	r.Handle("/users/authenticate", api.CORSHandle(api.PreFlightHandleFunc(api.Authenticate))).Methods("POST", "OPTIONS")
	r.Handle("/organizations",
		api.AuthNHandleFunc(api.NewOrganization)).Methods("POST")
	r.Handle("/organizations",
		api.AuthNHandleFunc(api.ListOrganizations)).Methods("GET")
	r.Handle("/organizations/{name}",
		api.AuthNHandleFunc(api.DetailOrganization)).Methods("GET")
	r.Handle("/organizations/{name}",
		api.AuthNHandleFunc(api.ModifyOrganization)).Methods("PUT")
	r.Handle("/organizations/{name}",
		api.AuthNHandleFunc(api.DeleteOrganization)).Methods("DELETE")
	r.Handle("/devices", api.AuthNHandleFunc(api.NewDevice)).Methods("POST")
	r.Handle("/devices", api.AuthNHandleFunc(api.ListDevices)).Methods("GET")
	r.Handle("/devices/{codename}",
		api.AuthNHandleFunc(api.DetailDevice)).Methods("GET")
	r.Handle("/devices/{codename}",
		api.AuthNHandleFunc(api.ModifyDevice)).Methods("PUT")
	r.Handle("/devices/{codename}",
		api.AuthNHandleFunc(api.DeleteDevice)).Methods("DELETE")
	r.Handle("/tests", api.AuthNHandleFunc(api.NewCicle)).Methods("POST")
	r.Handle("/tests", api.AuthNHandleFunc(api.ListCicles)).Methods("GET")
	r.Handle("/tests/{id}", api.AuthNHandleFunc(api.DetailCicle)).Methods("GET")
	r.Handle("/tests/{id}", api.AuthNHandleFunc(api.ModifyCicle)).Methods("PUT")
	r.Handle("/tests/{id}",
		api.AuthNHandleFunc(api.DeleteCicle)).Methods("DELETE")
	r.HandleFunc("/ws/web/{apikey}", ws.Web)
	r.HandleFunc("/ws/agent/{apikey}", ws.Agent)

	http.Handle("/", r)

	bind, err := config.GetString("daemon:bind")
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("Error apollod ListenAndServe: ", err)
	}
}
