package rest

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/user"
	"net/http"
)

type UserHandler struct {
	Service *user.Service
}

func NewUserHandler(r *chi.Mux, service *user.Service) {
	handler := &UserHandler{
		Service: service,
	}

	r.Route("/api/user", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post("/register", handler.Register)
	})
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}
	if err := render.Bind(r, user); err != nil {

	}
	err := u.Service.Register()
	render.Status(r, http.StatusOK)
}

type userRequest struct {
	*domain.User
}

func (u *userRequest) Bind(r *http.Request) error {
	if u.User == nil {
		return errors.New("user fields absent")
	}

	return nil
}
