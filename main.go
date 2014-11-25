package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	PORT         = os.Getenv("PORT")
	DATABASE_URL = os.Getenv("DATABASE_URL")
)

// apiError define structure of API error
type apiError struct {
	Error   error  `json:"-"`
	Message string `json:"error"`
	Code    int    `json:"code"`
}

// apiHandler global API mux
type apiHandler func(w http.ResponseWriter, r *http.Request) *apiError

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// add header on every response
	w.Header().Add("Server", "Automata/0.1")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// if handler return an &apiError
	err := fn(w, r)
	if err != nil {
		// http log
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL, err.Error)

		// response proper http status code
		w.WriteHeader(err.Code)

		// response JSON
		resp := json.NewEncoder(w)
		err_json := resp.Encode(err)
		if err_json != nil {
			log.Println("Encode JSON for error response was failed.")

			return
		}

		return
	}

	// http log
	// TODO: print response
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

// index handle '/' request
func index(w http.ResponseWriter, r *http.Request) *apiError {

	// response "404 not found" on every undefined
	// URL pattern handler
	if r.URL.Path != "/" {
		return &apiError{
			errors.New("Not Found"),
			"Not Found",
			http.StatusNotFound,
		}
	}

	fmt.Fprintln(w, "Hello World! - automata.")
	return nil
}

// TODO: define User type
// User type
type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// users handle '/users' request
// TODO: validate request header must Accept: application/json
func users(w http.ResponseWriter, r *http.Request) *apiError {

	switch r.Method {
	case "GET":
		return usersGET(w, r)
	case "POST":
		return usersPOST(w, r)
	default:
		return &apiError{
			errors.New("Not Found"),
			"Not Found",
			http.StatusNotFound,
		}
	}
	return nil
}

// usersGET handle 'GET' request on '/users'
// TODO: query data from database
func usersGET(w http.ResponseWriter, r *http.Request) *apiError {
	fmt.Fprintln(w, "usersGET executed")

	return nil
}

// usersPOST handle 'POST' request on '/users'
// TODO: insert data into database
// TODO: accept data json with format like this
// {"name": "bayu", "email": "bayu@gmail.com"}
// TODO: decode received JSON -> insert data into database
func usersPOST(w http.ResponseWriter, r *http.Request) *apiError {

	// decode received JSON
	var u User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		return &apiError{
			err,
			"Internal server error",
			http.StatusInternalServerError,
		}
	}

	// TODO: insert data into database
	// TODO: transfer var db ke sini
	return nil
}

func main() {

	log.Println("Opening connection to database ... ")
	db, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Ping database connection ... ")
	err = db.Ping()
	if err != nil {
		log.Println("Ping database connection: failure :(")
		log.Fatal(err)
	}
	log.Println("Ping database connection: success!")

	// register index handler
	// TODO: transfer db var to handler
	http.Handle("/", apiHandler(index))
	http.Handle("/users", apiHandler(users))

	// server listener
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
