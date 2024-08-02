package router

import (
	"database/sql"
	"net/http"

	"github.com/tomaslobato/mailist/handlers"
)

func SetupRoutes(db *sql.DB, appPwd string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /emails", func(w http.ResponseWriter, r *http.Request) { handlers.GetEmails(w, db) })
	mux.HandleFunc("GET /emails/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.GetEmailById(w, r, db) })
	mux.HandleFunc("POST /emails", func(w http.ResponseWriter, r *http.Request) { handlers.AddEmail(w, r, db) })
	mux.HandleFunc("POST /emails/send", func(w http.ResponseWriter, r *http.Request) { handlers.SendEmail(w, r, db, appPwd) })

	return mux
}
