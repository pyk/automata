package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	PORT = os.Getenv("AUTOMATA_SERVER_PORT")
)

// HelloWorld handle '/' request
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World! - automata.")
}

func main() {

	// rest api
	http.HandleFunc("/", HelloWorld)

	// server listener
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
