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
	r.Handle("/tests/packages",
		authNHandler(api.ListPackages)).Methods("GET")
	r.Handle("/tests/packages",
		authNHandler(api.UploadPackage)).Methods("POST")
	r.Handle("/tests/packages/{filename}",
		authNHandler(api.DetailPackage)).Methods("GET")
	r.Handle("/tests/packages/downloads/{filename}",
		authNHandler(api.DownloadPackage)).Methods("GET")
	r.HandleFunc("/users", api.NewUser).Methods("POST")
	r.Handle("/users/{email}", authNHandler(api.DetailUser)).Methods("GET")
	r.HandleFunc("/users/authenticate", api.Authenticate).Methods("POST")
	r.Handle("/organizations",
		authNHandler(api.NewOrganization)).Methods("POST")
	r.Handle("/organizations",
		authNHandler(api.ListOrganizations)).Methods("GET")
	r.Handle("/organizations/{name}",
		authNHandler(api.DetailOrganization)).Methods("GET")
	r.Handle("/organizations/{name}",
		authNHandler(api.ModifyOrganization)).Methods("PUT")
	r.Handle("/organizations/{name}",
		authNHandler(api.DeleteOrganization)).Methods("DELETE")
	r.Handle("/devices", authNHandler(api.NewDevice)).Methods("POST")
	r.Handle("/devices", authNHandler(api.ListDevices)).Methods("GET")
	r.Handle("/devices/{codename}",
		authNHandler(api.DetailDevice)).Methods("GET")
	r.Handle("/devices/{codename}",
		authNHandler(api.ModifyDevice)).Methods("PUT")
	r.Handle("/devices/{codename}",
		authNHandler(api.DeleteDevice)).Methods("DELETE")
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
