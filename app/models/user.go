package models

import (
	log "github.com/go-pkgz/lgr"
)

// User model
type User struct {
	UserID int
	Token  string
}

// AllUsers returns all records from Users table
func AllUsers() ([]*User, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.UserID, &user.Token)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// GetToken for user
func GetToken(userID int) (string, error) {
	var token string
	err := db.QueryRow("SELECT token FROM users WHERE user_id = $1", userID).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateOrUpdate user token
func CreateOrUpdate(userID int, token string) error {
	log.Printf("[DEBUG] CreateOrUpdate")
	stmt, err := db.Prepare("INSERT INTO users (user_id, token) VALUES  ($1, $2) " +
		"ON CONFLICT (user_id) " +
		"DO UPDATE SET token = EXCLUDED.token;")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(userID, token)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Affected = %d\n", rowCnt)
	return nil
}
