package internal

import (
	"context"

	"littlerollingsushi.com/example/entity"
)

type RegisterUsecase struct {
	config  RegisterUsecaseConfig
	gateway RegisterGateway
}

type RegisterUsecaseConfig struct {
	SaltLength int
}

//go:generate mockery --name=RegisterGateway --output=./mocks
type RegisterGateway interface {
	EncryptPassword(password string, saltLength int) (cryptedPassword string, err error)
	InsertUser(context.Context, entity.User) error
}

func NewRegisterUsecase(config RegisterUsecaseConfig, gateway RegisterGateway) *RegisterUsecase {
	return &RegisterUsecase{
		config:  config,
		gateway: gateway,
	}
}

func (u *RegisterUsecase) Register(ctx context.Context, in RegisterUsecaseInput) error {
	cryptedPassword, err := u.gateway.EncryptPassword(in.Password, u.config.SaltLength)
	if err != nil {
		return err
	}

	user := entity.User{
		FirstName:       in.FirstName,
		LastName:        in.LastName,
		Email:           in.Email,
		CryptedPassword: string(cryptedPassword),
	}

	return u.gateway.InsertUser(ctx, user)
}
