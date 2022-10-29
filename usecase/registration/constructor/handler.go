package constructor

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"littlerollingsushi.com/example/usecase/helper"
	"littlerollingsushi.com/example/usecase/registration/handler"
	"littlerollingsushi.com/example/usecase/registration/internal"
)

func ConstructRegisterHandler(db *sql.DB) *handler.RegisterHandler {
	gateway := internal.NewInsertUserGateway(db)
	usecase := internal.NewRegisterUsecase(
		internal.RegisterUsecaseConfig{SaltLength: bcrypt.DefaultCost},
		struct {
			*internal.InsertUserGateway
			*helper.PasswordEncrypter
		}{
			InsertUserGateway: gateway,
			PasswordEncrypter: &helper.PasswordEncrypter{},
		},
	)
	timer := &helper.TimerImplementation{}
	return handler.NewRegisterHandler(usecase, timer)
}
