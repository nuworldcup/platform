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
func (db *DB) IsValidTournament(tournamentId string) (bool, error) {
	if tournamentId == "" {
		return false, errors.New("tournamentId must not be empty")
	}
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM tournament WHERE tournament_id=$1)`, tournamentId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// This could probably be used in place of the query above (might be faster) then just check for the noRow err
// Gets one tournament given a tournament_name
func (db *DB) GetTournament(tournamentId string) (*types.Tournament, error) {
	if tournamentId == "" {
		return nil, errors.New("tournamentId must not be empty")
	}

	var t types.Tournament

	err := db.QueryRow(`SELECT tournament_id, tournament_name, tournament_timestamp, registration_timestamp FROM tournament WHERE tournament_id=$1`, tournamentId).Scan(&t.TournamentId, &t.TournamentName, &t.TournamentTime, &t.RegistrationTime)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
