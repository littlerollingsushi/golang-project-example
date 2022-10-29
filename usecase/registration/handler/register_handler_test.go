package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	helperMocks "littlerollingsushi.com/example/usecase/helper/mocks"
	"littlerollingsushi.com/example/usecase/registration/handler"
	"littlerollingsushi.com/example/usecase/registration/handler/mocks"
	"littlerollingsushi.com/example/usecase/registration/internal"
)

type RegisterHandlerSuite struct {
	suite.Suite

	request        *http.Request
	requestParams  map[string]string
	responseWriter *httptest.ResponseRecorder

	usecase *mocks.RegisterUsecase
	timer   *helperMocks.Timer
	handler *handler.RegisterHandler

	expectedUsecaseInput              internal.RegisterUsecaseInput
	expectedTimestamp                 time.Time
	expectedSuccessResponseBody       string
	expectedDuplicateUserResponseBody string
	errMock                           error
}

func TestRegisterHandlerSuite(t *testing.T) {
	suite.Run(t, &RegisterHandlerSuite{})
}

func (s *RegisterHandlerSuite) SetupTest() {
	s.request = httptest.NewRequest("POST", "http://test.com/v1/register", strings.NewReader("first_name=john&last_name=doe&email=john.doe@email.com&password=verysecure"))
	s.request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.responseWriter = httptest.NewRecorder()

	s.requestParams = map[string]string{}

	s.usecase = mocks.NewRegisterUsecase(s.T())
	s.timer = helperMocks.NewTimer(s.T())
	s.handler = handler.NewRegisterHandler(s.usecase, s.timer)

	s.expectedUsecaseInput = internal.RegisterUsecaseInput{
		FirstName: "john",
		LastName:  "doe",
		Email:     "john.doe@email.com",
		Password:  "verysecure",
	}
	s.expectedTimestamp = time.Date(2022, 10, 29, 23, 59, 59, 123000000, time.UTC)
	s.expectedSuccessResponseBody = `
		{
			"message": "User registered. Continue to login.",
			"meta": {
				"http_status": 201,
				"server_time": "2022-10-29T23:59:59.123Z"
			}
		}
	`
	s.expectedDuplicateUserResponseBody = `
		{
			"message": "User already exists. Choose different email.",
			"meta": {
				"http_status": 422,
				"server_time": "2022-10-29T23:59:59.123Z"
			}
		}
	`

	s.errMock = errors.New("mock error")
}

func (s *RegisterHandlerSuite) TestRegister_UsecaseUnknownError_ReturnInternalServerError() {
	s.usecase.On("Register", s.request.Context(), s.expectedUsecaseInput).Return(s.errMock)

	s.handler.Register(s.responseWriter, s.request, s.requestParams)

	resp := s.responseWriter.Result()
	body, _ := io.ReadAll(resp.Body)
	a := s.Assert()
	a.Equal(http.StatusInternalServerError, resp.StatusCode)
	a.Equal("Oops! Something went wrong.", string(body))
}

func (s *RegisterHandlerSuite) TestRegister_DuplicateUser_ReturnInternalServerError() {
	s.usecase.On("Register", s.request.Context(), s.expectedUsecaseInput).Return(internal.ErrUserAlreadyExist)
	s.timer.On("NowInUTC").Return(s.expectedTimestamp)

	s.handler.Register(s.responseWriter, s.request, s.requestParams)

	resp := s.responseWriter.Result()
	body, _ := io.ReadAll(resp.Body)
	a := s.Assert()
	a.Equal(http.StatusUnprocessableEntity, resp.StatusCode)
	a.JSONEq(s.expectedDuplicateUserResponseBody, string(body))
}

func (s *RegisterHandlerSuite) TestRegister_UsecaseSuccess_ReturnCreated() {
	s.usecase.On("Register", s.request.Context(), s.expectedUsecaseInput).Return(nil)
	s.timer.On("NowInUTC").Return(s.expectedTimestamp)

	s.handler.Register(s.responseWriter, s.request, s.requestParams)

	resp := s.responseWriter.Result()
	body, _ := io.ReadAll(resp.Body)
	a := s.Assert()
	a.Equal(http.StatusCreated, resp.StatusCode)
	a.JSONEq(s.expectedSuccessResponseBody, string(body))
}
