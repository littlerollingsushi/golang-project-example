package internal_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"littlerollingsushi.com/example/entity"
	"littlerollingsushi.com/example/usecase/registration/internal"
	"littlerollingsushi.com/example/usecase/registration/internal/mocks"
)

type RegisterUsecaseSuite struct {
	suite.Suite

	config  internal.RegisterUsecaseConfig
	gateway *mocks.RegisterGateway
	usecase *internal.RegisterUsecase

	context                context.Context
	input                  internal.RegisterUsecaseInput
	cryptedPassword        string
	errMock                error
	expectedInsertUserData entity.User
}

func TestRegisterUsecaseSuite(t *testing.T) {
	suite.Run(t, &RegisterUsecaseSuite{})
}

func (s *RegisterUsecaseSuite) SetupTest() {
	s.config = internal.RegisterUsecaseConfig{SaltLength: 12}
	s.gateway = mocks.NewRegisterGateway(s.T())
	s.usecase = internal.NewRegisterUsecase(s.config, s.gateway)

	s.context = context.Background()
	s.input = internal.RegisterUsecaseInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@email.com",
		Password:  "verysecure",
	}
	s.cryptedPassword = "verysecureencrypted"
	s.errMock = errors.New("mock error")
	s.expectedInsertUserData = entity.User{
		FirstName:       s.input.FirstName,
		LastName:        s.input.LastName,
		Email:           s.input.Email,
		CryptedPassword: s.cryptedPassword,
	}
}

func (s *RegisterUsecaseSuite) TestRegister_GeneratePasswordFailed_ReturnOriginalError() {
	s.gateway.On("EncryptPassword", s.input.Password, s.config.SaltLength).Return("", s.errMock)

	err := s.usecase.Register(s.context, s.input)

	s.Assert().ErrorIs(err, s.errMock)
}

func (s *RegisterUsecaseSuite) TestRegister_InsertUserFailed_ReturnOriginalError() {
	s.gateway.On("EncryptPassword", s.input.Password, s.config.SaltLength).Return(s.cryptedPassword, nil)
	s.gateway.On("InsertUser", s.context, s.expectedInsertUserData).Return(s.errMock)

	err := s.usecase.Register(s.context, s.input)

	s.Assert().ErrorIs(err, s.errMock)
}

func (s *RegisterUsecaseSuite) TestRegister_InsertUserSuccess_ReturnNil() {
	s.gateway.On("EncryptPassword", s.input.Password, s.config.SaltLength).Return(s.cryptedPassword, nil)
	s.gateway.On("InsertUser", s.context, s.expectedInsertUserData).Return(nil)

	err := s.usecase.Register(s.context, s.input)

	s.Assert().Nil(err)
}
