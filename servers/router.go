package servers

import (
	"github.com/fasthttp/router"
	"github.com/spf13/viper"
	"go-api-gateway/entities"
	"go-api-gateway/proxy"
)

func NewRouter(v *viper.Viper, cfg *entities.Config) *router.Router {
	r := router.New()
	for _, route := range cfg.Routes {
		for _, method := range route.Methods {
			r.Handle(method, route.IngressPath, proxy.New(v, route, cfg))
		}
	}
	return r
}
