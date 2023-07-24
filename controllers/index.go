// Controllers package contains request handling controllers
package controllers

import "net/http"

// RegisterIndexController maps '/' endpoint to serve directory '/static' containing index.html file
func RegisterIndexController() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
}
