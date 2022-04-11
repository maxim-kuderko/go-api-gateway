package middlewares

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Compress struct {
	Level int `json:"level"`
}

func compress(cfg json.RawMessage) Middleware {
	var config *Compress
	if err := json.Unmarshal(cfg, &config); err != nil {
		logrus.Fatal(err)
	}
	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return fasthttp.CompressHandlerLevel(handler, config.Level)
	}
}
