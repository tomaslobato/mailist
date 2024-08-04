package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tomaslobato/mailist/router"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Failed to load .env.local file: %v", err)
	}

	tursoUrl := os.Getenv("TURSO_URL")
	appPwd := os.Getenv("GMAIL_APP_PASSWORD")
	adminCode := os.Getenv("ADMIN_CODE")

	db, err := sql.Open("libsql", tursoUrl)
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err = http.ListenAndServe(":"+port, router.SetupRoutes(db, appPwd, adminCode))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
