package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

// the key for JWT signing 
var SECRET_KEY = generateKey() 

var users = make([]Credentials, 10)

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
	user := Credentials{}
	payload, err_payload := io.ReadAll(r.Body)
	if err_payload != nil {
		log.Fatalf("%s\n", err_payload)
	}
	json.Unmarshal([]byte(payload),&user)
	// chech if thew new user has is using an already used username
	for _, u := range users {
		if u.Username == user.Username {
			if compHash(u.Password,user.Password) {
				// create JWT 
				// the jwt expires after 5 minutes
				tokenExp := time.Now().Add(5 * time.Minute)
				// create the claims -> JWT payload
				claims := Claims{
					UserId: u.Id,
					Username: u.Username,	
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(tokenExp),
					},
				}
				// declare the token + the signing algorithm
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				// create the JWT string
				token_string, err_token_string := token.SignedString(SECRET_KEY)
				eventLogger.Printf("Generated new token_string -> %s\n", token_string)
				if err_token_string != nil {
					eventLogger.Printf("%s\n", err_token_string)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				
				// set the client cookie for token to the JWT
				http.SetCookie(w, &http.Cookie{
					Name: "token",
					Value: token_string,
					Expires: tokenExp,
				})
				
				io.WriteString(w, fmt.Sprintf("Sucessful login: %s\n", user.Username))
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
	user := Credentials{}
	payload, err_payload := io.ReadAll(r.Body)
	if err_payload != nil {
		log.Fatalf("%s\n", err_payload)
	}
	json.Unmarshal([]byte(payload),&user)
	// chech if thew new user has is using an already used username
	for _, u := range users {
		if user.Username == u.Username {
			w.WriteHeader(http.StatusConflict)
			io.WriteString(w, fmt.Sprintf("The username %s is already in use.", user.Username))
			eventLogger.Printf("New user used an already used username %s\n", user.Username)
			return
		}
	}
	
	// generate a new user ID	
	userId := uuid.New().String()
	user.Id = userId 
	// hash the password of the user
	pswd, err_pswd := genHash(user.Password)
	if err_pswd != nil {
		eventLogger.Fatalf("%s\n", err_pswd)	
	}
	user.Password = pswd	
	// append the users
	users = append(users, user)

	io.WriteString(w, fmt.Sprintf("Successfully added the user %s\n", user.Username))
	eventLogger.Printf("Added new user %s\n", user.Username)
}

var aHandler = func(w http.ResponseWriter, r *http.Request) {
	eventLogger.Printf("Request on /api#app\n")
	// read the token cookie -> jwt
	cookie, err_cookie := r.Cookie("token")	
	
	if err_cookie != nil {
		if err_cookie == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)	
			return
		}
		// for any other error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the JWT from the cookie
	token_string := cookie.Value
	
	// reconstruct 'Claims'
	claims := &Claims{}
		
	token, err_token := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token)(interface{}, error){
		return SECRET_KEY,nil
	})

	if err_token != nil {
		if err_token == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	if !token.Valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// Finally -> This user is authorized using JWT set in cookies
	io.WriteString(w,fmt.Sprintf("Username: %s\nId: %s\n", claims.Username, claims.UserId))
}

func main() {
	http.HandleFunc("/", fHandler)
	http.HandleFunc("/api_login", lHandler)	
	http.HandleFunc("/api_signup", sHandler)
	http.HandleFunc("/api_app", aHandler)

	eventLogger.Printf("Server live on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
