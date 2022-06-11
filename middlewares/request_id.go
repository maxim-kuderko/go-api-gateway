package middlewares

import (
	"encoding/json"
	"github.com/savsgio/gotils/uuid"
	"github.com/valyala/fasthttp"
)

func requestID(cfg json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			ctx.Request.Header.Set(`X-Request-ID`, uuid.V4())
			handler(ctx)
		}
	}
}
