package customer

import (
	"encoding/base64"
	"errors"
	"regexp"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/wiliamsouza/apollo/db"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Name      string    `bson:"name"`
	Email     string    `bson:"_id"`
	Password  string    `bson:"password,omitempty"`
	ApiKey    string    `bson:"apikey,omitempty"`
	Created   time.Time `bson:"created,omitempty"`
	LastLogin time.Time `bson:"lastlogin,omitempty"`
}

func NewUser(name, email, password string) (*User, error) {
	u := &User{Name: name, Email: email, Password: password}
	v, err := u.ValidateEmail()
	if !v {
		return u, err
	}
	if err != nil {
		return u, err
	}
	u.EncryptPassword()
	u.Created = time.Now()
	u.GenerateApiKey()
	if err := db.Session.User().Insert(&u); err != nil {
		return u, err
	}
	var user *User
	_ = db.Session.User().Find(bson.M{"_id": email}).Select(bson.M{"name": 1, "email": 1}).One(&user)
	return user, nil
}

func (u *User) ValidateEmail() (bool, error) {
	m, err := regexp.Match(`^[^@]+@[^@]+\.[^@]+$`, []byte(u.Email))
	if err != nil {
		panic(err)
	}
	if !m {
		return false, errors.New("Validation Error: Email is not valid")
	}
	return true, nil
}

func (u *User) EncryptPassword() {
	if passwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost); err == nil {
		u.Password = string(passwd)
	}
}

func (u *User) GenerateApiKey() {
	token := uuid.New()
	u.ApiKey = base64.StdEncoding.EncodeToString([]byte(token))
}
