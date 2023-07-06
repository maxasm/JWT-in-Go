package main

import (
	"fmt"
	mux "github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Printf("hey .....\n")
	})
	http.Handle("/", r)
}
