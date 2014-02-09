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
)

const version = "0.0.1"

type MuxHandler func(http.ResponseWriter, *http.Request, map[string]string)

func (h MuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h(w, r, vars)
}

func main() {
	configFile := flag.String("config", "/etc/apollo.conf", "Apollo webserver configuration file")
	gVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *gVersion {
		fmt.Printf("apollo-webserver version %s\n", version)
		return
	}
	err := config.ReadAndWatchConfigFile(*configFile)
	if err != nil {
		msg := `Could not find apollo config file. Searched on %s. For an example conf check /etc/apollo.conf file.\n %s`
		log.Panicf(msg, *configFile, err)
	}

	db.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/tests/packages", api.ListPackages).Methods("GET")
	r.HandleFunc("/tests/packages", api.UploadPackage).Methods("POST")
	r.HandleFunc("/tests/packages/{filename}", api.DetailPackage).Methods("GET")
	r.HandleFunc("/tests/packages/downloads/{filename}", api.DownloadPackage).Methods("GET")
	r.HandleFunc("/users", api.NewUser).Methods("POST")
	r.HandleFunc("/organizations", api.NewOrganization).Methods("POST")
	r.HandleFunc("/organizations", api.ListOrganizations).Methods("GET")
	r.Handle("/organizations/{name}", MuxHandler(api.DetailOrganization)).Methods("GET")
	r.Handle("/organizations/{name}", MuxHandler(api.ModifyOrganization)).Methods("PUT")
	r.Handle("/organizations/{name}", MuxHandler(api.DeleteOrganization)).Methods("DELETE")

	http.Handle("/", r)

	bind, err := config.GetString("webserver:bind")
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("Error apollo-webserver ListenAndServe: ", err)
	}
}
