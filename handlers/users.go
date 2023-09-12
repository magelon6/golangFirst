package handlers

import (
	"context"
	"log"
	"microservicespetprod/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	l *log.Logger
}

func NewUser (l *log.Logger) *User {
	return &User{l}
}

func(u *User) GetUsers(rw http.ResponseWriter, r *http.Request) {
	lu := data.GetUsers()
	err := lu.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Cannot unmarshall Users", http.StatusInternalServerError)
	}
}

func(u *User) AddUser(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(KeyUser{}).(data.User)
	data.AddUser(&user)
}

func(u User) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := r.Context().Value(KeyUser{}).(data.User)

	uid, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Cannot convert UID to int", http.StatusBadRequest)
		return
	}

	err = data.UpdateUser(uid, &user)

	if err == data.UserNotFoundError {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
	
}

type KeyUser struct {}

func(u User) ValidateMiddlewareUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		user := data.User{}

		err := user.FromJSON(r.Body)
	
		if err != nil {
			http.Error(rw, "Cannot unmarshall user", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}