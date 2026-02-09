package Jee

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
		router   *router
		groups   []*RouteGroup
		funMap   template.FuncMap
		template *template.Template
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

func Default() *engine {
	engine := NewEngine()
	engine.Use(Logger(), Recovery())
	return engine
}

func (routeGroup *RouteGroup) Use(handlers ...HandlerFunc) {
	routeGroup.middlewares = append(routeGroup.middlewares, handlers...)
}
func (group *RouteGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("path")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
func (group *RouteGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}
func (engine *engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funMap = funcMap
}

func (engine *engine) LoadHTMLGlob(pattern string) {
	engine.template = template.Must(template.New("").Funcs(engine.funMap).ParseGlob(pattern))
}

func (routegroup *RouteGroup) addRoute(method string, comb string, handler HandlerFunc) {
	pattern := routegroup.prefix + comb
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
	c.engine = engine
	engine.router.handle(c)

}
