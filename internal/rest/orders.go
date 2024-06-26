package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"github.com/sotavant/yandex-diplom-one/order"
	"io"
	"net/http"
)

type OrdersHandler struct {
	Service *order.Service
}

type OrderResponse struct {
	domain.Order
}

func NewOrdersHandler(r *chi.Mux, service *order.Service) {
	handler := &OrdersHandler{
		Service: service,
	}

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Route("/api/user/orders", func(r chi.Router) {
			r.Post("/", handler.AddOrder)
			r.With(render.SetContentType(render.ContentTypeJSON)).Get("/", handler.Orders)
		})
	})
}

func NewOrderResponse(order domain.Order) *OrderResponse {
	resp := &OrderResponse{Order: order}

	return resp
}

func (rd *OrderResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewOrderListResponse(orders []domain.Order) []render.Renderer {
	var list []render.Renderer
	for _, v := range orders {
		list = append(list, NewOrderResponse(v))
	}

	return list
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

	response, err := o.Service.Add(r.Context(), textBody)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
	}

	responseCode := getResponseCode(response)
	w.WriteHeader(responseCode)
	_, err = w.Write([]byte(response))
	if err != nil {
		internal.Logger.Infow("error in write response", err, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (o *OrdersHandler) Orders(w http.ResponseWriter, r *http.Request) {
	orders, response, err := o.Service.List(r.Context())
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

	if err = render.RenderList(w, r, NewOrderListResponse(orders)); err != nil {
		render.Status(r, http.StatusInternalServerError)
		internal.Logger.Infoln(err)
		return
	}
}

func isTextPlainRequest(r *http.Request) bool {
	return render.GetRequestContentType(r) == render.ContentTypePlainText
}

func getResponseCode(response string) int {
	switch response {
	case domain.RespOrderAlreadyUploaded:
		return http.StatusOK
	case domain.RespOrderAdmitted:
		return http.StatusAccepted
	case domain.RespNoDataToResponse:
		return http.StatusNoContent
	default:
		return http.StatusOK
	}
}
