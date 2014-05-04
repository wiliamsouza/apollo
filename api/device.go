package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-martini/martini"

	"github.com/wiliamsouza/apollo/device"
	"github.com/wiliamsouza/apollo/token"
)

// NewDevice create new device
func NewDevice(w http.ResponseWriter, r *http.Request, token *token.Token) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var o device.Device
	err = json.Unmarshal(b, &o)
	if err != nil {
		msg := "Error parssing json request, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	newDevice, err := device.NewDevice(o)
	if err != nil {
		msg := "Error creating new device, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(newDevice)
	if err != nil {
		msg := "Error generating json result, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

// ListDevices list devices
func ListDevices(w http.ResponseWriter, r *http.Request, token *token.Token) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	devices, err := device.ListDevices()
	if err != nil {
		msg := "Error getting device list: "
		http.Error(w, msg+err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&devices)
	if err != nil {
		msg := "Error generating json result: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// DetailDevice detail device
func DetailDevice(w http.ResponseWriter, r *http.Request, token *token.Token,
	p martini.Params) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id := p["codename"]
	device, err := device.DetailDevice(id)
	if err != nil {
		msg := "Error getting device detail: "
		http.Error(w, msg+err.Error(), http.StatusNotFound)
		return
	}
	result, err := json.Marshal(&device)
	if err != nil {
		msg := "Error generating json result: "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// ModifyDevice modify device
func ModifyDevice(w http.ResponseWriter, r *http.Request, token *token.Token,
	p martini.Params) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "Error parssing request body, "
		http.Error(w, msg+err.Error(), http.StatusInternalServerError)
		return
	}
	var o device.Device
	err = json.Unmarshal(b, &o)
	if err != nil {
		msg := "Error parssing json request, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	id := p["codename"]
	err = device.ModifyDevice(id, o)
	if err != nil {
		msg := "Error updating device, "
		http.Error(w, msg+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteDevice delete device
func DeleteDevice(w http.ResponseWriter, r *http.Request, token *token.Token,
	p martini.Params) {

	id := p["codename"]
	err := device.RemoveDevice(id)
	if err != nil {
		http.Error(w, "Error deleting device, "+err.Error(),
			http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
