package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/globocom/config"
	"github.com/gorilla/mux"

	"github.com/wiliamsouza/apollo/api"
	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/ws"
)

const version = "0.0.1"

type muxHandler func(http.ResponseWriter, *http.Request, map[string]string)

func (h muxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h(w, r, vars)
}

func main() {
	configFile := flag.String("config", "/etc/apollo/apollod.conf", "Apollo daemon configuration file")
	gVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *gVersion {
		fmt.Printf("apollod version %s\n", version)
		return
	}
	err := config.ReadAndWatchConfigFile(*configFile)
	if err != nil {
		msg := `Could not find apollo config file. Searched on %s. For an example conf check /etc/apollo/apollod.conf file.\n %s`
		log.Panicf(msg, *configFile, err)
	}

	db.Connect()

	go ws.Bridge.Run()

	r := mux.NewRouter()
	r.HandleFunc("/tests/packages", api.ListPackages).Methods("GET")
	r.HandleFunc("/tests/packages", api.UploadPackage).Methods("POST")
	r.HandleFunc("/tests/packages/{filename}", api.DetailPackage).Methods("GET")
	r.HandleFunc("/tests/packages/downloads/{filename}", api.DownloadPackage).Methods("GET")
	r.HandleFunc("/users", api.NewUser).Methods("POST")
	r.HandleFunc("/organizations", api.NewOrganization).Methods("POST")
	r.HandleFunc("/organizations", api.ListOrganizations).Methods("GET")
	r.Handle("/organizations/{name}", muxHandler(api.DetailOrganization)).Methods("GET")
	r.Handle("/organizations/{name}", muxHandler(api.ModifyOrganization)).Methods("PUT")
	r.Handle("/organizations/{name}", muxHandler(api.DeleteOrganization)).Methods("DELETE")
	r.Handle("/ws/web/{apikey}", muxHandler(ws.Web))
	r.Handle("/ws/runner/{apikey}", muxHandler(ws.Runner))

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
