package internal

type LoginUsecaseInput struct {
	Email    string
	Password string
}

type LoginUsecaseOutput struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
}
