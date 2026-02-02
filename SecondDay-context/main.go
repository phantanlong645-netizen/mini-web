package main

import (
	"Jee/Jee"
	"net/http"
)

func main() {
	r := Jee.NewEngine()
	r.GET("/", func(c *Jee.Context) {
		c.HTML(200, "<h1>Hello Jee</h1>")
	})
	r.GET("/hello", func(c *Jee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *Jee.Context) {
		c.JSON(http.StatusOK, Jee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":8021")
}
