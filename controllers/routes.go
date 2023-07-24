package controllers

import "net/http"

// RegisterControllers contain other available endpoints
func RegisterControllers() {
	userCont := newUserController()

	// index page
	RegisterIndexController()

	// api
	http.Handle("/users", *userCont)
	http.Handle("/users/", *userCont)

}

func sub(i, j int) int {
	if i < j {
		return j - i
	} else {
		return i - j
	}
}
