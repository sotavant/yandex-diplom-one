package rest

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"github.com/sotavant/yandex-diplom-one/user"
	"github.com/sotavant/yandex-diplom-one/withdrawn"
	"net/http"
)

type WithdrawnHandler struct {
	Service *withdrawn.Service
}

type withdrawnRequest struct {
	*domain.Withdrawn
}

type WithdrawnResponse struct {
	domain.Withdrawn
}

func NewWithdrawnHandler(r *chi.Mux, service *withdrawn.Service) {
	handler := &WithdrawnHandler{
		Service: service,
	}

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post("/api/user/balance/withdraw", handler.Add)
		r.Get("/api/user/withdrawals", handler.List)
	})
}

func (wd *WithdrawnHandler) Add(w http.ResponseWriter, r *http.Request) {
	wdReq := &withdrawnRequest{}
	if err := render.Bind(r, wdReq); err != nil {
		err = render.Render(w, r, errorRender(http.StatusBadRequest, err))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}

	userID := r.Context().Value(user.ContextUserIDKey{}).(int64)
	if userID == 0 {
		err := render.Render(w, r, errorRender(getStatusCode(domain.ErrUserNotAuthorized), domain.ErrUserNotAuthorized))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}

	wdReq.Withdrawn.UserID = userID
	err := wd.Service.Add(r.Context(), wdReq.Withdrawn)
	if err != nil {
		err = render.Render(w, r, errorRender(getStatusCode(err), err))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}
}

func (wd *WithdrawnHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(user.ContextUserIDKey{}).(int64)
	if userID == 0 {
		err := render.Render(w, r, errorRender(getStatusCode(domain.ErrUserNotAuthorized), domain.ErrUserNotAuthorized))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}

	wds, response, err := wd.Service.List(r.Context(), userID)
	if err != nil {
		err = render.Render(w, r, errorRender(http.StatusInternalServerError, err))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			internal.Logger.Infoln(err)
		}
		return
	}

	if response != "" {
		w.WriteHeader(getResponseCode(response))
		return
	}

	if err = render.RenderList(w, r, NewWithdrawnListResponse(wds)); err != nil {
		render.Status(r, http.StatusInternalServerError)
		internal.Logger.Infoln(err)
		return
	}
}

func (u *withdrawnRequest) Bind(r *http.Request) error {
	if u.Withdrawn == nil {
		return errors.New("withdrawn fields absent")
	}

	if u.Withdrawn.OrderNum == "" || u.Withdrawn.Sum <= 0 {
		return errors.New("bad params")
	}

	return nil
}

func NewWithdrawnResponse(wd domain.Withdrawn) *WithdrawnResponse {
	resp := &WithdrawnResponse{Withdrawn: wd}

	return resp
}

func (wd *WithdrawnResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewWithdrawnListResponse(wds []domain.Withdrawn) []render.Renderer {
	var list []render.Renderer
	for _, v := range wds {
		list = append(list, NewWithdrawnResponse(v))
	}

	return list
}
