package customer

import (
	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
)

func (s *S) TestNewOrganization(c *gocheck.C) {
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	id := "doegroup"
	t := Team{Name: "Test", Users: []string{jhon, jane}}
	o := Organization{Name: id, Teams: []Team{t}, Admins: []string{jhon}}
	org, err := NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(id)
	var orgDb Organization
	err = db.Session.Organization().FindId(id).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestListOrganizations(c *gocheck.C) {
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	id1 := "jhoncorp"
	id2 := "janecorp"
	t := Team{Name: "Test", Users: []string{jhon, jane}}
	o1 := Organization{Name: id1, Teams: []Team{t}, Admins: []string{jhon}}
	o2 := Organization{Name: id2, Teams: []Team{t}, Admins: []string{jane}}
	_, err := NewOrganization(o1)
	c.Assert(err, gocheck.IsNil)
	_, err = NewOrganization(o2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(id1)
	defer db.Session.Organization().RemoveId(id2)
	// BUG: The order here is relevant cause .Sort("_id")
	orgList := OrganizationList{o2, o1}
	orgListDb, err := ListOrganizations()
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgListDb, gocheck.DeepEquals, orgList)
}

func (s *S) TestDetailOrganization(c *gocheck.C) {
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	id := "doegroup"
	t := Team{Name: "Test", Users: []string{jhon, jane}}
	o := Organization{Name: id, Teams: []Team{t}, Admins: []string{jhon}}
	org, err := NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(id)
	orgDb, err := DetailOrganization(id)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestModifyOrganization(c *gocheck.C) {
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	id := "jhoncorp"
	t1 := Team{Name: "Test", Users: []string{jhon}}
	t2 := Team{Name: "Test", Users: []string{jhon, jane}}
	o1 := Organization{Name: id, Teams: []Team{t1}, Admins: []string{jhon}}
	o2 := Organization{Name: id, Teams: []Team{t2}, Admins: []string{jane}}
	_, err := NewOrganization(o1)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Organization().RemoveId(id)
	err = ModifyOrganization(id, o2)
	c.Assert(err, gocheck.IsNil)
	var orgDb Organization
	err = db.Session.Organization().FindId(id).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, o2)
}

func (s *S) TestRemoveOrganization(c *gocheck.C) {
	jhon := "jhon@doe.com"
	jane := "jane@doe.com"
	id := "doegroup"
	t := Team{Name: "Test", Users: []string{jhon, jane}}
	o := Organization{Name: id, Teams: []Team{t}, Admins: []string{jhon}}
	org, err := NewOrganization(o)
	c.Assert(err, gocheck.IsNil)
	err = RemoveOrganization(org.Name)
	c.Assert(err, gocheck.IsNil)
	lenght, err := db.Session.Organization().FindId(org.Name).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)

}
