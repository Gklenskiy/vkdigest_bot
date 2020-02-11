package models

import log "github.com/go-pkgz/lgr"

// UserSource model
type UserSource struct {
	UserID int
	Source string
}

// GetSources for user
func GetSources(userID int) ([]string, error) {
	rows, err := db.Query("SELECT source FROM user_sources WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := make([]string, 0)
	var source string
	for rows.Next() {
		err := rows.Scan(&source)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sources, nil
}

// CreateSource for user
func CreateSource(userID int, source string) error {
	stmt, err := db.Prepare("INSERT INTO user_sources (user_id, source) VALUES  ($1, $2) " +
		"ON CONFLICT " +
		"ON CONSTRAINT user_sources_user_id_source_key " +
		"DO NOTHING;")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(userID, source)
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
