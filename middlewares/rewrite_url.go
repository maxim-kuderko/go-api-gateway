package middlewares

import (
	"encoding/json"
	"github.com/savsgio/gotils/strconv"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strings"
)

type RewriteRule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func rewriteUrl(cfg json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	var rule *RewriteRule
	if err := json.Unmarshal(cfg, &rule); err != nil {
		logrus.Fatal(err)
	}
	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			u := strings.Replace(strconv.B2S(ctx.Request.URI().PathOriginal()), rule.From, rule.To, 1) // todo: remove alloc
			ctx.Request.URI().SetPath(u)
			handler(ctx)
		}
	}
}
