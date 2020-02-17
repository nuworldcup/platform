package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rojaswestall/platform/dblib"
	"github.com/rojaswestall/platform/lib"
)

func RegistrationOpenHandler(db *dblib.DB, name string) (bool, error) {

	if exists, err := db.IsValidTournament(name); err != nil {
		return false, err
	} else if !exists {
		msg := fmt.Sprintf("invalid tournament_name")
		return false, &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	if tournament, err := db.GetTournament(name); err != nil {
		return false, err
	} else {
		now := time.Now()
		// loc, _ := time.LoadLocation("Asia/Shanghai")
		// now := time.Now().In(loc)
		if now.Before(tournament.RegistrationTime) {
			return false, nil
		}
	}
	// Need way to return good message
	return true, nil
}
