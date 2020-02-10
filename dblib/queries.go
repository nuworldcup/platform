package dblib

import (
	"errors"

	"github.com/rojaswestall/platform/types"
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

// Checks if a tournament with a certain tournament_name exists
func (db *DB) IsValidTournament(tournamentName string) (bool, error) {
	if tournamentName == "" {
		return false, errors.New("tournamentName must not be empty")
	}
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM tournament WHERE tournament_name=$1)`, tournamentName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// This could probably be used in place of the query above (might be faster) then just check for the noRow err
// Gets one tournament given a tournament_name
func (db *DB) GetTournament(tournamentName string) (*types.Tournament, error) {
	if tournamentName == "" {
		return nil, errors.New("tournamentName must not be empty")
	}

	var t types.Tournament

	err := db.QueryRow(`SELECT tournament_name, display_name, tournament_timestamp, registration_timestamp FROM tournament WHERE tournament_name=$1`, tournamentName).Scan(&t.TournamentName, &t.DisplayName, &t.TournamentTime, &t.RegistrationTime)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
