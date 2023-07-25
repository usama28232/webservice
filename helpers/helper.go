// Helpers package (shared package) contain all helper functions
package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"webservice/constants"
	"webservice/models"
)

// ParseRequest decodes incoming request and extracts User Model
//
// Returns User model or error
func ParseRequest(r *http.Request) (models.User, error) {
	dec := json.NewDecoder(r.Body)
	var u models.User
	err := dec.Decode(&u)
	if err == nil {
		return u, nil
	} else {
		return u, err
	}
}

// ParseDebugRequest decodes incoming request and extracts Debug Param
//
// Returns Config Struct
func ParseDebugRequest(r *http.Request) (models.Param, error) {
	body, err := io.ReadAll(r.Body)
	p := models.Param{}
	p.Debug = false
	if err == nil && len(body) > 0 {
		json.Unmarshal(body, &p)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		return p, nil
	}
	return p, nil
}

// Encode incoming data model
//
// Returns error if encoding fails
func EncodeResponse(data interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	err := enc.Encode(data)
	if err == nil {
		return nil
	} else {
		return err
	}
}

// Extracts Username from Request Header
//
// Returns AppUser
func ExtractAppUser(r *http.Request) (models.AppUser, error) {
	username := r.Header.Get(constants.USER_HEADER_KEY)
	u := models.AppUser{}
	if username == "" {
		return u, errors.New("could not extract '" + constants.USER_HEADER_KEY + "' from request header")
	} else {
		u.Username = username
	}
	return u, nil
}
