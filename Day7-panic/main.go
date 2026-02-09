package main

import (
	"net/http"

	"Day7-panic/Jee"
)

func main() {
	r := Jee.Default()
	r.GET("/", func(c *Jee.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *Jee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9997")
}
