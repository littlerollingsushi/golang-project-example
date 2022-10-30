package constructor

import (
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"littlerollingsushi.com/example/usecase/helper"
	"littlerollingsushi.com/example/usecase/login/handler"
	"littlerollingsushi.com/example/usecase/login/internal"
)

type Config struct {
	PrivateKeyPath string `envconfig:"PRIVATE_KEY_PATH"`
}

func ConstructLoginHandler(db *sql.DB) *handler.LoginHandler {
	cfg := Config{}
	envconfig.Process("RSA", &cfg)

	rawPrivateKey, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		log.Fatalf("invalid raw private key PEM path: %v\n", err)
	}

	privPem, _ := pem.Decode(rawPrivateKey)
	privKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		log.Fatalf("invalid PCKS1 private key: %v\n", err)
	}

	gateway := internal.NewGetUserByEmailGateway(db)
	usecase := internal.NewLoginUsecase(
		struct {
			*internal.GetUserByEmailGateway
			*helper.PasswordEncrypter
			helper.Timer
		}{
			GetUserByEmailGateway: gateway,
			PasswordEncrypter:     &helper.PasswordEncrypter{},
			Timer:                 &helper.TimerImplementation{},
		},
		privKey,
	)
	timer := &helper.TimerImplementation{}
	return handler.NewLoginHandler(usecase, timer)
}
