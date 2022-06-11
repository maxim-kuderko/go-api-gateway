package middlewares

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/savsgio/gotils/strconv"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
)

type jwtConfig struct {
	Skip              bool   `json:"skip"`
	SecretEnvVarName  string `json:"secret_env_var_name"`
	TokenFormat       string `json:"token_format"`
	TokenHeader       string `json:"token_header"`
	ParsedTokenHeader string `json:"parsed_token_header"`
}

func jwtAuth(cfg json.RawMessage) func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	var config *jwtConfig
	if err := json.Unmarshal(cfg, &config); err != nil {
		logrus.Fatal(err)
	}
	secret := []byte(os.Getenv(config.SecretEnvVarName))
	return func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if config.Skip {
				handler(ctx)
				return
			}
			tokenString, err := getToken(ctx, config)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusForbidden)
				ctx.WriteString(err.Error())
				return
			}
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil || !token.Valid {
				ctx.SetStatusCode(fasthttp.StatusForbidden)
				ctx.WriteString(err.Error())
				return
			}
			segment, _ := jwt.DecodeSegment(strings.Split(tokenString, ".")[1])
			ctx.Request.Header.SetBytesV(config.ParsedTokenHeader, segment)
			handler(ctx)
		}
	}
}

var (
	badAuthorizationHeader     = errors.New(`bad authorization header`)
	mustSpecifyJwtHeaderFormat = errors.New(`must specify header format`)
)

func getToken(ctx *fasthttp.RequestCtx, config *jwtConfig) (string, error) {
	switch config.TokenFormat {
	case `Bearer`:
		header := ctx.Request.Header.Peek(`Authorization`)
		if len(header) < 7 { // length of "Bearer
			return ``, badAuthorizationHeader
		}
		return strconv.B2S(header[7:]), nil
	case "Custom":
		return strconv.B2S(ctx.Request.Header.Peek(config.TokenHeader)), nil
	default:
		return ``, mustSpecifyJwtHeaderFormat
	}
}
