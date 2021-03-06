package middlewares

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Compress struct {
	Level int `json:"level"`
}

func compress(cfg json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	var config *Compress
	if err := json.Unmarshal(cfg, &config); err != nil {
		logrus.Fatal(err)
	}
	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return fasthttp.CompressHandlerBrotliLevel(handler, config.Level, config.Level)
	}
}
