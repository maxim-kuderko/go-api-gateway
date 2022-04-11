package proxy

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go-api-gateway/entities"
	"go-api-gateway/middlewares"
	"net"
	"time"
)

func New(v *viper.Viper, route *entities.Route, cfg *entities.Config) fasthttp.RequestHandler {
	c := getClient(route)
	handler := getHandler(v, route, c)
	withMiddlewares := buildMiddlewares(route, cfg, handler)
	return withMiddlewares
}

func getClient(route *entities.Route) *fasthttp.HostClient {
	return &fasthttp.HostClient{
		Addr:                          route.Origin,
		ReadTimeout:                   time.Second, //todo: move to config
		WriteTimeout:                  time.Second, //todo: move to config
		DisableHeaderNamesNormalizing: true,        //todo: move to config
		DisablePathNormalizing:        true,        //todo: move to config
		MaxConns:                      1024,        //todo: move to config
	}
}

func getHandler(v *viper.Viper, route *entities.Route, client *fasthttp.HostClient) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		req := &ctx.Request
		resp := &ctx.Response
		prepareRequest(ctx, route)
		if err := client.Do(req, resp); err != nil {
			logrus.Debug("error when proxying the request: %s", err)
			resp.SetStatusCode(504)
			resp.SetBodyString(err.Error()) //todo: move to config
		}
		postprocessResponse(resp)
	}
}

func buildMiddlewares(route *entities.Route, cfg *entities.Config, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	middleWaresArr := make([]middlewares.Middleware, 0, len(route.Middlewares)+len(cfg.DefaultMiddlewares))
	unique := map[string]bool{}
	allMiddlewareConfigs := append(route.Middlewares, cfg.DefaultMiddlewares...)
	for _, m := range allMiddlewareConfigs {
		if _, ok := unique[m.Name]; ok {
			continue
		}
		unique[m.Name] = true
		middleWaresArr = append(middleWaresArr, middlewares.Factory(m.Name, m.Settings))
	}

	for i := len(middleWaresArr) - 1; i >= 0; i-- {
		handler = middleWaresArr[i](handler)
	}
	return handler
}

func prepareRequest(ctx *fasthttp.RequestCtx, route *entities.Route) {
	for _, h := range headersToblock {
		ctx.Request.Header.Del(h)
	}
	ctx.Request.SetHost(route.Origin)
	if ip, _, err := net.SplitHostPort(ctx.RemoteAddr().String()); err == nil { //todo: remove alloc
		ctx.Request.Header.Add("X-Forwarded-For", ip)
	}
}

func postprocessResponse(resp *fasthttp.Response) {
	for _, h := range headersToblock {
		resp.Header.Del(h)
	}
}

var headersToblock = []string{
	"Connection",          // Connection
	"Proxy-Connection",    // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",          // Keep-Alive
	"Proxy-Authenticate",  // Proxy-Authenticate
	"Proxy-Authorization", // Proxy-Authorization
	"Te",                  // canonicalized version of "TE"
	"Trailer",             // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",   // Transfer-Encoding
	"Upgrade",             // Upgrade
}
