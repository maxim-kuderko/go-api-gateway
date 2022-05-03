package middlewares

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Middleware func(handler fasthttp.RequestHandler) fasthttp.RequestHandler

func Factory(t string, settings json.RawMessage) Middleware {
	var fn func(settings json.RawMessage) Middleware
	switch t {
	case `request_id`:
		fn = requestID
	case `rewrite_url`:
		fn = rewriteUrl
	case `compress`:
		fn = compress
	case `jwtAuth`:
		fn = jwtAuth
	case `prometheus`:
		fn = prometheusMonitor
	default:
		logrus.Fatalf("unknown middleware %s, exiting", t)
	}
	return fn(settings)
}
