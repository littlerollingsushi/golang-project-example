package internal_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"littlerollingsushi.com/example/entity"
	"littlerollingsushi.com/example/usecase/login/internal"
	"littlerollingsushi.com/example/usecase/login/internal/mocks"
)

type LoginUsecaseSuite struct {
	suite.Suite

	context context.Context
	input   internal.LoginUsecaseInput
	output  internal.LoginUsecaseOutput

	priv    *rsa.PrivateKey
	gateway *mocks.LoginGateway
	usecase *internal.LoginUsecase

	user    entity.User
	now     time.Time
	errMock error
}

func TestLoginUsecaseSuite(t *testing.T) {
	suite.Run(t, &LoginUsecaseSuite{})
}

func (s *LoginUsecaseSuite) SetupTest() {
	s.context = context.Background()
	s.input = internal.LoginUsecaseInput{
		Email:    "john.doe@email.com",
		Password: "verysecure",
	}
	s.output = internal.LoginUsecaseOutput{
		TokenType: "Bearer",
		ExpiresIn: 3600,
	}

	s.priv, _ = rsa.GenerateKey(rand.Reader, 2048)
	s.gateway = mocks.NewLoginGateway(s.T())
	s.usecase = internal.NewLoginUsecase(s.gateway, s.priv)

	encrypted, _ := bcrypt.GenerateFromPassword([]byte(s.input.Password), bcrypt.DefaultCost)
	s.user = entity.User{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@email.com",
		CryptedPassword: string(encrypted),
	}
	s.now = time.Now()
	s.errMock = errors.New("mock error")
}

func (s *LoginUsecaseSuite) TestLogin_EmptyEmail_ReturnError() {
	s.input.Email = ""

	output, err := s.usecase.Login(s.context, s.input)

	a := s.Assert()
	a.Empty(output)
	a.ErrorIs(err, internal.ErrEmptyEmail)
}

func (s *LoginUsecaseSuite) TestLogin_EmptyPassowrd_ReturnError() {
	s.input.Password = ""

	output, err := s.usecase.Login(s.context, s.input)

	a := s.Assert()
	a.Empty(output)
	a.ErrorIs(err, internal.ErrEmptyPassword)
}

func (s *LoginUsecaseSuite) TestLogin_GetUserByEmailError_ReturnError() {
	s.gateway.On("GetUserByEmail", s.context, s.input.Email).Return(entity.User{}, s.errMock)

	output, err := s.usecase.Login(s.context, s.input)

	a := s.Assert()
	a.Empty(output)
	a.ErrorIs(err, s.errMock)
}

func (s *LoginUsecaseSuite) TestLogin_ComparePasswordError_ReturnError() {
	s.gateway.On("GetUserByEmail", s.context, s.input.Email).Return(s.user, nil)
	s.gateway.On("IsHashAndPasswordEqual", s.user.CryptedPassword, s.input.Password).Return(false)

	output, err := s.usecase.Login(s.context, s.input)

	a := s.Assert()
	a.Empty(output)
	a.ErrorIs(err, internal.ErrInvalidPassword)
}

func (s *LoginUsecaseSuite) TestLogin_ValidCredentialsInvalidPrivateKey_ReturnErrInvalidKey() {
	priv, _ := rsa.GenerateKey(strings.NewReader("random bytes."), 13)
	s.usecase = internal.NewLoginUsecase(s.gateway, priv)

	s.gateway.On("GetUserByEmail", s.context, s.input.Email).Return(s.user, nil)
	s.gateway.On("IsHashAndPasswordEqual", s.user.CryptedPassword, s.input.Password).Return(true)
	s.gateway.On("NowInUTC").Return(s.now)

	output, err := s.usecase.Login(s.context, s.input)

	a := s.Assert()
	a.Empty(output)
	a.ErrorIs(err, internal.ErrInvalidPrivateKey)
}

func (s *LoginUsecaseSuite) TestLogin_ValidCredentialsValidKey_ReturnAccessToken() {
	s.gateway.On("GetUserByEmail", s.context, s.input.Email).Return(s.user, nil)
	s.gateway.On("IsHashAndPasswordEqual", s.user.CryptedPassword, s.input.Password).Return(true)
	s.gateway.On("NowInUTC").Return(s.now)

	output, err := s.usecase.Login(s.context, s.input)

	a := s.Assert()
	a.Nil(err)
	parsed, err := jwt.ParseWithClaims(output.AccessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return &s.priv.PublicKey, nil
	})
	a.Nil(err)
	claims := parsed.Claims.(*jwt.RegisteredClaims)
	a.Nil(claims.Valid())
	a.Equal("littlerollingsushi.com", claims.Issuer)
	a.Equal(jwt.ClaimStrings{"littlerollingsushi.com"}, claims.Audience)
	a.Equal(s.input.Email, claims.Subject)
	a.Equal(time.Unix(s.now.Unix(), 0), claims.NotBefore.Time)
	a.Equal(time.Unix(s.now.Unix(), 0), claims.IssuedAt.Time)
	a.Equal(time.Unix(s.now.Unix(), 0).Add(1*time.Hour), claims.ExpiresAt.Time)
	a.Equal(s.output.ExpiresIn, output.ExpiresIn)
	a.Equal(s.output.TokenType, output.TokenType)
}
