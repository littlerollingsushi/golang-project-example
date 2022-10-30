package integration_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"littlerollingsushi.com/example/integration_test/helper"
)

type LoginSuite struct {
	suite.Suite
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, &LoginSuite{})
}

func (s *LoginSuite) TestLogin_RegisteredUser_ReturnAccessToken() {
	randomString, _ := helper.GenerateRandomString(31)
	form := url.Values{}
	form.Add("first_name", "john")
	form.Add("last_name", "doe")
	form.Add("email", randomString+"@email.com")
	form.Add("password", randomString)

	_, err := http.Post("http://localhost:7070/v1/register", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalf("Error on registering user on login integration test: %v\n", err)
	}
	resp, respErr := http.Post("http://localhost:7070/v1/login", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	body, bodyErr := io.ReadAll(resp.Body)
	unmarshalledBody := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Meta        struct {
			HttpStatus int       `json:"http_status"`
			ServerTime time.Time `json:"server_time"`
		}
	}{}
	unmarshallErr := json.Unmarshal(body, &unmarshalledBody)

	a := s.Assert()
	a.Nil(respErr)
	a.Nil(bodyErr)
	a.Nil(unmarshallErr)
	a.True(len(unmarshalledBody.AccessToken) != 0)
	a.Equal(3600, unmarshalledBody.ExpiresIn)
	a.Equal("Bearer", unmarshalledBody.TokenType)
	a.Equal(http.StatusOK, unmarshalledBody.Meta.HttpStatus)
}

func (s *LoginSuite) TestLogin_UnregisteredUser_ReturnCreated() {
	randomString, _ := helper.GenerateRandomString(31)
	form := url.Values{}
	form.Add("first_name", "john")
	form.Add("last_name", "doe")
	form.Add("email", randomString+"@email.com")
	form.Add("password", randomString)

	resp, respErr := http.Post("http://localhost:7070/v1/login", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	body, bodyErr := io.ReadAll(resp.Body)
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
	a.Equal("Invalid credentials.", unmarshalledBody.Message)
	a.Equal(http.StatusUnauthorized, unmarshalledBody.Meta.HttpStatus)
}
