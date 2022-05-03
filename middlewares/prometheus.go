package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	strconv2 "github.com/savsgio/gotils/strconv"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type prometheusConfig struct {
	Path             string    `json:"path"`
	Port             string    `json:"Port"`
	HistogramBuckets []float64 `json:"histogram_buckets"`
}

func prometheusMonitor(cfg json.RawMessage) Middleware {
	var config *prometheusConfig
	if err := json.Unmarshal(cfg, &config); err != nil {
		logrus.Fatal(err)
	}
	http.Handle(config.Path, promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
	reqCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total",
		Help: "The HTTP request counts processed.",
	}, []string{"status_code", "method", "path"})

	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			handler(ctx)
			reqCounter.WithLabelValues(strconv.Itoa(ctx.Response.StatusCode()), strconv2.B2S(ctx.Request.Header.Method()), ctx.UserValue(router.MatchedRoutePathParam).(string)).Inc()
		}
	}
}
