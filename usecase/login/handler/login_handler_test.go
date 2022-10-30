package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	helperMocks "littlerollingsushi.com/example/usecase/helper/mocks"
	"littlerollingsushi.com/example/usecase/login/handler"
	"littlerollingsushi.com/example/usecase/login/handler/mocks"
	"littlerollingsushi.com/example/usecase/login/internal"
)

type LoginHandlerSuite struct {
	suite.Suite

	request        *http.Request
	requestParams  map[string]string
	responseWriter *httptest.ResponseRecorder

	usecase *mocks.LoginUsecase
	timer   *helperMocks.Timer
	handler *handler.LoginHandler

	expectedUsecaseInput         internal.LoginUsecaseInput
	expectedUsecaseOutput        internal.LoginUsecaseOutput
	expectedTimestamp            time.Time
	expectedSuccessResponseBody  string
	expectedUserNotFoundResponse string
	errMock                      error
}

func TestLoginHandlerSuite(t *testing.T) {
	suite.Run(t, &LoginHandlerSuite{})
}

func (s *LoginHandlerSuite) SetupTest() {
	form := url.Values{}
	form.Add("email", "john.doe@email.com")
	form.Add("password", "verysecure")
	s.request = httptest.NewRequest("POST", "http://test.com/v1/login", strings.NewReader(form.Encode()))
	s.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.responseWriter = httptest.NewRecorder()

	s.requestParams = map[string]string{}

	s.usecase = mocks.NewLoginUsecase(s.T())
	s.timer = helperMocks.NewTimer(s.T())
	s.handler = handler.NewLoginHandler(s.usecase, s.timer)

	s.expectedUsecaseInput = internal.LoginUsecaseInput{
		Email:    "john.doe@email.com",
		Password: "verysecure",
	}
	s.expectedUsecaseOutput = internal.LoginUsecaseOutput{
		AccessToken: "very secure access token",
		ExpiresIn:   3600,
		TokenType:   "Bearer",
	}
	s.expectedTimestamp = time.Date(2022, 10, 29, 23, 59, 59, 123000000, time.UTC)
	s.expectedSuccessResponseBody = `
		{
			"access_token": "very secure access token",
			"expires_in": 3600,
			"token_type": "Bearer",
			"meta": {
				"http_status": 200,
				"server_time": "2022-10-29T23:59:59.123Z"
			}
		}
	`
	s.expectedUserNotFoundResponse = `
		{
			"message": "Invalid credentials.",
			"meta": {
				"http_status": 401,
				"server_time": "2022-10-29T23:59:59.123Z"
			}
		}
	`

	s.errMock = errors.New("mock error")
}

func (s *LoginHandlerSuite) TestLogin_UsecaseUnknownError_ReturnInternalServerError() {
	s.usecase.On("Login", s.request.Context(), s.expectedUsecaseInput).Return(internal.LoginUsecaseOutput{}, s.errMock)

	s.handler.Login(s.responseWriter, s.request, s.requestParams)

	resp := s.responseWriter.Result()
	body, _ := io.ReadAll(resp.Body)
	a := s.Assert()
	a.Equal(http.StatusInternalServerError, resp.StatusCode)
	a.Equal("Oops! Something went wrong.", string(body))
}

func (s *LoginHandlerSuite) TestLogin_UserNotFound_ReturnInternalServerError() {
	s.usecase.On("Login", s.request.Context(), s.expectedUsecaseInput).Return(internal.LoginUsecaseOutput{}, internal.ErrUserNotFound)
	s.timer.On("NowInUTC").Return(s.expectedTimestamp)

	s.handler.Login(s.responseWriter, s.request, s.requestParams)

	resp := s.responseWriter.Result()
	body, _ := io.ReadAll(resp.Body)
	a := s.Assert()
	a.Equal(http.StatusUnauthorized, resp.StatusCode)
	a.JSONEq(s.expectedUserNotFoundResponse, string(body))
}

func (s *LoginHandlerSuite) TestLogin_UsecaseSuccess_ReturnCreated() {
	s.usecase.On("Login", s.request.Context(), s.expectedUsecaseInput).Return(s.expectedUsecaseOutput, nil)
	s.timer.On("NowInUTC").Return(s.expectedTimestamp)

	s.handler.Login(s.responseWriter, s.request, s.requestParams)

	resp := s.responseWriter.Result()
	body, _ := io.ReadAll(resp.Body)
	a := s.Assert()
	a.Equal(http.StatusOK, resp.StatusCode)
	a.JSONEq(s.expectedSuccessResponseBody, string(body))
}
