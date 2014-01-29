package device

type Device struct {
	Permission      Permissions
	Owner           string
	Status          string
	Name            string
	Codename        string
	Vendor          string
	Manufacturer    string
	Type            string
	Platform        string
	Cpu             string
	Gpu             string
	Ram             string
	Weight          string
	Dimensions      string
	ScreenDimension string
	Resolution      string
	ScreenDensity   string
	InternalStorage string
	SdCard          string
	Bluetooth       string
	WiFi            string
	MainCamera      string
	SecondaryCamera string
	Power           string
	Peripherals     string
}

type Permissions struct {
	Organization Permission
	Team         Permission
}

type Permission struct {
	Run     bool
	Results bool
	Info    bool
}
