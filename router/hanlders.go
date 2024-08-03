package router

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"unicode"

	DB "github.com/tomaslobato/mailist/db"
	"github.com/tomaslobato/mailist/models"
	"github.com/tomaslobato/mailist/utils"
)

func SendEmail(w http.ResponseWriter, r *http.Request, db *sql.DB, appPwd string, adminCode string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", 400)
	}
	defer r.Body.Close()

	var sendReq models.SendEmailReq
	json.Unmarshal(body, &sendReq)

	if sendReq.AdminCode == "" {
		http.Error(w, "Admin code not found", 400)
		return
	}
	if sendReq.AdminCode != adminCode {
		http.Error(w, "Admin code is wrong", 403)
		return
	}

	err = utils.SendEmail(db, appPwd, sendReq)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Sent emails successfully."))
}

func AddEmail(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", 400)
		return
	}
	defer r.Body.Close()

	var emailReq models.EmailReq
	json.Unmarshal(body, &emailReq)

	if emailReq.Email == "" {
		http.Error(w, "Email is required", 400)
		return
	}

	id, err := DB.AddUser(db, emailReq.Email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func GetEmailById(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	path := r.URL.Path
	id := path[len("/api/email/"):] //everything after /api/email/

	if !isNumeric(id) {
		http.Error(w, "id must be a number", 422)
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id must be a number", 422)
		return
	}

	user, err := DB.GetUserById(db, intId)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetEmails(w http.ResponseWriter, db *sql.DB) {
	users, err := DB.GetUsers(db)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

//

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
