package device

import (
	"fmt"

	"github.com/wiliamsouza/apollo/db"
)

// Device holds info about mobile devices
type Device struct {
	Codename        string      `json:"codename" bson:"_id"`
	Permission      Permissions `json:"permission" bson:"permission"`
	Owner           string      `json:"owner" bson:"owner"`
	Status          string      `json:"status" bson:"status"`
	Name            string      `json:"name" bson:"name"`
	Vendor          string      `json:"vendor" bson:"vendor"`
	Manufacturer    string      `json:"manufacturer" bson:"manufacturer"`
	Type            string      `json:"type" bson:"type"`
	Platform        string      `json:"platform" bson:"platform"`
	Cpu             string      `json:"cpu" bson:"cpu"`
	Gpu             string      `json:"gpu" bson:"gpu"`
	Ram             string      `json:"ram" bson:"ram"`
	Weight          string      `json:"weight" bson:"weight"`
	Dimensions      string      `json:"dimensions" bson:"dimentions"`
	ScreenDimension string      `json:"screenDimension" bson:"screenDimension"`
	Resolution      string      `json:"resolution" bson:"resolution"`
	ScreenDensity   string      `json:"screenDensity" bson:"screenDensity"`
	InternalStorage string      `json:"internalStorage" bson:"internalStorage"`
	SdCard          string      `json:"sdCard" bson:"sdCard"`
	Bluetooth       string      `json:"bluetooth" bson:"bluetooth"`
	WiFi            string      `json:"wifi" bson:"wifi"`
	MainCamera      string      `json:"mainCamera" bson:"mainCamera"`
	SecondaryCamera string      `json:"secondaryCamera" bson:"secondaryCamera"`
	Power           string      `json:"power" bson:"power"`
	Peripherals     string      `json:"peripherals" bson:"peripherals"`
}

type DeviceList []Device

// Permissions holds devices permission to organization and team
type Permissions struct {
	Organization Permission `json:"organization" bson:"organization"`
	Team         Permission `json:"team" bson:"team"`
}

// Permission holds what type of permission a organization and team have
type Permission struct {
	Run     bool `json:"run" bson:"run"`
	Results bool `json:"result" bson:"result"`
	Info    bool `json:"info" bson:"info"`
}

// NewDevice create new device
func NewDevice(device Device) (Device, error) {
	if err := db.Session.Device().Insert(&device); err != nil {
		return device, err
	}
	return device, nil
}

// ListDevices list devices
func ListDevices() (DeviceList, error) {
	var devices []Device
	err := db.Session.Device().Find(nil).Sort("_id").All(&devices)
	return DeviceList(devices), err
}

// DetailDevice detail device
func DetailDevice(codename string) (Device, error) {
	var device Device
	err := db.Session.Device().FindId(codename).One(&device)
	return device, err
}

// ModifyDevice modify device
func ModifyDevice(codename string, device Device) error {
	err := db.Session.Device().UpdateId(codename, device)
	if err != nil {
		return fmt.Errorf("error updating device: %s", err.Error())
	}
	return nil
}

// RemoveDevice remove device
func RemoveDevice(name string) error {
	err := db.Session.Device().RemoveId(name)
	if err != nil {
		return fmt.Errorf("error removing device: %s", err.Error())
	}
	return nil
}
