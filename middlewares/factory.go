package middlewares

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"plugin"
	"reflect"
)

var registry = map[string]func(settings json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler{}

func init() {
	registerInternal()
	registerExternal()
}

func Factory(t string, settings json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn, ok := registry[t]
	if !ok {
		logrus.Fatalf("unknown middleware %s, exiting", t)
	}
	return fn(settings)
}

func registerInternal() {
	registry[`request_id`] = requestID
	registry[`rewrite_url`] = rewriteUrl
	registry[`compress`] = compress
	registry[`jwtAuth`] = jwtAuth
	registry[`prometheus`] = prometheusMonitor
}
func registerExternal() {
	pluginsDir := os.Getenv(`PLUGINS_DIR`)
	if pluginsDir == `` {
		return
	}
	dir, err := os.ReadDir(pluginsDir)
	if err != nil {
		logrus.Fatalf("error while reading plugins directory: %s", err.Error())
	}
	for _, entry := range dir {
		pl, err := plugin.Open(pluginsDir + `/` + entry.Name())
		if err != nil {
			logrus.Fatalf("error while openeing plugin %s, error: %s", entry.Name(), err.Error())
		}
		nm, err := pl.Lookup(`Name`)
		if err != nil {
			logrus.Fatalf("error while looking up for Name function in plugin %s, name : %s", entry.Name(), err.Error())
		}

		ml, err := pl.Lookup(`Middleware`)
		if err != nil {
			logrus.Fatalf("error while looking up for Middleware function in plugin %s, name : %s", entry.Name(), err.Error())
		}
		name, ok := nm.(func() string)
		if !ok {
			logrus.Fatalf("error while calling Name function in plugin %s, error: %s", entry.Name(), reflect.TypeOf(nm))
		}

		middleware, ok := ml.(func(settings json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler)
		if !ok {
			logrus.Fatalf("error while calling Middleware function in plugin %s, error: %s", entry.Name(), reflect.TypeOf(ml))
		}
		registry[name()] = middleware
	}
}
