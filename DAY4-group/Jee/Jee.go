package Jee

import (
	"log"
	"net/http"
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
	c := newContext(w, req)
	engine.router.handle(c)
}
