package customer

import (
	"fmt"

	"github.com/wiliamsouza/apollo/db"
)

type Organization struct {
	Name   string   `json:"name" bson:"_id"`
	Teams  []Team   `json:"teams" bson:"teams"`
	Admins []string `json:"admins" bson:"admins"`
}

type Team struct {
	Name  string   `json:"name" bson:"name"`
	Users []string `json:"users" bson:"users"`
}

type OrganizationList []Organization

func NewOrganization(organization Organization) (Organization, error) {
	if err := db.Session.Organization().Insert(&organization); err != nil {
		return organization, err
	}
	return organization, nil
}

func ListOrganizations() (OrganizationList, error) {
	var organizations []Organization
	err := db.Session.Organization().Find(nil).Sort("_id").All(&organizations)
	return OrganizationList(organizations), err

}

func DetailOrganization(name string) (Organization, error) {
	var organization Organization
	err := db.Session.Organization().FindId(name).One(&organization)
	return organization, err

}

func ModifyOrganization(name string, organization Organization) error {
	if len(organization.Admins) == 0 {
		return fmt.Errorf("Can not remove all organization admins")
	}
	err := db.Session.Organization().UpdateId(name, organization)
	if err != nil {
		return fmt.Errorf("Error updating organization: %s", err.Error())
	}
	return nil
}

func RemoveOrganization(name string) error {
	err := db.Session.Organization().RemoveId(name)
	if err != nil {
		return fmt.Errorf("Error removing organization: %s", err.Error())
	}
	return nil
}
