package main

import (
	"context"
	"go-api-gateway/initializers"
	"go-api-gateway/servers"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			initializers.NewConfig,
			initializers.RouterConfig,
			servers.NewRouter,
		),
		fx.Invoke(servers.NewFastHTTP),
	)
	app.Start(context.Background())
}
