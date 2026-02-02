package Jee

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (router *router) Addroute(Method string, url string, handler HandlerFunc) {
	key := Method + "-" + url
	router.handlers[key] = handler
}

func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	handler := router.handlers[key]
	if handler != nil {
		handler(c)
	} else {
		c.String(404, "404 NOT FOUND")
	}
}
