package internal_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"littlerollingsushi.com/example/entity"
	"littlerollingsushi.com/example/usecase/login/internal"
)

type GetUserByEmailGatewaySuite struct {
	suite.Suite

	db            *sql.DB
	mockDb        sqlmock.Sqlmock
	expectedQuery string
	errMock       error

	context context.Context
	email   string
	user    entity.User
	gateway *internal.GetUserByEmailGateway
}

func TestGetUserByEmailGatewaySuite(t *testing.T) {
	suite.Run(t, &GetUserByEmailGatewaySuite{})
}

func (s *GetUserByEmailGatewaySuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error occured on opening a stub database: %v\n", err)
	}

	s.db = db
	s.mockDb = mock
	s.expectedQuery = "SELECT first_name, last_name, email, crypted_password FROM user WHERE email = ?"
	s.errMock = errors.New("mocked error")

	s.gateway = internal.NewGetUserByEmailGateway(s.db)
	s.context = context.Background()
	s.email = "john.doe@email.com"
	s.user = entity.User{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           s.email,
		CryptedPassword: "verysecureencrypted",
	}
}

func (s *GetUserByEmailGatewaySuite) TearDownTest() {
	s.db.Close()
}

func (s *GetUserByEmailGatewaySuite) TestGetUserByEmail_UnknownError_ReturnOriginalError() {
	s.mockDb.ExpectQuery(regexp.QuoteMeta(s.expectedQuery)).WillReturnError(s.errMock)

	user, err := s.gateway.GetUserByEmail(s.context, s.email)

	a := s.Assert()
	a.Empty(user)
	a.ErrorIs(err, s.errMock)
}

func (s *GetUserByEmailGatewaySuite) TestGetUserByEmail_DuplicateUser_ReturnUserAlreadyExistErr() {
	s.mockDb.ExpectQuery(regexp.QuoteMeta(s.expectedQuery)).WillReturnError(sql.ErrNoRows)

	user, err := s.gateway.GetUserByEmail(s.context, s.email)

	a := s.Assert()
	a.Empty(user)
	a.ErrorIs(err, internal.ErrUserNotFound)
}

func (s *GetUserByEmailGatewaySuite) TestGetUserByEmail_InsertSuccess_ReturnNil() {
	rows := sqlmock.NewRows([]string{"first_name", "last_name", "email", "crypted_password"})
	rows.AddRow(s.user.FirstName, s.user.LastName, s.user.Email, s.user.CryptedPassword)
	s.mockDb.ExpectQuery(regexp.QuoteMeta(s.expectedQuery)).WillReturnRows(rows)

	user, err := s.gateway.GetUserByEmail(s.context, s.email)

	a := s.Assert()
	a.EqualValues(s.user, user)
	a.Nil(err)
}
