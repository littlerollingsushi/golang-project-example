package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"littlerollingsushi.com/example/usecase/helper"
	"littlerollingsushi.com/example/usecase/registration/internal"
)

type RegisterHandler struct {
	usecase RegisterUsecase
	timer   helper.Timer
}

//go:generate mockery --name=RegisterUsecase --output=./mocks
type RegisterUsecase interface {
	Register(context.Context, internal.RegisterUsecaseInput) error
}

func NewRegisterHandler(usecase RegisterUsecase, timer helper.Timer) *RegisterHandler {
	return &RegisterHandler{usecase: usecase, timer: timer}
}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	in := internal.RegisterUsecaseInput{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
	}

	err := h.usecase.Register(r.Context(), in)
	if err != nil {
		h.processError(w, err)
		return
	}

	h.writeRegisterResponse(w)
}

func (h *RegisterHandler) processError(w http.ResponseWriter, err error) {
	switch err {
	case internal.ErrUserAlreadyExist:
		h.processUserAlreadyExistErr(w, err)
	default:
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Oops! Something went wrong."))
	}
}

func (h *RegisterHandler) processUserAlreadyExistErr(w http.ResponseWriter, err error) {
	data := map[string]interface{}{
		"message": "User already exists. Choose different email.",
		"meta": map[string]interface{}{
			"http_status": http.StatusUnprocessableEntity,
			"server_time": h.timer.NowInUTC().Format("2006-01-02T15:04:05.999Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(data)
}

func (h *RegisterHandler) writeRegisterResponse(w http.ResponseWriter) {
	data := map[string]interface{}{
		"message": "User registered. Continue to login.",
		"meta": map[string]interface{}{
			"http_status": http.StatusCreated,
			"server_time": h.timer.NowInUTC().Format("2006-01-02T15:04:05.999Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
