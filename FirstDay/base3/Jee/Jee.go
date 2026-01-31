package Jee

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type engine struct {
	router map[string]HandlerFunc
}

func NewEngine() *engine {
	engine := &engine{
		router: make(map[string]HandlerFunc),
	}
	return engine
}

func (engine *engine) AddRoute(Method string, url string, handler HandlerFunc) {
	key := url + "-" + Method
	log.Printf("新加了路由 %s    %s", key, url)
	engine.router[key] = handler

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
	key := req.URL.Path + "-" + req.Method
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 Not Found")
	}

}
