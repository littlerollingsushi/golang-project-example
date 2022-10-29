package integration_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"littlerollingsushi.com/example/integration_test/helper"
)

type RegisterSuite struct {
	suite.Suite
}

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, &RegisterSuite{})
}

func (s *RegisterSuite) TestRegister_NewUser_ReturnCreated() {
	randomString, _ := helper.GenerateRandomString(31)
	form := url.Values{}
	form.Add("first_name", "john")
	form.Add("last_name", "doe")
	form.Add("email", randomString+"@email.com")
	form.Add("password", randomString)

	resp, respErr := http.Post("http://localhost:7070/v1/register", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	body, bodyErr := ioutil.ReadAll(resp.Body)
	unmarshalledBody := struct {
		Message string `json:"message"`
		Meta    struct {
			HttpStatus int       `json:"http_status"`
			ServerTime time.Time `json:"server_time"`
		}
	}{}
	unmarshallErr := json.Unmarshal(body, &unmarshalledBody)

	a := s.Assert()
	a.Nil(respErr)
	a.Nil(bodyErr)
	a.Nil(unmarshallErr)
	a.Equal("User registered. Continue to login.", unmarshalledBody.Message)
	a.Equal(http.StatusCreated, unmarshalledBody.Meta.HttpStatus)
}

func (s *RegisterSuite) TestRegister_DuplicateUser_ReturnUnprocessable() {
	randomString, _ := helper.GenerateRandomString(31)
	form := url.Values{}
	form.Add("first_name", "john")
	form.Add("last_name", "doe")
	form.Add("email", randomString+"@email.com")
	form.Add("password", randomString)

	// first post
	_, _ = http.Post("http://localhost:7070/v1/register", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

	// post using the same email
	resp, respErr := http.Post("http://localhost:7070/v1/register", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	body, bodyErr := ioutil.ReadAll(resp.Body)
	unmarshalledBody := struct {
		Message string `json:"message"`
		Meta    struct {
			HttpStatus int       `json:"http_status"`
			ServerTime time.Time `json:"server_time"`
		}
	}{}
	unmarshallErr := json.Unmarshal(body, &unmarshalledBody)

	a := s.Assert()
	a.Nil(respErr)
	a.Nil(bodyErr)
	a.Nil(unmarshallErr)
	a.Equal("User already exists. Choose different email.", unmarshalledBody.Message)
	a.Equal(http.StatusUnprocessableEntity, unmarshalledBody.Meta.HttpStatus)
}
