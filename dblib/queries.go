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

func (db *DB) GetCountries(tournamentId string) ([]types.Country, error) {
	// join the tournament
	if tournamentId == "" {
		return nil, errors.New("tournamentId must not be empty")
	}

	countries := []types.Country{}

	// still need to only get the teams for a specific tournament
	rows, err := db.Query(`SELECT c.country_pkey, display_name, two_letter_iso, three_letter_iso FROM country c LEFT JOIN team t ON t.team_name = c.country_pkey WHERE t.team_name IS NULL`)
	if err != nil {
		return nil, err
	}

	// Iterate through results of query and store them in countries
	defer rows.Close()
	for rows.Next() {
		var c types.Country
		err = rows.Scan(&c.CountryKey, &c.DisplayName, &c.TwoLetterIso, &c.ThreeLetterIso)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return countries, nil
}
