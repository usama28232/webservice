package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"webservice/helpers"
	"webservice/models"
)

var server *httptest.Server

func spinTestServer() {
	fmt.Println("____ SERVER STARTED ____")
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newUserController().ServeHTTP(&CustomRespWriter{ResponseWriter: w}, r)
	}))

}

func TestMain(m *testing.M) {
	fmt.Println("**** Test Main Before ****")
	spinTestServer()
	intCode := m.Run()
	fmt.Println("**** Test Main After ****")
	fmt.Println("Closing Server OBJ")
	defer server.Close()
	os.Exit(intCode)
}

func TestHttptestServer(t *testing.T) {
	// server.Start()
	defer server.Close()

	// init logger explicitly
	helpers.GetDefaultLogger()

	fmt.Println("TestHttptestServer Exec")
	client := server.Client()
	body := `{
		"FirstName": "John",
		"LastName": "Wick"
	}`
	resp, err := client.Post(server.URL+"/users", "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
	} else {
		fmt.Printf("Http Response Code: %v\n", resp.StatusCode)
	}
	_body, errResp := io.ReadAll(resp.Body)
	if errResp == nil {
		fmt.Printf("Body: %v", string(_body))
		user := models.User{}
		parseErr := json.Unmarshal([]byte(_body), &user)
		if parseErr != nil {
			t.Errorf("Parse Err: %v", err.Error())
		} else {
			fmt.Printf("User Object %v\n", user)

			// ....
			// client.Get(server.URL + "/users/" + string(user.Id))

		}

	} else {
		t.Errorf("Response Err: %v", errResp.Error())
	}

}
