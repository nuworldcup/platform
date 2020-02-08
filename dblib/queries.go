package dblib

import (
	"errors"
)

// Checks if a team with a certain team_name exists
func (db *DB) TeamExists(teamName string) (bool, error) {
	if teamName == "" {
		return false, errors.New("teamName must not be empty")
	}
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM team WHERE team_name=$1)`, teamName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
