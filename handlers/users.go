package handlers

import (
	"log"
	"microservicespetprod/data"
	"net/http"
)

type User struct {
	l *log.Logger
}

func NewUser (l *log.Logger) *User {
	return &User{l}
}

func(u *User) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		u.getUsers(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		u.addUser(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func(u *User) getUsers(rw http.ResponseWriter, r *http.Request) {
	lu := data.GetUsers()
	err := lu.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Cannot unmarshall Users", http.StatusInternalServerError)
	}
}

func(u *User) addUser(rw http.ResponseWriter, r *http.Request) {
	user := &data.User{}
	err := user.FromJSON(r.Body)
	
	if err != nil {
		http.Error(rw, "Cannot unmarshall user", http.StatusBadRequest)
	}
	
	u.l.Printf("User: %#v", user)
	data.AddUser(user)
}