package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sotavant/yandex-diplom-one/internal"
	"github.com/sotavant/yandex-diplom-one/internal/rest"
	"github.com/sotavant/yandex-diplom-one/internal/rest/middleware"
	"github.com/sotavant/yandex-diplom-one/user"
)

// test encoding
func main() {
	ctx := context.Background()
	_, err := internal.InitApp(ctx)

	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Gzip)

	userService := user.NewService()
	rest.NewUserHandler(r, userService)
}
