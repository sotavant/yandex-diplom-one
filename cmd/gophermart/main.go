package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/repository/pgsql"
	"github.com/sotavant/yandex-diplom-one/internal/rest"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"github.com/sotavant/yandex-diplom-one/order"
	"github.com/sotavant/yandex-diplom-one/user"
	"github.com/sotavant/yandex-diplom-one/withdrawn"
	"net/http"
)

// test encoding
func main() {
	ctx := context.Background()
	app, err := internal.InitApp(ctx)

	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Gzip)

	userRepo, err := pgsql.NewUserRepository(ctx, app.DBPool)
	if err != nil {
		panic(err)
	}
	ordersRepo, err := pgsql.NewOrderRepository(ctx, app.DBPool)
	if err != nil {
		panic(err)
	}
	wdRepo, err := pgsql.NewWithdrawnRepository(ctx, app.DBPool)
	if err != nil {
		panic(err)
	}

	userService := user.NewService(userRepo)
	orderService := order.NewOrderService(ordersRepo)
	wdService := withdrawn.NewService(wdRepo, userRepo)

	rest.NewUserHandler(r, userService)
	rest.NewOrdersHandler(r, orderService)
	rest.NewWithdrawnHandler(r, wdService)

	err = http.ListenAndServe(app.Address, r)
	if err != nil {
		panic(err)
	}
}
