package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rojaswestall/platform/gtools"
	"github.com/rojaswestall/platform/migrate"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home : )")
}

// Add omit empty so unmarshalling json won't require that we have things like
// id or phone number

type Person struct {
	Id          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"number"`
}

type Ref struct {
	Person
}

type Player struct {
	Person
}

type Coach struct {
	Person
}

type Team struct {
	Id       string   `json:"id"`
	Captains []Player `json:"captains"`
	Players  []Player `json:"players"` // have this be a list of ids
	Name     string   `json:"name"`
	Wins     int      `json:"wins"`
	Losses   int      `json:"losses"`
	Draws    int      `json:"draws:"`
	Icon     string   `json:"icon"` // url of icon or flag to represent team
}

type GameState int

const (
	unknown GameState = iota
	FirstHalf
	SecondHalf
	NotStarted
	HalfTime
	ExtraTime
	OverTime
	Complete
	sentinel
)

func (s GameState) isValid() bool {
	return s > unknown && s < sentinel
}

// to get the string representations of the gameState
func (s GameState) String() string {
	states := [...]string{
		"First Half",
		"Second Half",
		"Not Started",
		"Half Time",
		"Extra Time",
		"Overtime",
		"Complete"}
	if !s.isValid() {
		return "Unknown"
	}
	return states[s]
}

type Game struct {
	Id            string `json:"id"`
	HomeTeam      Team
	AwayTeam      Team
	HomeTeamScore string
	AwayTeamScore string
	Date          string
	Time          string
	Venue         string
	GameState     GameState // could be something like first half, second half, overtime ...
	// Refs          []string  // ids of the refs
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
	TeamName string `json:"team_name"`
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	// TODO:: Use AWS secrets to set spreadsheetId for sheets
	spreadsheetId := ""
	// TODO:: Use AWS secrets to get credentials/token

	var info RegistrationInfo

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = gtools.AddSheet(spreadsheetId, info.TeamName)
	if err != nil {
		log.Fatal(err)
	}

	// respond back to the client that the team was registered
	fmt.Fprintf(w, "Team Info: %+v", info)
	fmt.Println("register endpoint was hit")
}

func main() {
	// create db instance
	// TODO:: Use AWS secrets to get username and password
	db, err := sql.Open("postgres", "postgres://nuwcuser:password@localhost:5432/nuwc?sslmode=disable")
	// Want sslmode to be enable as some point, for now disable
	if err != nil {
		log.Fatal(err)
	}
	// migrate
	migrate.Migrate(db)

	// Create router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/serveWs", serveWs)
	router.HandleFunc("/register", handleRegister).Methods("POST")
	//add any new endpoints here

	// Start listening
	fmt.Println("The server is listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
