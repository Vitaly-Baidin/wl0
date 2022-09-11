package route

import (
	"context"
	"errors"
	"github.com/Vitaly-Baidin/l0/internal/sub/route/httpErr"
	"github.com/Vitaly-Baidin/l0/internal/sub/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"net/http"
)

type route struct {
	Service *service.OrderService
}

func NewOrderRoute(service *service.OrderService) *route {
	return &route{
		Service: service,
	}
}

func (r *route) RegisterRoute(router *chi.Mux) {
	router.Route("/api/v1/orders", func(router chi.Router) {
		router.Get("/", r.getAllOrders)
		router.Get("/{id}", r.getOrderByUID)
	})
}

func (r *route) getAllOrders(w http.ResponseWriter, request *http.Request) {
	orders, err := r.Service.GetAllOrders(context.TODO())

	if !checkError(w, request, err) {
		render.JSON(w, request, orders)
	}

}

func (r *route) getOrderByUID(w http.ResponseWriter, request *http.Request) {
	orderUID := chi.URLParam(request, "id")
	order, err := r.Service.GetOrderByUID(context.TODO(), orderUID)

	if !checkError(w, request, err) {
		render.JSON(w, request, order)
	}
}

func checkError(w http.ResponseWriter, request *http.Request, err error) bool {
	if errors.As(err, &pgx.ErrNoRows) {
		err = render.Render(w, request, httpErr.Err404Render(err))
		if err != nil {
			return true
		}
		return true
	} else if err != nil {
		err = render.Render(w, request, httpErr.Err500Render(err))
		if err != nil {
			return true
		}
		return true
	}
	return false
}
