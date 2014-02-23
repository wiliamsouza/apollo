package token

import (
	"io/ioutil"

	"github.com/globocom/config"
)

var (
	PublicKey  []byte
	PrivateKey []byte
)

func LoadKeys() {
	loadPublicKeyBytes()
	loadPrivateKeyBytes()
}

func loadPublicKeyBytes() {
	file, err := config.GetString("rsa:public")
	if err != nil {
		panic(err)
	}
	PublicKey, err = ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
}

func loadPrivateKeyBytes() {
	file, err := config.GetString("rsa:private")
	if err != nil {
		panic(err)
	}
	PrivateKey, err = ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
}
