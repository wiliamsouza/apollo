package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"launchpad.net/gocheck"

	"github.com/wiliamsouza/apollo/db"
	"github.com/wiliamsouza/apollo/test"
	"github.com/wiliamsouza/apollo/token"
)

func (s *S) TestNewCicle(c *gocheck.C) {
	// It returns Id too we remove it cause we can't predict before hand what it will be
	b := `{"name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]}`
	defer db.Session.Device().DropCollection()
	body := strings.NewReader(b)
	request, err := http.NewRequest("POST", "cicles", body)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewCicle(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusCreated)
	var r test.Cicle
	err = json.Unmarshal([]byte(response.Body.String()), &r)
	c.Assert(err, gocheck.IsNil)
	cicle := test.Cicle{Name: r.Name, Device: r.Device, Packages: r.Packages}
	result, err := json.Marshal(cicle)
	c.Assert(err, gocheck.IsNil)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(string(result), gocheck.Equals, b)
}

func (s *S) TestNewCicleInvalidJson(c *gocheck.C) {
	result := "Error parssing json request, unexpected end of JSON input\n"
	invalid := `{"name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]`
	body := strings.NewReader(invalid)
	request, err := http.NewRequest("POST", "cicles", body)
	c.Assert(err, gocheck.IsNil)
	t, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewCicle(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusBadRequest)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "text/plain; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestListCicles(c *gocheck.C) {
	r := `[{"id":"%v","name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]},{"id":"%v","name":"Test yakju","device":"yakju","packages":["bluetooth.jar","wifi.jar"]}]`
	cicle1 := test.Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	cicle2 := test.Cicle{Name: "Test yakju", Device: "yakju",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	test1, err := test.NewCicle(cicle1)
	c.Assert(err, gocheck.IsNil)
	test2, err := test.NewCicle(cicle2)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Cicle().RemoveId(test1.Id)
	defer db.Session.Cicle().RemoveId(test2.Id)
	request, err := http.NewRequest("GET", "cicles", nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	ListCicles(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	results := fmt.Sprintf(r, test1.Id.Hex(), test2.Id.Hex())
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, results)
}

func (s *S) TestDetailCicle(c *gocheck.C) {
	r := `{"id":"%v","name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]}`
	cicle := test.Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	testCicle, err := test.NewCicle(cicle)
	c.Assert(err, gocheck.IsNil)
	result := fmt.Sprintf(r, testCicle.Id.Hex())
	defer db.Session.Cicle().RemoveId(testCicle.Id)
	url := fmt.Sprintf("cicles/%s", testCicle.Id.Hex())
	request, err := http.NewRequest("GET", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	DetailCicle(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	ct := response.HeaderMap["Content-Type"][0]
	c.Assert(ct, gocheck.Equals, "application/json; charset=utf-8")
	c.Assert(response.Body.String(), gocheck.Equals, result)
}

func (s *S) TestModifyCicle(c *gocheck.C) {
	cicle := test.Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	testCicle, err := test.NewCicle(cicle)
	c.Assert(err, gocheck.IsNil)
	r := `{"id":"%v","name":"Test yakju","device":"yakju","packages":["bluetooth.jar","wifi.jar"]}`
	b := fmt.Sprintf(r, testCicle.Id.Hex())
	body := strings.NewReader(b)
	defer db.Session.Cicle().RemoveId(testCicle.Id)
	url := fmt.Sprintf("cicles/%s", testCicle.Id.Hex())
	var org test.Cicle
	err = json.Unmarshal([]byte(b), &org)
	c.Assert(err, gocheck.IsNil)
	request, err := http.NewRequest("PUT", url, body)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	ModifyCicle(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	var orgDb test.Cicle
	err = db.Session.Cicle().FindId(testCicle.Id).One(&orgDb)
	c.Assert(err, gocheck.IsNil)
	c.Assert(orgDb, gocheck.DeepEquals, org)
}

func (s *S) TestDeleteCicle(c *gocheck.C) {
	cicle := test.Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	testCicle, err := test.NewCicle(cicle)
	c.Assert(err, gocheck.IsNil)
	url := fmt.Sprintf("cicles/%s", testCicle.Id.Hex())
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	DeleteCicle(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusOK)
	lenght, err := db.Session.Cicle().FindId(testCicle.Id).Count()
	c.Assert(err, gocheck.IsNil)
	c.Assert(lenght, gocheck.Equals, 0)
}

func (s *S) TestDeleteCicleNoExist(c *gocheck.C) {
	result := "Error deleting cicle, Invalid cicle object id hex\n"
	cicle := test.Cicle{Name: "Test maguro", Device: "maguro",
		Packages: []string{"bluetooth.jar", "wifi.jar"}}
	testCicle, err := test.NewCicle(cicle)
	c.Assert(err, gocheck.IsNil)
	defer db.Session.Cicle().RemoveId(testCicle.Id)
	url := fmt.Sprintf("cicles/%s", "noexist")
	request, err := http.NewRequest("DELETE", url, nil)
	c.Assert(err, gocheck.IsNil)
	tkn, err := token.New("jhon@doe.com")
	c.Assert(err, gocheck.IsNil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tkn))
	tk, err := token.Validate(request)
	c.Assert(err, gocheck.IsNil)
	response := httptest.NewRecorder()
	DeleteCicle(response, request, tk)
	c.Assert(response.Code, gocheck.Equals, http.StatusBadRequest)
	c.Assert(response.Body.String(), gocheck.Equals, result)
}
