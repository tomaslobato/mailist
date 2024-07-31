package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"unicode"

	"github.com/joho/godotenv"
	DB "github.com/tomaslobato/mailist/db"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gopkg.in/gomail.v2"
)

type EmailReq struct {
	Email string `json:"email"`
}

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func sendEmail(db *sql.DB, appPwd string) error {
	users, err := DB.GetUsers(db)
	if err != nil {
		return err
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, "tlobatodev@gmail.com", appPwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", "tlobatodev@gmail.com")
	m.SetHeader("Subject", "Test")
	m.SetBody("text/html", "this is a test for mailist-go.vercel.app")

	for _, user := range users {
		m.SetHeader("To", user.Email)
		err := d.DialAndSend(m)
		if err != nil {
			fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
			continue
		} else {
			fmt.Printf("Email sent to %s\n", user.Email)
		}
	}

	return nil
}

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load .env.local file")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	token := os.Getenv("TURSO_AUTH_TOKEN")
	appPwd := os.Getenv("GMAIL_APP_PASSWORD")

	tursoUrl := fmt.Sprintf("libsql://mailist-tomaslobato.turso.io?authToken=%s", token)

	db, err := sql.Open("libsql", tursoUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", tursoUrl, err)
		os.Exit(1)
	}
	defer db.Close()

	mux.HandleFunc("GET /emails", func(w http.ResponseWriter, r *http.Request) {
		users, err := DB.GetUsers(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	mux.HandleFunc("GET /emails/{id}", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		id := path[len("/emails/"):] //everything after /emails/

		if !isNumeric(id) {
			http.Error(w, "id must be a number", 422)
			return
		}

		intId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "id must be a number", 422)
			return
		}

		user, err := DB.GetUserById(db, intId)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("POST /emails", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var emailReq EmailReq
		json.Unmarshal(body, &emailReq)

		if emailReq.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		id, err := DB.AddUser(db, emailReq.Email)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	})

	mux.HandleFunc("POST /emails/send", func(w http.ResponseWriter, r *http.Request) {
		err := sendEmail(db, appPwd)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("Sent emails successfully."))
	})

	http.ListenAndServe(":3000", mux)
}
