// Package main is the entry point of application
//
// Spins up a web server and listen for requests
package main

import (
	"flag"
	"fmt"
	"net/http"
	"webservice/controllers"
)

// This is the entry point for application.
//
// You can change web port by specifying -port flag at runtime.
func main() {
	var port string
	flag.StringVar(&port, "port", "3000", "Specifies Webservice Port")
	flag.Parse()
	fmt.Println("Hello webservice", port)
	controllers.RegisterControllers()
	http.ListenAndServe(":"+port, nil)
}

func add(i, j int) int {
	return i + j
}
