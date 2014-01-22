package user

import (
	"code.google.com/p/go.crypto/bcrypt"
)

type User struct {
	Name     string
	Email    string
	Password string
	ApiKey   string
}

func (u *User) CryptPassword() {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Password = string(password)
	}
}
