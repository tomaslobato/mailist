package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tomaslobato/mailist/router"
)

var db *sql.DB
var appPwd string

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load .env.local file")
		os.Exit(1)
	}

	dbUrl := os.Getenv("TURSO_URL")
	appPwd = os.Getenv("GMAIL_APP_PASSWORD")

	db, err = sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.SetupRoutes(db, appPwd).ServeHTTP(w, r)
}
