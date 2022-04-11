package entities

type Route struct {
	Origin      string        `json:"origin"`
	Methods     []string      `json:"methods"`
	IngressPath string        `json:"ingress_path"`
	Middlewares []*Middleware `json:"middlewares"`
}

type Routes []*Route
