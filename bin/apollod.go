package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/globocom/config"
	"github.com/gorilla/mux"

	"github.com/wiliamsouza/apollo/api"
	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/token"
	"github.com/wiliamsouza/apollo/ws"
)

const version = "0.0.1"

type authNHandler func(http.ResponseWriter, *http.Request, *jwt.Token)

func (h authNHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := token.Validate(r)
	if err != nil {
		msg := "Error not authorized: "
		http.Error(w, msg+err.Error(), http.StatusUnauthorized)
		return
	}
	h(w, r, token)
}

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
	r.HandleFunc("/tests/packages", api.ListPackages).Methods("GET")
	r.HandleFunc("/tests/packages", api.UploadPackage).Methods("POST")
	r.HandleFunc("/tests/packages/{filename}",
		api.DetailPackage).Methods("GET")
	r.HandleFunc("/tests/packages/downloads/{filename}",
		api.DownloadPackage).Methods("GET")
	r.HandleFunc("/users", api.NewUser).Methods("POST")
	r.HandleFunc("/users/{email}", api.DetailUser).Methods("GET")
	r.HandleFunc("/users/authenticate", api.Authenticate).Methods("POST")
	r.HandleFunc("/organizations", api.NewOrganization).Methods("POST")
	r.HandleFunc("/organizations", api.ListOrganizations).Methods("GET")
	r.HandleFunc("/organizations/{name}", api.DetailOrganization).Methods("GET")
	r.HandleFunc("/organizations/{name}", api.ModifyOrganization).Methods("PUT")
	r.HandleFunc("/organizations/{name}", api.DeleteOrganization).Methods("DELETE")
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
