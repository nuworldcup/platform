package gtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "gtools/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

////// The above methods are provided by google (but we modified a lil) for quickstart use ^
////// The below methods were made to make using google sheets easy : )
////// Add more as needed!

// Set up a service using credentials to perform sheets operations
func getService() *sheets.Service {
	b, err := ioutil.ReadFile("gtools/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// if changing this line, delete token.json / reset in AWS secrets
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}
	return srv
}

// Add sheet to a spreadsheet given sheetId and name of sheet
func AddSheet(spreadsheetId string, name string) error {
	srv := getService()
	request := &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: name,
			},
		},
	}
	batchUpdateRequest := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{request},
	}
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetId, &batchUpdateRequest).Do()
	// the _ is the response which takes this form https://godoc.org/google.golang.org/api/sheets/v4#BatchUpdateSpreadsheetResponse
	return err
}

// Add strings to a sheet, given sheet name, spreadsheetId, and slice of strings to add
func AddSheetRow(sheetName string, spreadsheetId string, values []string) error {
	srv := getService()
	// wrap name in single quotes to account for spaces
	writeRange := "'" + sheetName + "'"

	var vr sheets.ValueRange

	newRow := make([]interface{}, len(values))
	for i, v := range values {
		newRow[i] = v
	}
	vr.Values = append(vr.Values, newRow)

	valueInputOption := "RAW"
	insertDataOption := "INSERT_ROWS"

	_, err := srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, &vr).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve data from sheet. %v", err)
	// }
	return err
}
