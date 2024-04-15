package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"github.com/sotavant/yandex-diplom-one/order"
	"io"
	"net/http"
)

type OrdersHandler struct {
	Service *order.Service
}

func NewOrdersHandler(r *chi.Mux, service *order.Service) {
	handler := &OrdersHandler{
		Service: service,
	}

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Route("/api/user/orders", func(r chi.Router) {
			r.Post("/", handler.AddOrder)
		})
	})
}

func (o *OrdersHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	if !isTextPlainRequest(r) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		internal.Logger.Infow("error in io.readAll", err, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = o.Service.Add(textBody)
}

func isTextPlainRequest(r *http.Request) bool {
	return render.GetRequestContentType(r) == render.ContentTypePlainText
}
