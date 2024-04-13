package rest

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"net/http"
)

type OrdersHandler struct{}

func NewOrdersHandler(r *chi.Mux) {
	handler := &OrdersHandler{}

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
	fmt.Println("ok")
}

func isTextPlainRequest(r *http.Request) bool {
	return render.GetRequestContentType(r) == render.ContentTypePlainText
}
