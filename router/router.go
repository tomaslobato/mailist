package router

import (
	"database/sql"
	"net/http"
)

func SetupRoutes(db *sql.DB, appPwd, adminCode string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /emails", func(w http.ResponseWriter, r *http.Request) { GetEmails(w, db) })
	mux.HandleFunc("GET /emails/{id}", func(w http.ResponseWriter, r *http.Request) { GetEmailById(w, r, db) })
	mux.HandleFunc("POST /emails", func(w http.ResponseWriter, r *http.Request) { AddEmail(w, r, db) })
	mux.HandleFunc("POST /emails/send", func(w http.ResponseWriter, r *http.Request) { SendEmail(w, r, db, appPwd, adminCode) })

	return mux
}
