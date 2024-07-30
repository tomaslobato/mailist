package DB

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID    int
	Email string
}

func QueryUsers(db *sql.DB) ([]User, error) {
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
