// TITLE: Login functions
// Created by (Jacky: FafnirZ) (09/21)
// Last modified by (Jacky: FafnirZ)(04/09/21)
// # # #
/*
This module handles the Login logic, with some input validation, as well
as querying from `database` module and comparing the user's hash if the
user exists
**/
package auth

import (
	"DiffSync/database"
	_session "DiffSync/internal/session"
	"log"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

var httpDbPool database.Pool

func init() {
	var err error
	httpDbPool, err = database.NewPool(database.Config{
		HostAndPort: "db:5432",
		User:        "postgres",
		Password:    "postgres",
		Database:    "test_db",
	})

	if err != nil {
		log.Print(err.Error())
	}
}

type User struct {
	Email    string
	Password string
}

// data is in string because it could be
// stringified json object
// e.g. {"data":"email format invalid"}
type response struct {
	Data string `json:"data"`
}

// generic function for returning error in a suitable
// api format
func throwErr(err error, w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	r := response{Data: err.Error()}
	// jsonify data
	response, _ := json.Marshal(r)
	fmt.Fprint(w, string(response))
}

// EXPECT it to be from a form (handle non form requests)
// email=___&password=____
// content-type: x-www-urlencoded
/*
	curl -v -d email=john.smith@gmail.com \
	-d password=password \
	localhost:8080/login
*/
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get fields from form
		email := r.FormValue("email")
		password := r.FormValue("password")

		// initialise user class
		var user *User = &User{Email: email, Password: password}

		// input validation
		err = user.isValidEmail()
		if err != nil {
			throwErr(err, w)
			return
		}

		err = user.checkPassword()
		if err != nil {
			throwErr(err, w)
			return
		}

		// else create a session if user's session isnt already created
		_session.CreateSession(w, r, user.Email)

		// will change to FRONTEND_URI soon
		http.Redirect(w, r, "http://localhost:3000/dashboard", http.StatusMovedPermanently)
		break
	case "DEFAULT":
		// only post requests are allowed
		http.Redirect(w, r, "http://localhost:3000/login", http.StatusMovedPermanently)
		break
	}

}

// check username is valid
// uses a relatively strict regular expression
func (u *User) isValidEmail() error {
	// email must be > 2 characters before the domain
	// white lists a few domains which are allowed
	// match alphanumeric greater than 2 letters
	// followed by optional (.something) greater than 0 times
	// I know it will fail the z{8} z{6} cases
	regex := `^(z[0-9]{7}|([a-zA-Z0-9]{2,})+(\.[a-zA-Z0-9]+)*)@(gmail.com|ad.unsw.edu|student.unsw.edu|hotmail.com|outlook.com)(.au)?$`
	r, _ := regexp.Compile(regex)

	// if the email doesnt match, throw error
	if match := r.MatchString(u.Email); !match {
		return errors.New("email format invalid")
	}
	return nil
}

// inports getCredentials from internal database package
func (u *User) checkPassword() error {

	// hash the user's password first
	hashedPassword := u.hashPassword()

	// do not need input validation for password, since the password
	// gets hashed before being placed into sql query
	// the user does not have direct control over the sql query
	// i.e. even if the user puts ' or 1= '1
	// the hashed value will not be an injection command

	matches := CredentialsMatch(u.Email, hashedPassword)
	if matches == 0 {
		return errors.New("invalid credentials")
	} else if matches > 1 { // if more than 1 result is returned
		log.Print("!!alert!! duplicate results found")
		// will handle this differently later
		return errors.New("internal server error")
	}

	return nil

}

// todo hash function
func (u *User) hashPassword() string {
	// todo currently does not hash the password
	return u.Password
}