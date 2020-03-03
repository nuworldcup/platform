package api

import (
	"fmt"
	"net/http"

	"github.com/rojaswestall/platform/dblib"
	"github.com/rojaswestall/platform/lib"
	"github.com/rojaswestall/platform/types"
)

func AvailableCountriesHandler(db *dblib.DB, name string) ([]types.Country, error) {

	if exists, err := db.IsValidTournament(name); err != nil {
		return nil, err
	} else if !exists {
		msg := fmt.Sprintf("invalid tournament_name")
		return nil, &lib.MalformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	countries, err := db.GetCountries(name)
	if err != nil {
		return nil, err
	}

	filtered := []types.Country{}
	for i, v := range countries {
		// filter unsupported countries
		if types.IsSupportedCountry(v.CountryKey) {
			filtered = append(filtered, countries[i])
		}
	}
	// Need way to return good message
	return filtered, nil
}
