package customer

import (
	"encoding/base64"
	"errors"
	"regexp"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"labix.org/v2/mgo/bson"

	"github.com/wiliamsouza/apollo/db"
)

// User represent a system user
type User struct {
	Name      string    `bson:"name"`
	Email     string    `bson:"_id"`
	Password  string    `bson:"password"`
	APIKey    string    `bson:"apikey"`
	Created   time.Time `bson:"created"`
	LastLogin time.Time `bson:"lastlogin"`
}

// ValidateEmail check if email is valid
func (u *User) ValidateEmail() (bool, error) {
	m, err := regexp.Match(`^[^@]+@[^@]+\.[^@]+$`, []byte(u.Email))
	if err != nil {
		panic(err)
	}
	if !m {
		return false, errors.New("validation Error: Email is not valid")
	}
	return true, nil
}

// EncryptPassword before store on DB encrypt user password
func (u *User) EncryptPassword() {
	if passwd, err := bcrypt.GenerateFromPassword([]byte(u.Password),
		bcrypt.MinCost); err == nil {
		u.Password = string(passwd)
	}
}

// GenerateAPIKey for new users
func (u *User) GenerateAPIKey() {
	token := uuid.New()
	u.APIKey = base64.StdEncoding.EncodeToString([]byte(token))
}

// NewUser create a new user, check email, encrypt pass and generate APIKey
func NewUser(name, email, password string) (User, error) {
	u := User{Name: name, Email: email, Password: password}
	v, err := u.ValidateEmail()
	if !v {
		return u, err
	}
	if err != nil {
		return u, err
	}
	u.EncryptPassword()
	u.Created = time.Now()
	u.GenerateAPIKey()
	err = db.Session.User().Insert(&u)
	if err != nil {
		return u, err
	}
	var user User
	err = db.Session.User().FindId(email).One(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByAPIKey find user by APIKey
func GetUserByAPIKey(APIKey string) (User, error) {
	var u User
	err := db.Session.User().Find(bson.M{"apikey": APIKey}).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetUserByEmail find user by Email
func GetUserByEmail(email string) (User, error) {
	var u User
	err := db.Session.User().FindId(email).One(&u)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Detailuser detail user
func DetailUser(email string) (User, error) {
	var user User
	err := db.Session.User().FindId(email).One(&user)
	return user, err

}
