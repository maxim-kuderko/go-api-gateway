package entities

type Config struct {
	DefaultMiddlewares []*Middleware `json:"default_middlewares"`
	Routes             Routes        `json:"routes"`
}
