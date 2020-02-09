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
func (db *DB) CreateTeamIfNotExists(t *types.Team) (int, error) {
	var teamId int
	// Need to do things to validate the input
	if t == nil {
		return teamId, errors.New("team required")
	}
	if t.Name == "" {
		return teamId, errors.New("team_name must not be empty")
	}
	if t.Tournament == "" {
		return teamId, errors.New("tournament_type must not be empty")
	}
	// upsert
	err := db.QueryRow(`INSERT INTO team(team_name, fk_tournament_name) VALUES($1, $2) ON CONFLICT (team_name, fk_tournament_name) DO NOTHING RETURNING (team_id)`, t.Name, t.Tournament).Scan(&teamId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, errors.New("team already exists with given name")
		}
		return teamId, err
	}
	return teamId, nil
}

////////////////////////////////////////////////////
/////////////////// TRANSACTIONS ///////////////////
////////////////////////////////////////////////////

// Creates a new player
// Returns an error if player is invalid or the tx fails.
func (tx *Tx) CreatePlayer(p *types.Player) (int, error) {
	var playerId int
	// Need to do things to validate the input
	if p == nil {
		return playerId, errors.New("player required")
	}
	err := tx.QueryRow(`INSERT INTO player(first_name, last_name, email, club) VALUES($1, $2, $3, $4) RETURNING player_id`, p.FirstName, p.LastName, p.Email, p.Club).Scan(&playerId)
	return playerId, err
}

// Creates a new player
// Returns an error if player is invalid or the tx fails.
func (tx *Tx) CreateCaptain(c *types.Captain) (int, error) {
	var captainId int
	// Need to do things to validate the input
	if c == nil {
		return captainId, errors.New("captain required")
	}
	err := tx.QueryRow(`INSERT INTO player(first_name, last_name, email, phone, club) VALUES($1, $2, $3, $4, $5) RETURNING player_id`, c.FirstName, c.LastName, c.Email, c.PhoneNumber, c.Club).Scan(&captainId)
	return captainId, err
}

// Adds a player to a team
// Returns an error if the player_id or team_id is invalid or if the tx fails
func (tx *Tx) AssignPlayerToTeam(tId int, pId int) (sql.Result, error) {
	// Need to do things to validate the input
	if tId == 0 {
		return nil, errors.New("team_id required")
	}
	if pId == 0 {
		return nil, errors.New("player_id required")
	}
	res, err := tx.Exec(`INSERT INTO team_player(team_id, player_id, captain) VALUES($1, $2, $3)`, tId, pId, false)
	return res, err
}

// Adds a captain to a team
// Returns an error if the captain_id or team_id is invalid or if the tx fails
func (tx *Tx) AssignCaptainToTeam(tId int, cId int) (sql.Result, error) {
	// Need to do things to validate the input
	if tId == 0 {
		return nil, errors.New("team_id required")
	}
	if cId == 0 {
		return nil, errors.New("player_id required")
	}
	res, err := tx.Exec(`INSERT INTO team_player(team_id, player_id, captain) VALUES($1, $2, $3)`, tId, cId, true)
	return res, err
}
