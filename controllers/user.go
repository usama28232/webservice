package controllers

import (
	"net/http"
	"regexp"
	"strconv"
	"webservice/helpers"
	"webservice/models"
)

// UserController Struct
type UserController struct {
	pattern *regexp.Regexp
}

// Landing point for Http requests to /users
//
// Returns void
func (controller UserController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/users" {
		switch request.Method {
		case http.MethodGet:
			controller.getAll(writer, request)
		case http.MethodPost:
			controller.post(writer, request)
		}
	} else {
		matches := controller.pattern.FindStringSubmatch(request.URL.Path)
		if len(matches) == 0 {
			writer.WriteHeader(http.StatusNotFound)
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
		}
		switch request.Method {
		case http.MethodGet:
			controller.getById(id, writer)
		case http.MethodPut:
			controller.put(id, writer, request)
		case http.MethodDelete:
			controller.remove(id, writer)
		default:
			writer.WriteHeader(http.StatusNotFound)
		}
	}
}

func (controller UserController) getAll(w http.ResponseWriter, r *http.Request) {
	helpers.EncodeResponse(models.GetAllUsers(), w)
}

func (controller UserController) getById(id int, w http.ResponseWriter) {
	data, err := models.GetUserById(id)
	if err == nil {
		helpers.EncodeResponse(data, w)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller UserController) remove(id int, w http.ResponseWriter) {
	err := models.RemoveUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (controller UserController) put(id int, w http.ResponseWriter, r *http.Request) {
	user, err := helpers.ParseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	if id != user.Id {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id mismatch"))
	}
	nuser, nerr := models.UpdateUser(user)
	if nerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		helpers.EncodeResponse(nuser, w)
	}
}

func (controller UserController) post(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.ParseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	nuser, nerr := models.AddNewUser(user)
	if nerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		helpers.EncodeResponse(nuser, w)
	}
}

func newUserController() *UserController {
	return &UserController{
		pattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}
