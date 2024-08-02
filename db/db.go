package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID    int
	Email string
}

func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM emails")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return users, nil
}

func AddUser(db *sql.DB, email string) (int, error) {
	result, err := db.Exec("INSERT INTO emails (email) VALUES (?)", email)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new user:%v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %v", err)
	}

	return int(id), nil
}

func GetUserById(db *sql.DB, id int) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, email FROM emails WHERE id = ?", id).Scan(&user.ID, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with id: %d", id)
		}
		return nil, fmt.Errorf("Error getting email by ID: %s", err)
	}

	return &user, nil
}
