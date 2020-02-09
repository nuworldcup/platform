package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rojaswestall/platform/api"
	"github.com/rojaswestall/platform/dblib"
	"github.com/rojaswestall/platform/migrate"

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

// TODO :: add way to pass nuwc data without have to split up nuwc.db, nuwc.spreadsheetId etc.
func (nuwc *NUWCData) registerHandler(w http.ResponseWriter, r *http.Request) {
	// will need to add spreadsheetId(s), token, maybe email creds, slack
	err := api.RegisterHandler(w, r, nuwc.db)
	if err != nil {
		// internal error
		fmt.Println("There was an error, still need to implement a graceful handler")
		// http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
