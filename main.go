// Package main is the entry point of application
//
// Spins up a web server and listen for requests
package main

import (
	"flag"
	"net/http"
	"webservice/constants"
	"webservice/controllers"
	"webservice/loggers"
)

// This is the entry point for application.
//
// You can change web port by specifying -port flag at runtime.
func main() {
	logger := loggers.GetLogger(constants.Info)

	var port string
	flag.StringVar(&port, "port", "3000", "Specifies Webservice Port")
	flag.Parse()

	logger.Infow("Starting Webservice with", "port", port)

	http.ListenAndServe(":"+port, controllers.RegisterControllers())
}
