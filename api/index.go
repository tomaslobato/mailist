package api

//needed for vercel deployment

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
var adminCode string

func init() {
	dbUrl := os.Getenv("TURSO_URL")
	appPwd = os.Getenv("GMAIL_APP_PASSWORD")
	adminCode = os.Getenv("ADMIN_CODE")

	var err error
	db, err = sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %s", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.SetupRoutes(db, appPwd, adminCode).ServeHTTP(w, r)
}
