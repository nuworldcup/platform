package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rojaswestall/platform/dblib"
	"github.com/rojaswestall/platform/migrate"
	"github.com/rojaswestall/platform/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Want to add secrets to this too
type NUWCData struct {
	db *dblib.DB
	// slack key
	// spreadsheet id
	// sheets token
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home : )")
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

type RegistrationInfo struct {
	RegistrationType types.RegistrationType `json:"registration_type"`          // cannot be empty
	TournamentType   types.TournamentType   `json:"tournament_type"`            // cannot be empty
	NamePreferences  []string               `json:"name_preferences,omitempty"` // can have as many preferences as they want
	Captains         []types.Captain        `json:"captains,omitempty"`         // can have as many captains as they want
	Players          []types.Player         `json:"players"`                    // cannot be empty
}

// Creat a map from country name to 3 digit code that tells us what flag to use OR keep that on the front end

func (nuwc *NUWCData) registerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:: Use AWS secrets to set spreadsheetId for sheets
	// spreadsheetId := "1jDCdULFKmxmgCsJTJgqzKloCvnE85r8PyLvXDAlKLcA"
	// TODO:: Use AWS secrets to get credentials/token in gtools lib

	// Add everything to our db and then start the google sheets stuffs,
	// if there is a google sheets error don't return error to the client
	// Only return error to client if we get errors adding to the db
	// We can always get the data from the db

	var info RegistrationInfo

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	///////////////////////////////////////////////////
	///////////// VALIDATE REQUEST FORMAT /////////////
	///////////////////////////////////////////////////

	// check for at least one captain
	// check for 8 players
	// check for 3 name preferences
	// check for reg type
	// check for tournament type
	// validate that it's after the time to register

	////////////////////////////////////////////////////
	//////////////////// DB UPDATES ////////////////////
	////////////////////////////////////////////////////

	var sheetName string
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
				id, err := nuwc.db.CreateTeamIfNotExists(&types.Team{Tournament: info.TournamentType, Name: "Team " + strconv.Itoa(count)})
				if err != nil {
					// insert failed: something wrong with the db, tell the client we're messed UP
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				if id == 0 {
					count++
					continue
				}
				// success
				teamId = id
				sheetName = "Team " + strconv.Itoa(count)
				added = false
			}
			break
		}

		// check the db to see if the name already exists
		exists, err := nuwc.db.TeamExists(info.NamePreferences[i])
		fmt.Println("Team "+info.NamePreferences[i]+" -- exists: %t", exists)
		if err != nil {
			// something wrong with the db, tell the client we're messed UP
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if exists {
			// the name already exists, try their next preference
			continue
		}

		// the name doesn't exist yet. Create the team
		id, err := nuwc.db.CreateTeamIfNotExists(&types.Team{Tournament: info.TournamentType, Name: info.NamePreferences[i]})
		if err != nil {
			// insert failed: something wrong with the db, tell the client we're messed UP
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if id == 0 {
			// there was no update because another team registered at basically the same time and got the name first, try again
			continue
		}
		teamId = id
		sheetName = info.NamePreferences[i]
		break
	}
	fmt.Println("teamName: %v", sheetName)

	// add all the players to the db
	for i := 0; i < len(info.Players); i++ {
		tx, _ := nuwc.db.Begin()
		playerId, err := tx.CreatePlayer(&info.Players[i])
		if err != nil {
			// something went wrong creating the player
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tx.AssignPlayerToTeam(teamId, playerId)
		tx.Commit()
	}
	for i := 0; i < len(info.Captains); i++ {
		tx, _ := nuwc.db.Begin()
		playerId, err := tx.CreateCaptain(&info.Captains[i])
		if err != nil {
			// something went wrong creating the player
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tx.AssignCaptainToTeam(teamId, playerId)
		tx.Commit()
	}

	///////////////////////////////////////////////////
	////////////////// GOOGLE SHEETS //////////////////
	///////////////////////////////////////////////////

	// Create new sheet
	// err = gtools.AddSheet(spreadsheetId, sheetName)
	// if err != nil {
	// 	// try for the top name preferences, and if always an error then assign default
	// 	if err.Error() == "googleapi: Error 400: Invalid requests[0].addSheet: A sheet with the name \""+sheetName+"\" already exists. Please enter another name., badRequest" {
	// 		// BAD means there's a db issue we need to figure out
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	}
	// 	// unknown error, internalServerError to client
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	// // add values to the new sheet
	// // Need function to do this
	// err = gtools.AddSheetRow(sheetName, spreadsheetId, info.NamePreferences)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	////////////////////////////////////////////////////
	//////////////// SLACK NOTIFICATION ////////////////
	////////////////////////////////////////////////////

	// Notify the slack channel that the team registered

	////////////////////////////////////////////////////
	//////////// SUCCESS RESPONSE TO CLIENT ////////////
	////////////////////////////////////////////////////

	// respond back to the client that the team was registered
	fmt.Fprintf(w, "Team Info: %+v", info)
}

func main() {
	// create db instance
	// TODO:: Use AWS secrets to get username and password
	db, err := dblib.Open("postgres://nuwcuser:password@localhost:5432/nuwc?sslmode=disable")
	// Want sslmode to be enable as some point, for now disable
	if err != nil {
		log.Fatal(err)
	}

	// migrate
	// might want to move this to db.go at some point?
	// migrate should be change to return an error so it can be handled however
	migrate.Migrate(db.DB)

	// establish a connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	// can get password and keep them here as well
	// inject from aws secrets
	nuwc := &NUWCData{db: db}

	// Create router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/serveWs", serveWs)
	router.HandleFunc("/register", nuwc.registerHandler).Methods("POST")
	//add any new endpoints here

	// Start listening
	fmt.Println("The server is listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
