package types

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
	Tournament string   `json:"tournament_type"`
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

// This type could be improved to be dymanic based on the year when the
// tournament is played. Right now this needs to be updated manually to
// account for many teams from different years being able to exist in the db
// Probably would want to create a small db table for this with migrations
// being added each year to make sure the type exists
type TournamentType string

const (
	Coed20   TournamentType = "coed2020"
	Womens20 TournamentType = "womens2020"
)

func (tt TournamentType) isValid() bool {
	return true
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
