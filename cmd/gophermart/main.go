package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/repository/pgsql"
	"github.com/sotavant/yandex-diplom-one/internal/rest"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"github.com/sotavant/yandex-diplom-one/user"
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

	userService := user.NewService(userRepo)
	rest.NewUserHandler(r, userService)

	err = http.ListenAndServe(app.Address, r)
	if err != nil {
		panic(err)
	}
}
