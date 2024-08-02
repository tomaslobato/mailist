package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/tomaslobato/mailist/router"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var db *sql.DB
var appPwd string

func init() {
	dbUrl := os.Getenv("TURSO_URL")
	appPwd = os.Getenv("GMAIL_APP_PASSWORD")

	db, err := sql.Open("libsql", dbUrl)
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
