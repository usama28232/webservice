package controllers

import (
	"net/http"
	"regexp"
	"strconv"
	"webservice/helpers"
	"webservice/loggers"
	"webservice/models"
)

// var logger *zap.SugaredLogger

// UserController Struct
type UserController struct {
	pattern *regexp.Regexp
}

// Landing point for Http requests to /users
//
// Returns void
func (controller UserController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	// get Logger instance by username
	appUser, _ := helpers.ExtractAppUser(request)
	logger := loggers.GetLoggerbyUsername(appUser.Username)

	logger.Info("Inside User Controller Entrypoint")
	logger.Debugw("User Controller Entrypoint", "Url", request.URL.Path, "Method", request.Method)
	if request.URL.Path == "/users" {
		switch request.Method {
		case http.MethodGet:
			logger.Info("Get All Users")
			controller.getAll(writer, request)
		case http.MethodPost:
			logger.Info("Add new User")
			controller.post(writer, request, appUser)
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
			controller.getById(id, writer, appUser)
		case http.MethodPut:
			logger.Info("Update User")
			controller.put(id, writer, request, appUser)
		case http.MethodDelete:
			logger.Info("Delete User")
			controller.remove(id, writer, appUser)
		default:
			logger.Debugw("Invalid Request", "Method", request.Method)
			writer.WriteHeader(http.StatusNotFound)
		}
	}
}

func (controller UserController) getAll(w http.ResponseWriter, r *http.Request) {
	helpers.EncodeResponse(models.GetAllUsers(), w)
}

func (controller UserController) getById(id int, w http.ResponseWriter, appUser models.AppUser) {
	logger := loggers.GetLoggerbyUsername(appUser.Username)
	data, err := models.GetUserById(id)
	if err == nil {
		logger.Debugw("Current Users", "v", data)
		helpers.EncodeResponse(data, w)
	} else {
		logger.Debugw("Error getting all users", "e", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

}

func (controller UserController) remove(id int, w http.ResponseWriter, appUser models.AppUser) {
	logger := loggers.GetLoggerbyUsername(appUser.Username)
	err := models.RemoveUserById(id)
	if err != nil {
		logger.Debugw("Error Removing User", "e", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		logger.Debugw("Removed user", "id", id)
	}
}

func (controller UserController) put(id int, w http.ResponseWriter, r *http.Request, appUser models.AppUser) {
	logger := loggers.GetLoggerbyUsername(appUser.Username)
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
			if nerr != nil {
				logger.Debugw("Error updating user", "e", nerr.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(nerr.Error()))
			} else {
				logger.Debugw("Updated User", "v", nuser)
				helpers.EncodeResponse(nuser, w)
			}
		}
	}

}

func (controller UserController) post(w http.ResponseWriter, r *http.Request, appUser models.AppUser) {
	logger := loggers.GetLoggerbyUsername(appUser.Username)
	user, err := helpers.ParseRequest(r)
	if err != nil {
		logger.Debugw("Parsing Error", "e", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		nuser, nerr := models.AddNewUser(user)

		if nerr != nil {
			logger.Debugw("Error adding user", "e", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			logger.Debugw("Added new user", "v", nuser)
			helpers.EncodeResponse(nuser, w)
		}
	}

}

func newUserController() *UserController {
	return &UserController{
		pattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}
