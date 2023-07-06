package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"encoding/json"
	"fmt"
)

// the username and password
type UserModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = make([]UserModel, 10)

var eventLogger = log.New(os.Stdout, "[event] ", log.Ltime)

var fHandler = func(w http.ResponseWriter, r *http.Request) {
	// get the requested file path
	req_path := r.URL.Path

	eventLogger.Printf("%s Request on %s\n", r.Method, req_path)
	// '/' -> index.html
	if req_path == "/" {
		eventLogger.Printf("Serving index.html page\n")
		req_path = "/index.html"
	}

	// files are served from the static folder
	req_path = "./client/dist" + req_path

	// check if the path exists. If not server index.html
	_, stat_err := os.Stat(req_path)

	// serve index.html if the requested file does not exist
	if os.IsNotExist(stat_err) {
		eventLogger.Printf("%s does not exist -> index.html\n", req_path)
		req_path = "./client/dist/index.html"
	}
	
	// set the correct content-type
	switch filepath.Ext(req_path) {
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		case ".js":
			w.Header().Set("Content-Type", "text/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".ico":
		w.Header().Set("Content-Type", "image/vnd.microsoft.icon")
	}

	// read data from the given file
	f, err_f := os.Open(req_path)
	defer f.Close()

	if err_f != nil {
		log.Fatalf("%s\n", err_f)
	}

	fdata, err_fdata := io.ReadAll(f)
	if err_fdata != nil {
		log.Fatalf("%s\n", err_fdata)
	}
	w.Write(fdata)
}



var lHandler = func(w http.ResponseWriter, r *http.Request) {
	eventLogger.Printf("Request on /login")
	usermodel := UserModel{}
	payload, err_payload := io.ReadAll(r.Body)
	if err_payload != nil {
		log.Fatalf("%s\n", err_payload)
	}
	json.Unmarshal([]byte(payload),&usermodel)
	// chech if thew new user has is using an already used username
	for _, u := range users {
		if usermodel.Username == u.Username {
			if usermodel.Password == u.Password {
				// crate jwt token
				io.WriteString(w, fmt.Sprintf("Sucessful login: %s\n", usermodel.Username))
				return
			} else {
				break
			} 
		}
	}
	
	w.WriteHeader(http.StatusUnauthorized)
	io.WriteString(w, "You tried to login using wrong credentials.")
}


var sHandler = func(w http.ResponseWriter, r *http.Request) {
	eventLogger.Printf("Request on /signup")
	usermodel := UserModel{}
	payload, err_payload := io.ReadAll(r.Body)
	if err_payload != nil {
		log.Fatalf("%s\n", err_payload)
	}
	json.Unmarshal([]byte(payload),&usermodel)
	// chech if thew new user has is using an already used username
	for _, u := range users {
		if usermodel.Username == u.Username {
			w.WriteHeader(http.StatusConflict)
			io.WriteString(w, fmt.Sprintf("The username %s is already in use.", usermodel.Username))
			eventLogger.Printf("New user used an already used username %s\n", usermodel.Username)
			return
		}
	}
	
	// add the user if the username is unique
	users = append(users, usermodel)	
	io.WriteString(w, fmt.Sprintf("Successfully added the user %s\n", usermodel.Username))
	eventLogger.Printf("Added new user %s\n", usermodel.Username)
}

func main() {
	http.HandleFunc("/", fHandler)
	http.HandleFunc("/login", lHandler)	
	http.HandleFunc("/signup", sHandler)

	eventLogger.Printf("Server live on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
