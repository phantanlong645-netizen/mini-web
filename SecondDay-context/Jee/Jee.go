package Jee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)
type engine struct {
	router *router
}

func NewEngine() *engine {
	engine := &engine{
		router: newRouter(),
	}
	return engine
}

func (engine *engine) AddRoute(Method string, url string, handler HandlerFunc) {
	key := Method + "-" + url
	log.Printf("新加了路由 %s    %s", key, url)
	engine.router.Addroute(Method, url, handler)
}
func (engine *engine) GET(url string, handler HandlerFunc) {
	engine.AddRoute("GET", url, handler)
}
func (engine *engine) POST(url string, handler HandlerFunc) {
	engine.AddRoute("POST", url, handler)
}

func (engine *engine) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, engine))
}
func (engine *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)

}
