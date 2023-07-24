package helpers

import (
	"bytes"
	"net/http"
	"testing"
	"webservice/models"
)

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

func TestParseRequest(t *testing.T) {
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
		t.Log(fakePostRequest)
		user, err := ParseRequest(fakePostRequest)
		if err == nil {
			t.Log(user)
		} else {
			t.Fatal(err)
		}
	}
}

func TestEncodeResponse(t *testing.T) {
	user := models.User{
		Id:        1,
		FirstName: "John",
		LastName:  "Smith",
	}
	rw := mockHttpRW{}
	err := EncodeResponse(user, &rw)

	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(rw.getData()))
	}

}
