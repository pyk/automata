package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	Tag     string `json:"-"`
	Error   error  `json:"-"`
	Message string `json:"error"`
	Code    int    `json:"code"`
}

// ApiHandler global API mux
type ApiHandler struct {
	DB      *sql.DB
	Handler func(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError
}

func (api ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// add header on every response
	w.Header().Add("Server", "Automata/0.1")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// if handler return an &apiError
	err := api.Handler(w, r, api.DB)
	if err != nil {
		// http log
		log.Printf("%s %s %s [%s] %s", r.RemoteAddr, r.Method, r.URL, err.Tag, err.Error)

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

// indexHandler handle '/' request
func indexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError {

	// response "404 not found" on every undefined
	// URL pattern handler
	if r.URL.Path != "/" {
		return &apiError{
			"indexHandler url",
			errors.New("Not Found"),
			"Not Found",
			http.StatusNotFound,
		}
	}

	err := db.Ping()
	if err != nil {
		return &apiError{
			"indexHandler db ping",
			err,
			"internal server error",
			http.StatusInternalServerError,
		}
	}
	log.Println("success ping database")

	return nil
}

// User type
// TODO: add more neede field
type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// users handle '/users' request
// TODO: validate request header must Accept: application/json
func usersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError {

	switch r.Method {
	case "GET":
		return usersHandlerGET(w, r, db)
	case "POST":
		return usersHandlerPOST(w, r, db)
	default:
		return &apiError{
			"usersHandler default case",
			errors.New("Not Found"),
			"Not Found",
			http.StatusNotFound,
		}
	}
	return nil
}

// usersHandlerGET handle 'GET' request on '/users'
// TODO: query data from database
func usersHandlerGET(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError {

	return nil
}

// usersHandlerPOST handle 'POST' request on '/users'
func usersHandlerPOST(w http.ResponseWriter, r *http.Request, db *sql.DB) *apiError {

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

	// TODO: aku pengen API nya kaya gini
	// api.NewDB(db) -> register database
	// api.Handler(indexHandler) -> register handler
	// api.Handler(user) -> register handler
	// akan tetapi DB bisa diakses oleh setiap
	// handler

	// Register handler
	http.Handle("/", ApiHandler{db, indexHandler})
	http.Handle("/users", ApiHandler{db, usersHandler})

	// server listener
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
