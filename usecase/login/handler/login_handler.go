package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"littlerollingsushi.com/example/usecase/helper"
	"littlerollingsushi.com/example/usecase/login/internal"
)

type LoginHandler struct {
	usecase LoginUsecase
	timer   helper.Timer
}

//go:generate mockery --name=LoginUsecase --output=./mocks
type LoginUsecase interface {
	Login(ctx context.Context, in internal.LoginUsecaseInput) (internal.LoginUsecaseOutput, error)
}

func NewLoginHandler(usecase LoginUsecase, timer helper.Timer) *LoginHandler {
	return &LoginHandler{usecase: usecase, timer: timer}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	in := internal.LoginUsecaseInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	out, err := h.usecase.Login(r.Context(), in)
	if err != nil {
		h.processError(w, err)
		return
	}

	h.writeLoginResponse(w, out)
}

func (h *LoginHandler) processError(w http.ResponseWriter, err error) {
	switch err {
	case internal.ErrUserNotFound:
		h.processUserNotFoundError(w, err)
	default:
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Oops! Something went wrong."))
	}
}

func (h *LoginHandler) processUserNotFoundError(w http.ResponseWriter, err error) {
	data := map[string]interface{}{
		"message": "Invalid credentials.",
		"meta": map[string]interface{}{
			"http_status": http.StatusUnauthorized,
			"server_time": h.timer.NowInUTC().Format("2006-01-02T15:04:05.999Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(data)
}

func (h *LoginHandler) writeLoginResponse(w http.ResponseWriter, out internal.LoginUsecaseOutput) {
	data := map[string]interface{}{
		"access_token": out.AccessToken,
		"expires_in":   out.ExpiresIn,
		"token_type":   out.TokenType,
		"meta": map[string]interface{}{
			"http_status": http.StatusOK,
			"server_time": h.timer.NowInUTC().Format("2006-01-02T15:04:05.999Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
