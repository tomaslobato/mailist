package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tomaslobato/mailist/router"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load .env.local file")
		os.Exit(1)
	}

	tursoUrl := os.Getenv("TURSO_URL")
	appPwd := os.Getenv("GMAIL_APP_PASSWORD")
	adminCode := os.Getenv("ADMIN_CODE")

	db, err := sql.Open("libsql", tursoUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %s", err)
		os.Exit(1)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.ListenAndServe(":"+port, router.SetupRoutes(db, appPwd, adminCode))
}
