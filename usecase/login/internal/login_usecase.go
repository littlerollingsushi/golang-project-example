package internal

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"littlerollingsushi.com/example/entity"
)

const accessTokenExpirationDurationSeconds = 3600

//go:generate mockery --name=LoginGateway --output=./mocks
type LoginGateway interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	IsHashAndPasswordEqual(hash, password string) bool
	NowInUTC() time.Time
}

type LoginUsecase struct {
	gateway    LoginGateway
	privateKey *rsa.PrivateKey
}

func NewLoginUsecase(gateway LoginGateway, privateKey *rsa.PrivateKey) *LoginUsecase {
	return &LoginUsecase{gateway: gateway, privateKey: privateKey}
}

func (u *LoginUsecase) Login(ctx context.Context, in LoginUsecaseInput) (LoginUsecaseOutput, error) {
	if in.Email == "" {
		return LoginUsecaseOutput{}, ErrEmptyEmail
	}

	if in.Password == "" {
		return LoginUsecaseOutput{}, ErrEmptyPassword
	}

	user, err := u.gateway.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return LoginUsecaseOutput{}, err
	}

	if !u.gateway.IsHashAndPasswordEqual(user.CryptedPassword, in.Password) {
		return LoginUsecaseOutput{}, ErrInvalidPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, u.buildJwtClaim(user.Email))
	signedToken, err := token.SignedString(u.privateKey)
	if err != nil {
		return LoginUsecaseOutput{}, ErrInvalidPrivateKey
	}

	return LoginUsecaseOutput{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   accessTokenExpirationDurationSeconds,
	}, nil
}

func (u *LoginUsecase) buildJwtClaim(email string) *jwt.RegisteredClaims {
	now := u.gateway.NowInUTC()

	return &jwt.RegisteredClaims{
		Issuer:    "littlerollingsushi.com",
		Audience:  jwt.ClaimStrings{"littlerollingsushi.com"},
		Subject:   email,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenExpirationDurationSeconds * time.Second)),
	}
}
