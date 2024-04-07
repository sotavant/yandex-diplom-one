package main

import (
	"context"
	"github.com/sotavant/yandex-diplom-one/internal"
)

func main() {
	ctx := context.Background()
	app, err := internal.InitApp(ctx)
	if err != nil {
		panic(err)
	}
}
