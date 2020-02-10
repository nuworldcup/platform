package types

import "time"

type Person struct {
	Id        *string `json:"id,omitempty"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
}

type Player struct {
	Person
	Club *bool `json:"club"`
}

type Captain struct {
	Person
	PhoneNumber *string `json:"phone"`
	Club        *bool   `json:"club"`
}

type Ref struct {
	Person
}

type Admin struct {
	Person
}

type Team struct {
	Id         string   `json:"id"`
	Captains   []Player `json:"captains"`
	Players    []Player `json:"players"` // have this be a list of ids
	Name       string   `json:"name"`
	Wins       int      `json:"wins"`
	Losses     int      `json:"losses"`
	Draws      int      `json:"draws:"`
	Icon       string   `json:"icon"` // url of icon or flag to represent team
	Tournament string   `json:"tournament_name"`
}

// there is some way to get this from the db and this will not be needed
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

type Tournament struct {
	TournamentName   string    `json:"tournament_name"`
	DisplayName      string    `json:"display_name"`
	TournamentTime   time.Time `json:"tournament_timestamp"`
	RegistrationTime time.Time `json:"registration_timestamp"`
}

// Probably need to create a small db table for this too
type RegistrationType string

const (
	IndividualReg RegistrationType = "individual"
	TeamReg       RegistrationType = "team"
)

func (rt RegistrationType) isValid() bool {
	return true
}
