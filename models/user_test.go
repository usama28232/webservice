package models

import "testing"

func TestAddNewUser(t *testing.T) {
	user := User{
		FirstName: "su",
		LastName:  "root",
	}
	got, err := AddNewUser(user)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	} else {
		t.Logf("success: %v", got)
	}
}
