package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	DB "github.com/tomaslobato/mailist/db"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load()
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
		users, err := DB.QueryUsers(db)
		if err != nil {
			log.Fatal(err)
		}

		var userStrings []string
		for _, user := range users {
			userStrings = append(userStrings, fmt.Sprintf("%d: %s\n", user.ID, user.Email))
		}

		uString := strings.Join(userStrings, "\n")
		w.Write([]byte(uString))
	})
	http.ListenAndServe(":3000", mux)
}
