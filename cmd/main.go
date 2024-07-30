package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	DB "github.com/tomaslobato/mailist/db"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type EmailReq struct {
	Email string `json:"email"`
}

func main() {
	err := godotenv.Load(".env.local")
	mux := http.NewServeMux()
	token := os.Getenv("TURSO_AUTH_TOKEN")
	url := fmt.Sprintf("libsql://mailist-tomaslobato.turso.io?authToken=%s", token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
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

	http.ListenAndServe(":3000", mux)
}
