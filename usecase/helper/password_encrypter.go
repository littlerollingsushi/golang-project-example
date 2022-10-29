package helper

import "golang.org/x/crypto/bcrypt"

type PasswordEncrypter struct {
}

func (*PasswordEncrypter) EncryptPassword(password string, saltLength int) (cryptedPassword string, err error) {
	crypted, err := bcrypt.GenerateFromPassword([]byte(password), saltLength)
	return string(crypted), err
}
