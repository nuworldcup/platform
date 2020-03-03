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
	// upsert
	// this does not work because there is no unique or exclusion constraint
	// matching the ON CONFLICT specification atm
	err := db.QueryRow(`INSERT INTO team(team_name) VALUES($1) ON CONFLICT (team_name) DO NOTHING RETURNING (team_id)`, t.Name).Scan(&teamId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, errors.New("team already exists with given name")
		}
		return teamId, err
	}
	return teamId, nil
}

// Creates a new player
// Returns an error if team is invalid or the tx fails
func (db *DB) CreatePlayer(p *types.Player) (sql.Result, error) {
	if p == nil {
		return nil, errors.New("player required")
	}
	res, err := db.Exec(`INSERT INTO player(first_name, last_name, email, club) VALUES($1, $2, $3, $4)`, p.FirstName, p.LastName, p.Email, p.Club)
	return res, err
}

////////////////////////////////////////////////////
/////////////////// TRANSACTIONS ///////////////////
////////////////////////////////////////////////////

// Creates a new team and returns the teamId
// Returns an error if team is invalid or the tx fails
func (tx *Tx) CreateTeam(tournamentId string, t *types.Team) (int, error) {
	var teamId int
	// Need to do things to validate the input
	if t == nil {
		return teamId, errors.New("teamId required")
	}
	if t.Name == "" {
		return teamId, errors.New("team_name must not be empty")
	}
	// get teams with the same names, make sure that none exist
	// upsert

	// This should be null if there is no team in the tournament
	// SELECT team.team_id FROM
	// team JOIN team_tournament ON
	// team.team_id = team_tournament.team_id;

	// Gets all teams with specific name
	// (SELECT * FROM team t WHERE t.team_name = $1)

	// Gets all team_tournament entries with specific tournamentId
	// (SELECT * FROM team_tournament t WHERE team_tournament = $2)

	// insert into orders
	// select col1, col2, col3
	// 	from srctable
	// where col3 is not null
	// err := db.QueryRow(`
	// 	INSERT INTO team(team_name)
	// 	VALUES($1)
	// 	ON CONFLICT (team_name) DO NOTHING RETURNING (team_id)
	// 	`,
	// 	t.Name
	// ).Scan(&teamId)

	// err := db.QueryRow(`
	// 	INSERT INTO team(team_name)
	// 	VALUES($1, $2)
	// 	ON CONFLICT (team_name) DO NOTHING RETURNING (team_id)
	// 	`,
	// 	t.Name,
	// 	tournamentId
	// ).Scan(&teamId)

	// insert := `
	// 	INSERT INTO team(team_name) VALUES($1)

	// 	SELECT team_id
	// 	FROM (SELECT * FROM team t WHERE t.team_name = $1) filtered_teams
	// 	JOIN (SELECT * FROM team_tournament t WHERE tournament_id = $2) filtered_team_tournament
	// 	ON filtered_teams.team_id = filtered_team_tournament.team_id;`

	insert := `
		INSERT INTO team(team_name)
		VALUES($1)
		RETURNING team_id
	`

	err := tx.QueryRow(insert, t.Name).Scan(&teamId)

	if err != nil {
		return teamId, err
	}
	return teamId, nil

	// upsert
	// this does not work because there is no unique or exclusion constraint
	// matching the ON CONFLICT specification atm
	// insert := `
	// 	INSERT INTO team_tournament(team_id, tournament_id, team_tournament_name)
	// 	VALUES($1, $2, $3)
	// `

	// err = db.QueryRow(insert, t.Name).Scan(&teamId)
	// if err != nil {
	// 	if err.Error() == "sql: no rows in result set" {
	// 		return 0, errors.New("team already exists with given name")
	// 	}
	// 	return teamId, err
	// }
	// return teamId, nil
}

// Adds a team to a tournament and returns the teamId and tournamentId
// Returns an error if the teamId or tournamentName is invalid or if the tx fails
func (tx *Tx) AssignTeamToTournament(teamId int, teamName string, tournamentId string) (sql.Result, error) {
	// Need to do things to validate the input
	if teamId == 0 {
		return nil, errors.New("teamId required")
	}
	if tournamentId == "" {
		return nil, errors.New("tournamentName required")
	}

	insert := `
		INSERT INTO team_tournament(team_id, tournament_id, team_tournament_name) 
		VALUES($1, $2, $3)
	`

	res, err := tx.Exec(insert, teamId, tournamentId, teamName)

	return res, err
}

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
