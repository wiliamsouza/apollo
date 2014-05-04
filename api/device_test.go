package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-martini/martini"
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/device"
	"github.com/wiliamsouza/apollo/token"
)

func (s *S) TestNewDevice(c *gocheck.C) {
	result := `{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}`
	defer db.Session.Device().DropCollection()
	body := strings.NewReader(result)
	request, err := http.NewRequest("POST", "devices", body)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	NewDevice(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestNewDeviceInvalidJson(c *gocheck.C) {
	invalid := `{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"`
	result := "Error parssing json request, unexpected end of JSON input\n"
	body := strings.NewReader(invalid)
	request, err := http.NewRequest("POST", "devices", body)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	NewDevice(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusBadRequest)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "text/plain; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestListDevices(c *gocheck.C) {
	results := `[{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"},{"codename":"a701","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}]`
	codename1 := "a700"
	codename2 := "a701"
	permission := device.Permission{Run: false, Results: false, Info: false}
	permissions := device.Permissions{Organization: permission, Team: permission}
	d1 := device.Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "acer",
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
	d2 := device.Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a701",
		Vendor:          "acer",
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
	_, err := device.NewDevice(d1)
	c.Assert(err, gocheck.IsNil)
	_, err = device.NewDevice(d2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename1)
	defer db.Session.Device().RemoveId(codename2)
	request, err := http.NewRequest("GET", "devices", nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	ListDevices(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestDetailDevice(c *gocheck.C) {
	result := `{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}`
	codename := "a700"
	permission := device.Permission{Run: false, Results: false, Info: false}
	permissions := device.Permissions{Organization: permission, Team: permission}
	d := device.Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "acer",
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
	_, err := device.NewDevice(d)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename)
	url := fmt.Sprintf("devices/%s", codename)
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	p := make(map[string]string)
	p["codename"] = codename
	params := martini.Params(p)
	DetailDevice(response, request, tk, params)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestModifyDevice(c *gocheck.C) {
	b := `{"codename":"a700","permission":{"organization":{"run":true,"result":true,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"jhon@doe.com","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}`
	body := strings.NewReader(b)
	codename := "a700"
	permission := device.Permission{Run: false, Results: false, Info: false}
	permissions := device.Permissions{Organization: permission, Team: permission}
	d := device.Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "acer",
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
	_, err := device.NewDevice(d)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename)
	var org device.Device
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("devices/%s", codename)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	p := make(map[string]string)
	p["codename"] = codename
	params := martini.Params(p)
	ModifyDevice(response, request, tk, params)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	var orgDb device.Device
	err = db.Session.Device().FindId(codename).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestDeleteDevice(c *gocheck.C) {
	codename := "a700"
	permission := device.Permission{Run: false, Results: false, Info: false}
	permissions := device.Permissions{Organization: permission, Team: permission}
	d := device.Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "acer",
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
	_, err := device.NewDevice(d)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("devices/%s", codename)
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	p := make(map[string]string)
	p["codename"] = codename
	params := martini.Params(p)
	DeleteDevice(response, request, tk, params)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	lenght, err := db.Session.Device().FindId(codename).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)
}

func (s *S) TestDeleteDeviceNoExist(c *gocheck.C) {
	codename := "a700"
	result := "Error deleting device, error removing device: not found\n"
	permission := device.Permission{Run: false, Results: false, Info: false}
	permissions := device.Permissions{Organization: permission, Team: permission}
	d := device.Device{
		Permission:      permissions,
		Name:            "Acer A700",
		Codename:        "a700",
		Vendor:          "acer",
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
	_, err := device.NewDevice(d)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Device().RemoveId(codename)
	url := fmt.Sprintf("devices/%s", "noexist")
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tt, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	tk := &token.Token{Email: tt.Claims["email"].(string), Exp: tt.Claims["exp"].(float64)}
	p := make(map[string]string)
	p["codename"] = "noexist"
	params := martini.Params(p)
	DeleteDevice(response, request, tk, params)
	lenght, err := db.Session.Device().FindId(codename).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 1)
	c.Assert(response.Code, gocheck.Equals, http.StatusBadRequest)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "text/plain; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}
