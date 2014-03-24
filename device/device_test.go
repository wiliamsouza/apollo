package device

import (
	"testing"

	"github.com/tsuru/config"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) SetUpSuite(c *gocheck.C) {
	err := config.ReadConfigFile("../etc/apollod.conf")
	c.Check(err, gocheck.IsNil)
	config.Set("database:name", "apollo_device_tests")
	db.Connect()
}

func (s *S) TearDownSuite(c *gocheck.C) {
	db.Session.DB.DropDatabase()
}

func (s *S) TestNewDevice(c *gocheck.C) {
	codename := "a700"
	permission := Permission{Run: true, Results: true, Info: true}
	permissions := Permissions{Organization: permission, Team: permission}
	d := Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "Acer",
		Manufacturer:    "acer",
		Type:            "tablet",
		Platform:        "NVIDIA Tegra 3",
		Cpu:             "1.3 GHz quad-core Cortex A9",
		Gpu:             "416 MHz twelve-core Nvidia GeForce ULP",
		Ram:             "1GB",
		Weight:          "665 g (1.47 lb)",
		Dimensions:      "259x175x11 mm (10.20x6.89x0.43 in)",
		ScreenDimension: "257 mm (10.1 in)",
		Resolution:      "1920x1200",
		ScreenDensity:   "224 PPI",
		InternalStorage: "32GB",
		SdCard:          "up to 32 GB",
		Bluetooth:       "yes",
		WiFi:            "802.11 b/g/n",
		MainCamera:      "5MP",
		SecondaryCamera: "2MP",
		Power:           "9800 mAh",
		Peripherals:     "accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone",
	}
	device, err := NewDevice(d)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename)
	var deviceDb Device
	err = db.Session.Device().FindId(codename).One(&deviceDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(deviceDb, gocheck.DeepEquals, device)
}

func (s *S) TestListDevice(c *gocheck.C) {
	codename1 := "a700"
	codename2 := "a701"
	permission := Permission{Run: true, Results: true, Info: true}
	permissions := Permissions{Organization: permission, Team: permission}
	d1 := Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "Acer",
		Manufacturer:    "acer",
		Type:            "tablet",
		Platform:        "NVIDIA Tegra 3",
		Cpu:             "1.3 GHz quad-core Cortex A9",
		Gpu:             "416 MHz twelve-core Nvidia GeForce ULP",
		Ram:             "1GB",
		Weight:          "665 g (1.47 lb)",
		Dimensions:      "259x175x11 mm (10.20x6.89x0.43 in)",
		ScreenDimension: "257 mm (10.1 in)",
		Resolution:      "1920x1200",
		ScreenDensity:   "224 PPI",
		InternalStorage: "32GB",
		SdCard:          "up to 32 GB",
		Bluetooth:       "yes",
		WiFi:            "802.11 b/g/n",
		MainCamera:      "5MP",
		SecondaryCamera: "2MP",
		Power:           "9800 mAh",
		Peripherals:     "accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone",
	}
	d2 := Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a701",
		Vendor:          "Acer",
		Manufacturer:    "acer",
		Type:            "tablet",
		Platform:        "NVIDIA Tegra 3",
		Cpu:             "1.3 GHz quad-core Cortex A9",
		Gpu:             "416 MHz twelve-core Nvidia GeForce ULP",
		Ram:             "1GB",
		Weight:          "665 g (1.47 lb)",
		Dimensions:      "259x175x11 mm (10.20x6.89x0.43 in)",
		ScreenDimension: "257 mm (10.1 in)",
		Resolution:      "1920x1200",
		ScreenDensity:   "224 PPI",
		InternalStorage: "32GB",
		SdCard:          "up to 32 GB",
		Bluetooth:       "yes",
		WiFi:            "802.11 b/g/n",
		MainCamera:      "5MP",
		SecondaryCamera: "2MP",
		Power:           "9800 mAh",
		Peripherals:     "accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone",
	}
	device1, err := NewDevice(d1)
	c.Assert(err, gocheck.IsNil)
	device2, err := NewDevice(d2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename1)
	defer db.Session.Device().RemoveId(codename2)
	deviceList := DeviceList{device1, device2}
	deviceListDb, err := ListDevices()
	c.Assert(err, gocheck.IsNil)
	c.Assert(deviceListDb, gocheck.DeepEquals, deviceList)
}

func (s *S) TestDetailDevice(c *gocheck.C) {
	codename := "a700"
	permission := Permission{Run: true, Results: true, Info: true}
	permissions := Permissions{Organization: permission, Team: permission}
	d := Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "Acer",
		Manufacturer:    "acer",
		Type:            "tablet",
		Platform:        "NVIDIA Tegra 3",
		Cpu:             "1.3 GHz quad-core Cortex A9",
		Gpu:             "416 MHz twelve-core Nvidia GeForce ULP",
		Ram:             "1GB",
		Weight:          "665 g (1.47 lb)",
		Dimensions:      "259x175x11 mm (10.20x6.89x0.43 in)",
		ScreenDimension: "257 mm (10.1 in)",
		Resolution:      "1920x1200",
		ScreenDensity:   "224 PPI",
		InternalStorage: "32GB",
		SdCard:          "up to 32 GB",
		Bluetooth:       "yes",
		WiFi:            "802.11 b/g/n",
		MainCamera:      "5MP",
		SecondaryCamera: "2MP",
		Power:           "9800 mAh",
		Peripherals:     "accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone",
	}
	device, err := NewDevice(d)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename)
	deviceDb, err := DetailDevice(codename)
	c.Assert(err, gocheck.IsNil)
	c.Assert(deviceDb, gocheck.DeepEquals, device)
}

func (s *S) TestModifyDevice(c *gocheck.C) {
	codename := "a700"
	permission := Permission{Run: true, Results: true, Info: true}
	permissions := Permissions{Organization: permission, Team: permission}
	d1 := Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "Acer",
		Manufacturer:    "acer",
		Type:            "tablet",
		Platform:        "NVIDIA Tegra 3",
		Cpu:             "1.3 GHz quad-core Cortex A9",
		Gpu:             "416 MHz twelve-core Nvidia GeForce ULP",
		Ram:             "1GB",
		Weight:          "665 g (1.47 lb)",
		Dimensions:      "259x175x11 mm (10.20x6.89x0.43 in)",
		ScreenDimension: "257 mm (10.1 in)",
		Resolution:      "1920x1200",
		ScreenDensity:   "224 PPI",
		InternalStorage: "32GB",
		SdCard:          "up to 32 GB",
		Bluetooth:       "yes",
		WiFi:            "802.11 b/g/n",
		MainCamera:      "5MP",
		SecondaryCamera: "2MP",
		Power:           "9800 mAh",
		Peripherals:     "accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone",
	}
	d2 := Device{
		Permission:      permissions,
		Name:            "Acer A701",
		Codename:        "a700",
		Vendor:          "Acer",
		Manufacturer:    "acer",
		Type:            "tablet",
		Platform:        "NVIDIA Tegra 4",
		Cpu:             "1.3 GHz quad-core Cortex A9",
		Gpu:             "416 MHz twelve-core Nvidia GeForce ULP",
		Ram:             "2GB",
		Weight:          "665 g (1.47 lb)",
		Dimensions:      "259x175x11 mm (10.20x6.89x0.43 in)",
		ScreenDimension: "257 mm (10.1 in)",
		Resolution:      "1920x1200",
		ScreenDensity:   "224 PPI",
		InternalStorage: "64GB",
		SdCard:          "up to 32 GB",
		Bluetooth:       "yes",
		WiFi:            "802.11 b/g/n",
		MainCamera:      "15MP",
		SecondaryCamera: "5MP",
		Power:           "9800 mAh",
		Peripherals:     "accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone",
	}
	_, err := NewDevice(d1)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename)
	err = ModifyDevice(codename, d2)
	c.Assert(err, gocheck.IsNil)
	var deviceDb Device
	err = db.Session.Device().FindId(codename).One(&deviceDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(deviceDb, gocheck.DeepEquals, d2)
}

func (s *S) TestRemoveDevice(c *gocheck.C) {
	permission := Permission{Run: true, Results: true, Info: true}
	permissions := Permissions{Organization: permission, Team: permission}
	d := Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "Acer",
		Manufacturer:    "acer",
		Type:            "tablet",
		Platform:        "NVIDIA Tegra 3",
		Cpu:             "1.3 GHz quad-core Cortex A9",
		Gpu:             "416 MHz twelve-core Nvidia GeForce ULP",
		Ram:             "1GB",
		Weight:          "665 g (1.47 lb)",
		Dimensions:      "259x175x11 mm (10.20x6.89x0.43 in)",
		ScreenDimension: "257 mm (10.1 in)",
		Resolution:      "1920x1200",
		ScreenDensity:   "224 PPI",
		InternalStorage: "32GB",
		SdCard:          "up to 32 GB",
		Bluetooth:       "yes",
		WiFi:            "802.11 b/g/n",
		MainCamera:      "5MP",
		SecondaryCamera: "2MP",
		Power:           "9800 mAh",
		Peripherals:     "accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone",
	}
	device, err := NewDevice(d)
	c.Assert(err, gocheck.IsNil)
	err = RemoveDevice(device.Codename)
	c.Assert(err, gocheck.IsNil)
	lenght, err := db.Session.Device().FindId(device.Codename).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)
}
