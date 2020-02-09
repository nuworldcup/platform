package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rojaswestall/platform/dblib"
	"github.com/rojaswestall/platform/gtools"
	"github.com/rojaswestall/platform/lib"
	"github.com/rojaswestall/platform/types"
)

type RegistrationInfo struct {
	RegistrationType types.RegistrationType `json:"registration_type"`          // cannot be empty
	TournamentType   string                 `json:"tournament_type"`            // cannot be empty
	NamePreferences  []string               `json:"name_preferences,omitempty"` // can have as many preferences as they want
	Captains         []types.Captain        `json:"captains,omitempty"`         // can have as many captains as they want
	Players          []types.Player         `json:"players"`                    // cannot be empty
}

// Creat a map from country name to 3 digit code that tells us what flag to use OR keep that on the front end

func convertBoolToStringForSheets(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// adds a team to the db, returns the teamName and any error
// This can definitely be split up a lot better
func findAvailableNameAndAddTeamToDb(info RegistrationInfo, db *dblib.DB) (string, error) {
	var teamName string
	var teamId int

	// add the team with unique team name
	for i := 0; i <= len(info.NamePreferences); i++ {

		// all the preferences have already been tried
		if i == len(info.NamePreferences) {
			// at this point we could also return a response to the client that all the names have already been taken
			added := true
			count := 1
			// try default names until it works
			for added {
				id, err := db.CreateTeamIfNotExists(&types.Team{Tournament: info.TournamentType, Name: "Team " + strconv.Itoa(count)})
				if err != nil {
					if err.Error() == "team already exists with given name" {
						count++
						continue
					}
					// insert failed: something wrong with the db, tell the client we're messed UP
					return "", err
				}

				// success
				teamId = id
				teamName = "Team " + strconv.Itoa(count)
				added = false
			}
			break
		}

		// check the db to see if the name already exists
		exists, err := db.TeamExists(info.NamePreferences[i])
		if err != nil {
			// something wrong with the db, tell the client we're messed UP
			return "", err
		} else if exists {
			// the name already exists, try their next preference
			continue
		}

		// the name doesn't exist yet. Create the team
		id, err := db.CreateTeamIfNotExists(&types.Team{Tournament: info.TournamentType, Name: info.NamePreferences[i]})
		if err != nil {
			if err.Error() == "team already exists with given name" {
				// there was no update because another team registered at basically the same time and got the name first, try again
				continue
			}
			// insert failed: something wrong with the db, tell the client we're messed UP
			return "", err
		}
		teamId = id
		teamName = info.NamePreferences[i]
		break
	}

	// add all the players to the db
	for i := 0; i < len(info.Players); i++ {
		tx, _ := db.Begin()
		playerId, err := tx.CreatePlayer(&info.Players[i])
		if err != nil {
			// something went wrong creating the player
			return "", err
		}
		tx.AssignPlayerToTeam(teamId, playerId)
		tx.Commit()
	}
	for i := 0; i < len(info.Captains); i++ {
		tx, _ := db.Begin()
		playerId, err := tx.CreateCaptain(&info.Captains[i])
		if err != nil {
			// something went wrong creating the player
			return "", err
		}
		tx.AssignCaptainToTeam(teamId, playerId)
		tx.Commit()
	}
	return teamName, nil
}

func addTeamToSheets(info RegistrationInfo, spreadsheetId string, teamName string) error {
	err := gtools.AddSheet(spreadsheetId, teamName)
	if err != nil {
		// try for the top name preferences, and if always an error then assign default
		if err.Error() == "googleapi: Error 400: Invalid requests[0].addSheet: A sheet with the name \""+teamName+"\" already exists. Please enter another name., badRequest" {
			// BAD means there's a db issue we need to figure out
			return err
		}

		// Add error for too many API calls here

		// unknown error, internalServerError to client
		return err
	}

	// Add team values
	var values [][]interface{}

	// This could be MUCH more dynamic but fine for now
	headers := []interface{}{"First Name", "Last Name", "Email", "Club", "Phone", "Captain"}
	// Add headers
	values = append(values, headers)

	// add all players, captains first
	for i := 0; i < len(info.Captains); i++ {
		rowVals := []interface{}{info.Captains[i].FirstName, info.Captains[i].LastName, info.Captains[i].Email, convertBoolToStringForSheets(*info.Captains[i].Club), info.Captains[i].PhoneNumber, convertBoolToStringForSheets(true)}
		values = append(values, rowVals)
	}
	for i := 0; i < len(info.Players); i++ {
		rowVals := []interface{}{info.Players[i].FirstName, info.Players[i].LastName, info.Players[i].Email, convertBoolToStringForSheets(*info.Players[i].Club), "", convertBoolToStringForSheets(false)}
		values = append(values, rowVals)
	}

	err = gtools.AddSheetData(teamName, spreadsheetId, values)
	if err != nil {
		// Add error for too many API calls here

		// unknown error, internalServerError to client
		return err
	}

	return nil
}

/////// VALIDATE REQUEST ///////

func validatePerson(p types.Person) error {
	if p.FirstName == nil {
		msg := fmt.Sprintf("first_name required")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}
	if p.LastName == nil {
		msg := fmt.Sprintf("last_name required")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}
	if p.Email == nil {
		msg := fmt.Sprintf("email required")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}
	return nil
}

func validatePlayer(p types.Player) error {
	err := validatePerson(p.Person)
	if err != nil {
		return err
	}
	if p.Club == nil {
		msg := fmt.Sprintf("club required")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}
	return nil
}

func validateCaptain(c types.Captain) error {
	err := validatePerson(c.Person)
	if err != nil {
		return err
	}
	if c.PhoneNumber == nil {
		msg := fmt.Sprintf("phone required")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}
	if c.Club == nil {
		msg := fmt.Sprintf("club required")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}
	return nil
}

func nuwcRegistrationValidationIndividual(info RegistrationInfo, db *dblib.DB) error {
	if exists, err := db.IsValidTournament(info.TournamentType); err != nil {
		return err
	} else if !exists {
		msg := fmt.Sprintf("invalid tournament_type")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	if len(info.Captains) > 1 {
		msg := fmt.Sprintf("length of captains > 1")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	if len(info.Players) != 1 {
		msg := fmt.Sprintf("length of players != 1")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	if len(info.NamePreferences) > 1 {
		msg := fmt.Sprintf("length of name_preferences > 1")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	// validate that it's after the time to register

	return nil
}

func nuwcRegistrationValidationTeam(info RegistrationInfo, db *dblib.DB) error {
	if exists, err := db.IsValidTournament(info.TournamentType); err != nil {
		return err
	} else if !exists {
		msg := fmt.Sprintf("invalid tournament_type")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	if len(info.Captains) < 1 {
		msg := fmt.Sprintf("length of captains < 1")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	for i := 0; i < len(info.Captains); i++ {
		err := validateCaptain(info.Captains[i])
		if err != nil {
			return err
		}
	}

	if len(info.Players) < 8 {
		msg := fmt.Sprintf("length of players < 8")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	for i := 0; i < len(info.Players); i++ {
		err := validatePlayer(info.Players[i])
		if err != nil {
			return err
		}
	}

	if len(info.NamePreferences) < 3 {
		msg := fmt.Sprintf("length of name_preferences < 3")
		return &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	// validate that it's after the time to register

	return nil
}

////////////////////////////////
////////////////////////////////
////////////////////////////////

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *dblib.DB) error {
	// TODO:: Use AWS secrets to set spreadsheetId for sheets
	spreadsheetId := "1jDCdULFKmxmgCsJTJgqzKloCvnE85r8PyLvXDAlKLcA"
	// TODO:: Use AWS secrets to get credentials/token in gtools lib

	// Add everything to our db and then start the google sheets stuffs,
	// if there is a google sheets error don't return error to the client
	// Only return error to client if we get errors adding to the db
	// We can always get the data from the db

	var info RegistrationInfo
	err := lib.DecodeJSONBody(w, r, &info)
	if err != nil {
		return err
	}

	///////////////////////////////////////////////////
	///////////// VALIDATE REQUEST FORMAT /////////////
	///////////////////////////////////////////////////
	err = nuwcRegistrationValidationTeam(info, db)
	if err != nil {
		return err
	}

	////////////////////////////////////////////////////
	//////////////////// DB UPDATES ////////////////////
	////////////////////////////////////////////////////

	// might want to create a different function that tries to add team if
	// the name doesn't exist and just tells the client that it failed if
	// none of the names were available

	teamName, err := findAvailableNameAndAddTeamToDb(info, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	///////////////////////////////////////////////////
	////////////////// GOOGLE SHEETS //////////////////
	///////////////////////////////////////////////////

	// Try to keep API calls AS LOW AS POSSIBLE. We only get 500 per hour
	// Right now it takes 2 API calls to create a team, and 1 to add an individual

	// TODO:: Account for reaching API limit -- add cron (?) or a queue to add to
	//        sheet once we have waited an hour. Could check on error

	// Create new sheet for team
	// We can either alter the sheetName to have the tournament_type too
	// or we can use a different spreadsheetId for the two tournaments
	// For now assuming two different sheets
	err = addTeamToSheets(info, spreadsheetId, teamName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	////////////////////////////////////////////////////
	//////////////// SLACK NOTIFICATION ////////////////
	////////////////////////////////////////////////////

	// Notify the slack channel that the team registered

	////////////////////////////////////////////////////
	/////////////// EMAIL TO THE CAPTAIN ///////////////
	////////////////////////////////////////////////////

	// Email the captain that we received their application to register
	// with all the info they submitted

	////////////////////////////////////////////////////
	//////////// SUCCESS RESPONSE TO CLIENT ////////////
	////////////////////////////////////////////////////

	// respond back to the client that the team was registered
	// need a better response. Maybe the team name
	fmt.Fprintf(w, "Team Info: %+v", info)
	return nil
}
