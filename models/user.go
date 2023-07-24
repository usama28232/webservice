// Models package contains Entities
package models

import (
	"errors"
	"fmt"
	"time"
)

// User struct for holding data
type User struct {
	Id        int
	FirstName string
	LastName  string
}

// Package level vars
var (
	users  []*User
	nextId = 1
)

// GetAllUsers represents all User collection
//
// returns all Users
func GetAllUsers() []*User {
	return users
}

// GetCurrentId provides maximum id from the collection
//
// returns max id
func GetCurrentId() int {
	return nextId
}

// AddNewUser adds a provided user to the collection
//
// returns newly added User object or error
func AddNewUser(u User) (User, error) {
	if u.Id != 0 {
		return User{}, errors.New("new user must not have `Id` field")
	}
	u.Id = nextId
	nextId++
	users = append(users, &u)
	return u, nil
}

// Updates existing users in collection
//
// returns updated User object or not found error
func UpdateUser(user User) (User, error) {
	for i, u := range users {
		if user.Id == u.Id {
			users[i] = &user
		}
	}
	return User{}, fmt.Errorf("user id `%v` not found", user.Id)
}

// Finds user by provided id
//
// returns user by id from collection or not found error
func GetUserById(id int) (User, error) {
	for _, user := range users {
		if user.Id == id {
			return *user, nil
		}
	}
	return User{}, fmt.Errorf("user id `%v` not found", id)
}

// Removes user by provided id
//
// returns not found error if any
func RemoveUserById(id int) error {
	for i, u := range users {
		if u.Id == id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user id `%v` not found", id)
}

// Fake delay
func ExecutionTimeSeconds() {
	time.Sleep(3 * time.Second)
}
