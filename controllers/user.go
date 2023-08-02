package controllers

import (
	"net/http"
	"regexp"
	"strconv"
	"webservice/helpers"
	"webservice/models"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// UserController Struct
type UserController struct {
	pattern *regexp.Regexp
}

// Landing point for Http requests to /users
//
// Returns void
func (controller UserController) ServeHTTP(writer *CustomRespWriter, request *http.Request) {
	logger = helpers.GetLoggerByRequest(request)
	logger.Info("Inside User Controller Entrypoint")
	logger.Debugw("User Controller Entrypoint", "URL.PATH", request.URL.Path, "Method", request.Method)
	if request.URL.Path == "/users" {
		switch request.Method {
		case http.MethodGet:
			logger.Info("Get All Users")
			controller.getAll(writer, request)
		case http.MethodPost:
			logger.Info("Add new User")
			controller.post(writer, request)
		}
	} else {
		matches := controller.pattern.FindStringSubmatch(request.URL.Path)
		if len(matches) == 0 {
			logger.Debugw("Cannot find method on", request.URL.Path)
			writer.WriteHeader(http.StatusNotFound)
		}
		id, err := strconv.Atoi(matches[1])
		if err != nil {
			logger.Debugw("Cannot find method on", id)
			writer.WriteHeader(http.StatusNotFound)
		}
		switch request.Method {
		case http.MethodGet:
			logger.Info("Get User by Id")
			controller.getById(id, writer, request)
		case http.MethodPut:
			logger.Info("Update User")
			controller.put(id, writer, request)
		case http.MethodDelete:
			logger.Info("Delete User")
			controller.remove(id, writer)
		default:
			logger.Debugw("Invalid Request", "Method", request.Method)
			writer.WriteHeader(http.StatusNotFound)
		}
	}
}

func (controller UserController) getAll(w *CustomRespWriter, r *http.Request) {
	helpers.EncodeResponse(models.GetAllUsers(), w.ResponseWriter)
}

func (controller UserController) getById(id int, w *CustomRespWriter, r *http.Request) {
	data, err := models.GetUserById(id)
	if err == nil {
		logger.Debugw("Current Users", "v", data)
		w.WriteTrxid(data.Id)
		helpers.EncodeResponse(data, w.ResponseWriter)
	} else {
		logger.Debugw("Error getting all users", "e", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (controller UserController) remove(id int, w *CustomRespWriter) {
	err := models.RemoveUserById(id)
	if err != nil {
		logger.Debugw("Error Removing User", "e", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		logger.Debugw("Removed user", "v", id)
		w.WriteTrxid(id)
	}
}

func (controller UserController) put(id int, w *CustomRespWriter, r *http.Request) {
	user, err := helpers.ParseRequest(r)
	if err != nil {
		logger.Debugw("Parsing Error", "e", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		if id != user.Id {
			logger.Debugw("Security Error", "Data mismatch in request URL & Request Body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Id mismatch"))
		} else {
			nuser, nerr := models.UpdateUser(user)
			w.WriteTrxid(id)
			if nerr != nil {
				logger.Debugw("Error updating user", "e", nerr.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(nerr.Error()))
			} else {
				logger.Debugw("Updated User", "v", nuser)
				helpers.EncodeResponse(nuser, w.ResponseWriter)
			}
		}
	}

}

func (controller UserController) post(w *CustomRespWriter, r *http.Request) {
	user, err := helpers.ParseRequest(r)
	if err != nil {
		logger.Debugw("Parsing Error", "e", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		nuser, nerr := models.AddNewUser(user)
		w.WriteTrxid(nuser.Id)

		if nerr != nil {
			logger.Debugw("Error adding user", "e", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			logger.Debugw("Added new user", "v", nuser)
			helpers.EncodeResponse(nuser, w.ResponseWriter)
		}
	}

}

func newUserController() *UserController {
	return &UserController{
		pattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}
