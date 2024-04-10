package rest

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
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
		render.Status(r, http.StatusBadRequest)
		internal.Logger.Infoln(err)
		return
	}

	token, err := u.Service.Register(r.Context(), *user.User)
	if err != nil {
		err = render.Render(w, r, errorRender(getStatusCode(err), err))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
			return
		}

		return
	}

	if err = render.Render(w, r, newTokenResponse(token)); err != nil {
		render.Status(r, http.StatusInternalServerError)
		internal.Logger.Infoln(err)
		return
	}
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

type tokenResponse struct {
	token string
}

func newTokenResponse(token string) *tokenResponse {
	return &tokenResponse{
		token: token,
	}
}

func (t *tokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrBadParams, domain.ErrPasswordTooWeak:
		return http.StatusBadRequest
	case domain.ErrLoginExist:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
