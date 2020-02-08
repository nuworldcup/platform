package dblib

import (
	"database/sql"
	"errors"

	"github.com/rojaswestall/platform/types"
)

////////////////////////////////////////////////////
//////////////////// DB ACTIONS ////////////////////
////////////////////////////////////////////////////

// Creates a new team
// Returns an error if team is invalid or the tx fails
func (db *DB) CreateTeamIfNotExists(t *types.Team) (sql.Result, error) {
	// Need to do things to validate the input
	if t == nil {
		return nil, errors.New("team required")
	}
	if t.Name == "" {
		return nil, errors.New("team_name must not be empty")
	}
	if t.Tournament == "" {
		return nil, errors.New("tournament_type must not be empty")
	}
	// upsert
	res, err := db.Exec(`INSERT INTO team(team_name, fk_tournament_name) VALUES($1, $2) ON CONFLICT (team_name, fk_tournament_name) DO NOTHING`, t.Name, t.Tournament)
	return res, err
}

////////////////////////////////////////////////////
/////////////////// TRANSACTIONS ///////////////////
////////////////////////////////////////////////////

// Creates a new player
// Returns an error if player is invalid or the tx fails.
func (tx *Tx) CreatePlayer(p *types.Player) error {
	// Need to do things to validate the input
	if p == nil {
		return errors.New("player required")
	}
	_, err := tx.Exec(`INSERT INTO player(first_name, last_name, email, club) VALUES($1, $2, $3, $4)`, p.FirstName, p.LastName, p.Email, p.Club)
	return err
}

// Creates a new player
// Returns an error if player is invalid or the tx fails.
func (tx *Tx) CreateCaptain(c *types.Captain) error {
	// Need to do things to validate the input
	if c == nil {
		return errors.New("captain required")
	}
	_, err := tx.Exec(`INSERT INTO player(first_name, last_name, email, phone, club) VALUES($1, $2, $3, $4, $5)`, c.FirstName, c.LastName, c.Email, c.PhoneNumber, c.Club)
	return err
}

// Adds a player to a team
// Returns an error if the player_id or team_id is invalid or if the tx fails
func (tx *Tx) AssignPlayerToTeam(tId string, pId string) (sql.Result, error) {
	// Need to do things to validate the input
	if tId == "" {
		return nil, errors.New("team_id required")
	}
	if pId == "" {
		return nil, errors.New("player_id required")
	}
	res, err := tx.Exec(`INSERT INTO team(team_id, player_id, captain) VALUES($1, $2, $3)`, tId, pId, false)
	return res, err
}

// Adds a captain to a team
// Returns an error if the captain_id or team_id is invalid or if the tx fails
func (tx *Tx) AssignCaptainToTeam(tId string, cId string) (sql.Result, error) {
	// Need to do things to validate the input
	if tId == "" {
		return nil, errors.New("team_id required")
	}
	if cId == "" {
		return nil, errors.New("player_id required")
	}
	res, err := tx.Exec(`INSERT INTO team(team_id, player_id, captain) VALUES($1, $2, $3)`, tId, cId, true)
	return res, err
}
