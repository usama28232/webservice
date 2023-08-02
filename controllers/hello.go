package controllers

import (
	"net/http"
	"webservice/helpers"
)

type HelloController struct {
}

func (controller HelloController) ServeHTTP(writer *CustomRespWriter, request *http.Request) {
	helpers.EncodeResponse("Hello World!", writer.ResponseWriter)
}

func newHelloController() *HelloController {
	return &HelloController{}
}
