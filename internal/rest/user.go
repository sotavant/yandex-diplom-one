package rest

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/user"
	"net/http"
	"strings"
)

type UserHandler struct {
	Service *user.Service
}

const (
	registerURI = "register"
	loginURI    = "login"
)

func NewUserHandler(r *chi.Mux, service *user.Service) {
	handler := &UserHandler{
		Service: service,
	}

	r.Route("/api/user", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post(fmt.Sprintf("/%s", registerURI), handler.Auth)
		r.Post(fmt.Sprintf("/%s", loginURI), handler.Auth)
	})
}

func (u *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	userRequest := &userRequest{}
	if err := render.Bind(r, userRequest); err != nil {
		err = render.Render(w, r, errorRender(http.StatusBadRequest, err))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}

	var token string
	var err error

	if isRegisterPage(r.RequestURI) {
		token, err = u.Service.Register(r.Context(), *userRequest.User)
	} else {
		token, err = u.Service.Login(r.Context(), *userRequest.User)
	}

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

	if u.User.Password == "" || u.User.Login == "" {
		return errors.New("bad params")
	}

	return nil
}

type tokenResponse struct {
	Token string `json:"token"`
}

func newTokenResponse(token string) *tokenResponse {
	return &tokenResponse{
		Token: token,
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
	case domain.ErrBadUserData:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func isRegisterPage(url string) bool {
	return strings.Contains(url, registerURI)
}
