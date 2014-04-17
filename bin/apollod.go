package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
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

	m := martini.Classic()
	m.Use(cors.Allow(&cors.Options{
		// TODO: Get allow origin list from apollod.conf
		AllowOrigins: []string{"*"},
	}))
	m.Options("**", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	/**
	m.Handle("/tests/packages",
		api.AuthN(api.ListPackages)).Methods("GET")
	m.Handle("/tests/packages",
		api.AuthN(api.UploadPackage)).Methods("POST")
	m.Handle("/tests/packages/{filename}",
		api.AuthN(api.DetailPackage)).Methods("GET")
	m.Handle("/tests/packages/downloads/{filename}",
		api.AuthN(api.DownloadPackage)).Methods("GET")
	**/
	m.Post("/users", api.NewUser)
	m.Get("/users/:email", api.AuthN(), api.DetailUser)
	m.Post("/users/authenticate", api.Authenticate)
	/**
	m.Handle("/organizations",
		api.AuthN(api.NewOrganization)).Methods("POST")
	m.Handle("/organizations",
		api.AuthN(api.ListOrganizations)).Methods("GET")
	m.Handle("/organizations/{name}",
		api.AuthN(api.DetailOrganization)).Methods("GET")
	m.Handle("/organizations/{name}",
		api.AuthN(api.ModifyOrganization)).Methods("PUT")
	m.Handle("/organizations/{name}",
		api.AuthN(api.DeleteOrganization)).Methods("DELETE")
	m.Handle("/devices", api.AuthN(api.NewDevice)).Methods("POST")
	m.Handle("/devices", api.AuthN(api.ListDevices)).Methods("GET")
	m.Handle("/devices/{codename}",
		api.AuthN(api.DetailDevice)).Methods("GET")
	m.Handle("/devices/{codename}",
		api.AuthN(api.ModifyDevice)).Methods("PUT")
	m.Handle("/devices/{codename}",
		api.AuthN(api.DeleteDevice)).Methods("DELETE")
	m.Handle("/tests", api.AuthN(api.NewCicle)).Methods("POST")
	m.Handle("/tests", api.AuthN(api.ListCicles)).Methods("GET")
	m.Handle("/tests/{id}", api.AuthN(api.DetailCicle)).Methods("GET")
	m.Handle("/tests/{id}", api.AuthN(api.ModifyCicle)).Methods("PUT")
	m.Handle("/tests/{id}",
		api.AuthN(api.DeleteCicle)).Methods("DELETE")
	**/
	m.Any("/ws/web/:apikey", ws.Web)
	m.Any("/ws/agent/:apikey", ws.Agent)
	http.Handle("/", m)

	bind, err := config.GetString("daemon:bind")
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("Error apollod ListenAndServe: ", err)
	}
}
