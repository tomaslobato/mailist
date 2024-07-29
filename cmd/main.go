package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		DB.QueryUsers(db)
		w.Write([]byte("opusk"))
	})
	http.ListenAndServe(":3000", mux)
}
