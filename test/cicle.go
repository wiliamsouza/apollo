package test

import (
	"errors"
	"fmt"

	"labix.org/v2/mgo/bson"

	"github.com/wiliamsouza/apollo/db"
)

type Cicle struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Device   string        `json:"device" bson:"device"`
	Packages []string      `json:"packages" bson:""packages`
}

type CicleList []Cicle

// NewCicle create new cicle
func NewCicle(cicle Cicle) (Cicle, error) {
	cicle.Id = bson.NewObjectId()
	if err := db.Session.Cicle().Insert(&cicle); err != nil {
		return cicle, err
	}
	return cicle, nil
}

// ListCicles list cicles
func ListCicles() (CicleList, error) {
	var cicles []Cicle
	err := db.Session.Cicle().Find(nil).Sort("_id").All(&cicles)
	return CicleList(cicles), err
}

// DetailCicle detail cicle
func DetailCicle(objIdHex string) (Cicle, error) {
	if !bson.IsObjectIdHex(objIdHex) {
		return Cicle{}, errors.New("Invalid cicle object id hex")
	}
	id := bson.ObjectIdHex(objIdHex)
	var cicle Cicle
	err := db.Session.Cicle().FindId(id).One(&cicle)
	return cicle, err
}

// ModifyCicle modify cicle
func ModifyCicle(objIdHex string, cicle Cicle) error {
	if !bson.IsObjectIdHex(objIdHex) {
		return fmt.Errorf("Invalid cicle object id hex")
	}
	id := bson.ObjectIdHex(objIdHex)
	err := db.Session.Cicle().UpdateId(id, cicle)
	if err != nil {
		return fmt.Errorf("error updating cicle: %s", err.Error())
	}
	return nil
}

// RemoveCicle remove cicle
func RemoveCicle(objIdHex string) error {
	if !bson.IsObjectIdHex(objIdHex) {
		return fmt.Errorf("Invalid cicle object id hex")
	}
	id := bson.ObjectIdHex(objIdHex)
	err := db.Session.Cicle().RemoveId(id)
	if err != nil {
		return fmt.Errorf("error removing cicle: %s", err.Error())
	}
	return nil
}
