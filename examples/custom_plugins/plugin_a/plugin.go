package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func Name() string {
	return `test-plugin`
}

func Middleware(settings json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			ctx.Request.Header.Add(Name(), `Custom-Value`)
		}
	}
}
