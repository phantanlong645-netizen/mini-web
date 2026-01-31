package main

import (
	"base3/Jee"
	"fmt"
	"net/http"
)

func main() {
	engine := Jee.NewEngine()

	engine.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Jee")

	})
	engine.POST("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Jee POST")
	})
	engine.Run(":8014")

}
