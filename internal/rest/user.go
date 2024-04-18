package rest

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
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

	r.Group(func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Route("/api/user", func(r chi.Router) {
			r.Post(fmt.Sprintf("/%s", registerURI), handler.Auth)
			r.Post(fmt.Sprintf("/%s", loginURI), handler.Auth)
			r.With(middleware.Auth).Get("/balance", handler.GetBalance)
		})
	})
}

func (u *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	userReq := &userRequest{}
	if err := render.Bind(r, userReq); err != nil {
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
		token, err = u.Service.Register(r.Context(), *userReq.User)
	} else {
		token, err = u.Service.Login(r.Context(), *userReq.User)
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

func (u *UserHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(user.ContextUserIdKey).(int64)
	if userId == 0 {
		err := render.Render(w, r, errorRender(getStatusCode(domain.ErrUserNotAuthorized), domain.ErrUserNotAuthorized))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}

	dbUser, err := u.Service.GetById(r.Context(), userId)
	if err != nil {
		err = render.Render(w, r, errorRender(getStatusCode(err), err))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}

	if err = render.Render(w, r, newUserBalanceResponse(&dbUser)); err != nil {
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

type userBalanceResponse struct {
	*domain.User
}

func newUserBalanceResponse(user *domain.User) *userBalanceResponse {
	resp := &userBalanceResponse{user}

	return resp
}

func (ub *userBalanceResponse) Render(w http.ResponseWriter, r *http.Request) error {
	ub.Login = ""
	ub.Password = ""

	return nil
}

func isRegisterPage(url string) bool {
	return strings.Contains(url, registerURI)
}
