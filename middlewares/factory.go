package middlewares

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Middleware func(handler fasthttp.RequestHandler) fasthttp.RequestHandler

func Factory(t string, settings json.RawMessage) Middleware {
	switch t {
	case `request_id`:
		return requestID(settings)
	case `rewrite_url`:
		return rewriteUrl(settings)
	}
	logrus.Fatalf("unknown middleware %s, exiting", t)
	return nil
}
