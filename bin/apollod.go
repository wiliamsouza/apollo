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

const version = "0.1.0"

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
	m.Get("/tests/packages", api.AuthN(), api.ListPackages)
	m.Post("/tests/packages", api.AuthN(), api.UploadPackage)
	m.Get("/tests/packages/:filename", api.AuthN(), api.DetailPackage)
	m.Get("/tests/packages/downloads/:filename", api.AuthN(), api.DownloadPackage)
	m.Post("/users", api.NewUser)
	m.Get("/users/:email", api.AuthN(), api.DetailUser)
	m.Post("/users/authenticate", api.Authenticate)
	m.Post("/organizations", api.AuthN(), api.NewOrganization)
	m.Get("/organizations", api.AuthN(), api.ListOrganizations)
	m.Get("/organizations/:name", api.AuthN(), api.DetailOrganization)
	m.Put("/organizations/:name", api.AuthN(), api.ModifyOrganization)
	m.Delete("/organizations/:name", api.AuthN(), api.DeleteOrganization)
	m.Post("/devices", api.AuthN(), api.NewDevice)
	m.Get("/devices", api.AuthN(), api.ListDevices)
	m.Get("/devices/:codename", api.AuthN(), api.DetailDevice)
	m.Put("/devices/:codename", api.AuthN(), api.ModifyDevice)
	m.Delete("/devices/:codename", api.AuthN(), api.DeleteDevice)

	m.Post("/tests", api.AuthN(), api.NewCicle)
	m.Get("/tests", api.AuthN(), api.ListCicles)
	m.Get("/tests/:id", api.AuthN(), api.DetailCicle)
	m.Put("/tests/:id", api.AuthN(), api.ModifyCicle)
	m.Delete("/tests/:id", api.AuthN(), api.DeleteCicle)

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
