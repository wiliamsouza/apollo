package customer

import (
	"fmt"

	"github.com/wiliamsouza/apollo/db"
)

// Organization holds its name, teams and admins
type Organization struct {
	Name   string   `json:"name" bson:"_id"`
	Teams  []Team   `json:"teams" bson:"teams"`
	Admins []string `json:"admins" bson:"admins"`
}

// Team holds its name and users list
type Team struct {
	Name  string   `json:"name" bson:"name"`
	Users []string `json:"users" bson:"users"`
}

// OrganizationList holds a list of organization
type OrganizationList []Organization

// NewOrganization create new organization
func NewOrganization(organization Organization) (Organization, error) {
	if err := db.Session.Organization().Insert(&organization); err != nil {
		return organization, err
	}
	return organization, nil
}

// ListOrganizations list organizations
func ListOrganizations() (OrganizationList, error) {
	var organizations []Organization
	err := db.Session.Organization().Find(nil).Sort("_id").All(&organizations)
	return OrganizationList(organizations), err

}

// DetailOrganization detail organization
func DetailOrganization(name string) (Organization, error) {
	var organization Organization
	err := db.Session.Organization().FindId(name).One(&organization)
	return organization, err

}

// ModifyOrganization modify organization
func ModifyOrganization(name string, organization Organization) error {
	if len(organization.Admins) == 0 {
		return fmt.Errorf("can not remove all organization admins")
	}
	err := db.Session.Organization().UpdateId(name, organization)
	if err != nil {
		return fmt.Errorf("error updating organization: %s", err.Error())
	}
	return nil
}

// RemoveOrganization remove organization
func RemoveOrganization(name string) error {
	err := db.Session.Organization().RemoveId(name)
	if err != nil {
		return fmt.Errorf("error removing organization: %s", err.Error())
	}
	return nil
}
