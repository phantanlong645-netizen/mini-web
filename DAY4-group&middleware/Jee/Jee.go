package Jee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type (
	RouteGroup struct {
		prefix      string
		parent      *RouteGroup
		engine      *engine
		middlewares []HandlerFunc
	}
	engine struct {
		*RouteGroup
		router *router
		groups []*RouteGroup
	}
)

func NewEngine() *engine {
	engine := &engine{router: newRouter()}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

func (routegroup *RouteGroup) Group(prefix string) *RouteGroup {
	engine := routegroup.engine
	newRouteGroup := &RouteGroup{engine: engine, prefix: routegroup.prefix + prefix, parent: routegroup}
	engine.groups = append(engine.groups, newRouteGroup)
	return newRouteGroup
}

func (routeGroup *RouteGroup) Use(handlers ...HandlerFunc) {
	routeGroup.middlewares = append(routeGroup.middlewares, handlers...)
}

func (routegroup *RouteGroup) addRoute(method string, comb string, handler HandlerFunc) {
	pattern := routegroup.prefix + "-" + comb
	log.Printf("Route %4s - %s", method, pattern)
	routegroup.engine.router.addRoute(method, pattern, handler)
}

func (routegroup *RouteGroup) GET(pattern string, handler HandlerFunc) {
	routegroup.addRoute("GET", pattern, handler)
}
func (routegroup *RouteGroup) POST(pattern string, handler HandlerFunc) {
	routegroup.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)

}
