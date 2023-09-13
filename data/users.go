package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	UserName  string `json:"userName"`
	Age       uint8  `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	SKU       string `json:"-"`
	CreatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

type Users []*User

func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

func (u *Users) ToJSON(rw io.Writer) error {
	e := json.NewEncoder(rw)
	return e.Encode(u)
}

func GetUsers() Users {
	return usersList
}

func AddUser(u *User) {
	u.ID = usersList[len(usersList)-1].ID + 1
	usersList = append(usersList, u)
}

func UpdateUser(id int, u *User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	u.ID = id
	usersList[pos] = u
	return nil
}

var UserNotFoundError = fmt.Errorf("Cannot find such user")

func findUser(id int) (*User, int, error) {
	for i, u := range usersList {
		if u.ID == id {
			return u, i, nil
		}
	}
	return nil, -1, UserNotFoundError
}

var usersList = []*User{
	&User{
		ID:        1,
		UserName:  "Kek1337",
		Age:       15,
		Email:     "kek@kek.com",
		Password:  "123123123",
		SKU:       "a23czw",
		CreatedOn: time.Now().UTC().String(),
		DeletedOn: time.Now().UTC().String(),
	},
}
