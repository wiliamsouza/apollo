package device

type Device struct {
	Permision Permisions
	Owner     string
	Status    string
}

type Permissions struct {
	Organization Permision
	Team         Permision
}

type Permision struct {
	Run     bool
	Results bool
	Info    bool
}
