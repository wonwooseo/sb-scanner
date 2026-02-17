package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code" example:"0000"`
	Message string `json:"message" example:"internal server error"`
}

func (res *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)
	return nil
}

func MakeInternalServerError() *ErrResponse {
	return &ErrResponse{http.StatusInternalServerError, "0000", "internal server error"}
}

func MakeBadRequestError(msg string) *ErrResponse {
	return &ErrResponse{http.StatusBadRequest, "0001", msg}
}
