package controllers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
	"webservice/constants"
	"webservice/loggers"
	"webservice/models"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ControllerBase interface {
	ServeHTTP(*CustomRespWriter, *http.Request)
}

// add other fields to capture here
type CustomRespWriter struct {
	http.ResponseWriter
	txid   string
	status int
}

func (c *CustomRespWriter) WriteHeader(code int) {
	c.status = code
	c.ResponseWriter.WriteHeader(code)
}

func (c *CustomRespWriter) WriteTrxid(value int) {
	c.txid = strconv.Itoa(value)
}

func (c *CustomRespWriter) WriteTrxidString(value string) {
	c.txid = value
}

func (c *CustomRespWriter) GetTrxid() string {
	return c.txid
}

// RegisterControllers contain other available endpoints
func RegisterControllers() *mux.Router {
	mux := mux.NewRouter()
	// index page
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	userController := newUserController()
	helloController := newHelloController()

	mux.HandleFunc("/users", middleware(userController)).Methods(http.MethodGet, http.MethodPost)
	mux.HandleFunc("/helloworld", middleware(helloController)).Methods(http.MethodGet, http.MethodPost)
	mux.HandleFunc("/users/{id:[0-9]+}", middleware(userController)).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	mux.StrictSlash(false)
	return mux
}

func middleware(controller ControllerBase) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		meta := models.HttpRequest{}

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

			meta.Method = r.Method
			meta.Url = r.URL.Path
			meta.Data = strBody
			meta.Agent = r.UserAgent()

			// Reset the request body so it can be read again by the actual handler
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		} else {
			meta.Method = r.Method
			meta.Url = r.URL.Path
			meta.Agent = r.UserAgent()
		}

		cw := &CustomRespWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		controller.ServeHTTP(cw, r)

		cookie, err := r.Cookie(constants.SESSION_ID)
		if err == nil {
			meta.SessionId = cookie.Value
		}

		meta.Duration = time.Since(startTime).Milliseconds()
		meta.Status = cw.status
		if len(cw.GetTrxid()) > 0 {
			meta.Trxid = cw.GetTrxid()
		}

		accessLogger := loggers.GetAccessLogger()
		accessLogger.Infow("http", zap.Any("v", meta))
	})

}
