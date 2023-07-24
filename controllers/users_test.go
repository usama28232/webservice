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
	"webservice/models"
)

var server *httptest.Server

type mockHttpRW struct {
	bytes.Buffer
}

func (w mockHttpRW) Header() http.Header {
	return http.Header{}
}

func (w mockHttpRW) WriteHeader(status int) {}

func (w mockHttpRW) getData() []byte {
	return w.Bytes()
}

func spinTestServer() {
	fmt.Println("____ SERVER STARTED ____")
	server = httptest.NewServer(http.HandlerFunc(newUserController().ServeHTTP))

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

func TestServeHTTP(t *testing.T) {

	// .. adding user
	body := `{
		"FirstName": "asdzxc",
		"LastName": "dfgdfg"
	}`
	bufferBody := bytes.NewBuffer([]byte(body))

	fakePostRequest, reqPostErr := http.NewRequest(http.MethodPost, "/users", bufferBody)
	if reqPostErr != nil {
		t.Fatal(reqPostErr)
	} else {
		mockPostrw := mockHttpRW{}
		UserController.ServeHTTP(*newUserController(), &mockPostrw, fakePostRequest)
		if len(mockPostrw.getData()) > 0 {
			t.Log("Post Success")
			fmt.Printf("post output: %v", string(mockPostrw.getData()))

			// .. getting newly created user
			fakeGetRequest, reqGetErr := http.NewRequest(http.MethodGet, "/users/1", nil)
			if reqGetErr != nil {
				t.Fatal(reqGetErr)
			} else {
				mockGetrw := mockHttpRW{}
				UserController.ServeHTTP(*newUserController(), &mockGetrw, fakeGetRequest)
				if len(mockGetrw.getData()) > 0 {
					t.Log("success")
				} else {
					t.Fatal("Could not simulate get request")
				}
				fmt.Printf("get output: %v", string(mockGetrw.getData()))
			}

		} else {
			t.Fatal("Could not simulate post request")
		}
	}

}

func TestHttptestServer(t *testing.T) {
	// server.Start()
	defer server.Close()
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
