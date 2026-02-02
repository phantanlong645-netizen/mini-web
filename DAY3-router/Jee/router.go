package Jee

import (
	"net/http"
	"strings"
)

type router struct {
	route    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		route:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func ParsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	var parts []string
	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (router *router) Addroute(Method string, url string, handler HandlerFunc) {
	Parts := ParsePattern(url)
	key := Method + "-" + url
	if _, ok := router.route[Method]; !ok {
		router.route[Method] = &node{}
	}
	router.route[Method].insert(url, Parts, 0)
	router.handlers[key] = handler
}

func (router *router) getRoute(method string, url string) (*node, map[string]string) {
	searchParts := ParsePattern(url)
	params := make(map[string]string)
	root, ok := router.route[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := ParsePattern(n.pattern)
		for index, v := range parts {
			if v[0] == ':' {
				params[v[1:]] = searchParts[index]
			}
			if v[0] == '*' {
				params[v[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

	}

	return n, params
}
func (router *router) getRoutes(method string) []*node {
	root, ok := router.route[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
