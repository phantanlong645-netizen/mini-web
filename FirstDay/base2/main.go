package main

import (
	"fmt"
	"log"
	"net/http"
)

type engine struct {
}

// 只要一个结构体实现了 ServeHTTP 这个方法，它就变成了一个“处理器”
func (engine *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "the url is %q", req.URL.Path)
	case "/hello":
		fmt.Printf("the url is %q\n", req.URL.Path)
		for k, v := range req.Header {
			fmt.Printf("the key is %q and the v is %q \n", k, v)
		}
	}
}

func main() {
	engine := new(engine)
	log.Fatal(http.ListenAndServe(":8012", engine))
}
