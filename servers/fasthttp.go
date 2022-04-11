package servers

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"log"
)

func NewFastHTTP(v *viper.Viper, router *router.Router) {
	srv := fasthttp.Server{
		Handler:               router.Handler,
		TCPKeepalive:          true,
		NoDefaultServerHeader: true,
		NoDefaultDate:         true,
		NoDefaultContentType:  true,
	}
	if err := srv.ListenAndServe(fmt.Sprintf(":%s", v.GetString(`HTTP_LISTEN_PORT`))); err != nil {
		log.Fatal(err)
	}
}
