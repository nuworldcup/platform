package types

type Country struct {
	CountryKey     string `json:"country"`
	DisplayName    string `json:"display_name"`
	TwoLetterIso   string `json:"two_letter_iso"`
	ThreeLetterIso string `json:"three_letter_iso"`
}

func IsSupportedCountry(c string) bool {
	switch c {
	case "antarctica", "kosovo":
		return false
	}
	return true
}
