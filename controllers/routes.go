package controllers

import (
	"bytes"
	"io"
	"net/http"
	"webservice/helpers"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// RegisterControllers contain other available endpoints
func RegisterControllers() http.Handler {
	userCont := newUserController()

	mux := mux.NewRouter()

	// index page
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	mux.HandleFunc("/users", userCont.ServeHTTP).Methods(http.MethodGet, http.MethodPost)
	mux.HandleFunc("/users/{id:[0-9]+}", userCont.ServeHTTP).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	mux.StrictSlash(false)

	accessLogger := helpers.GetAccessLogger()
	return loggingMiddleware(accessLogger, mux)
}

// loggingMiddleware is a middleware that logs incoming HTTP requests
func loggingMiddleware(log *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request using the Zap logger
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			body, err := io.ReadAll(r.Body)
			var strBody string
			if err != nil {
				strBody = "Err in request body"
				// http.Error(w, "Error reading request body", http.StatusInternalServerError)
				// return
			} else {
				// Log the request body
				strBody = string(body)
			}

			log.Infow("HttpRequest",
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.String("data", strBody),
				zap.String("user_agent", r.UserAgent()),
			)

			// Reset the request body so it can be read again by the actual handler
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		} else {
			log.Infow("HttpRequest",
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.String("user_agent", r.UserAgent()),
			)
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
