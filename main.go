package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	PORT = os.Getenv("AUTOMATA_SERVER_PORT")
)

// apiError define structure of API error
// TODO: API error response should return JSON
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
	// TODO: print response status
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

func main() {

	// register index handler
	http.Handle("/", apiHandler(index))

	// server listener
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
