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
	"time"
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
	reqDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "The HTTP request duration in seconds.",
			Buckets: config.HistogramBuckets,
		},
		[]string{"status_code", "method", "path"},
	)

	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			t := time.Now()
			handler(ctx)
			statusCode := strconv.Itoa(ctx.Response.StatusCode())
			method := strconv2.B2S(ctx.Request.Header.Method())
			route := ctx.UserValue(router.MatchedRoutePathParam).(string)
			reqCounter.WithLabelValues(statusCode, method, route).Inc()
			reqDuration.WithLabelValues(statusCode, method, route).Observe(time.Since(t).Seconds())
		}
	}
}