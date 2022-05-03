package main

import (
	"context"
	"go-api-gateway/initializers"
	"go-api-gateway/servers"
	"go.uber.org/fx"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	app := fx.New(
		fx.Provide(
			initializers.NewConfig,
			initializers.RouterConfig,
			servers.NewRouter,
		),
		fx.Invoke(servers.NewFastHTTP),
	)
	app.Start(context.Background())
}
