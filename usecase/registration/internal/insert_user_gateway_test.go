package internal_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"
	"littlerollingsushi.com/example/entity"
	"littlerollingsushi.com/example/usecase/registration/internal"
)

type InsertUserGatewaySuite struct {
	suite.Suite

	db                 *sql.DB
	mockDb             sqlmock.Sqlmock
	expectedQuery      string
	errMock            error
	errDuplicateRecord error

	context context.Context
	input   entity.User
	gateway *internal.InsertUserGateway
}

func TestInsertUserGatewaySuite(t *testing.T) {
	suite.Run(t, &InsertUserGatewaySuite{})
}

func (s *InsertUserGatewaySuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error occured on opening a stub database: %v\n", err)
	}

	s.db = db
	s.mockDb = mock
	s.expectedQuery = "INSERT INTO user (first_name, last_name, email, crypted_password, created_at) VALUES (?, ?, ?, ?, ?)"
	s.errMock = errors.New("mocked error")
	s.errDuplicateRecord = &mysql.MySQLError{Number: 1062, Message: "mock message"}

	s.gateway = internal.NewInsertUserGateway(s.db)
	s.context = context.Background()
	s.input = entity.User{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@email.com",
		CryptedPassword: "verysecureencrypted",
	}
}

func (s *InsertUserGatewaySuite) TearDownTest() {
	s.db.Close()
}

func (s *InsertUserGatewaySuite) TestInsertUser_UnknownError_ReturnOriginalError() {
	s.mockDb.ExpectExec(regexp.QuoteMeta(s.expectedQuery)).WillReturnError(s.errMock)

	err := s.gateway.InsertUser(s.context, s.input)

	s.Assert().ErrorIs(err, s.errMock)
}

func (s *InsertUserGatewaySuite) TestInsertUser_DuplicateUser_ReturnUserAlreadyExistErr() {
	s.mockDb.ExpectExec(regexp.QuoteMeta(s.expectedQuery)).WillReturnError(s.errDuplicateRecord)

	err := s.gateway.InsertUser(s.context, s.input)

	s.Assert().ErrorIs(err, internal.ErrUserAlreadyExist)
}

func (s *InsertUserGatewaySuite) TestInsertUser_InsertSuccess_ReturnNil() {
	s.mockDb.ExpectExec(regexp.QuoteMeta(s.expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.gateway.InsertUser(s.context, s.input)

	s.Assert().Nil(err)
}
