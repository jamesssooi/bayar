package bayar

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func loadGoogleClientConfig() (*oauth2.Config, error) {
	config, _ := LoadConfig()
	path := filepath.Join(config.ApplicationDirectory, config.GoogleConfigurationFilename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &oauth2.Config{}, err
	}
	return google.ConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
}

func processAuthorizationCode(code string) (*oauth2.Token, error) {
	config, _ := LoadConfig()
	googlecfg, _ := loadGoogleClientConfig()

	token, err := googlecfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return &oauth2.Token{}, err
	}

	cachepath := filepath.Join(config.ApplicationDirectory, "token.json")
	file, err := os.OpenFile(cachepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return &oauth2.Token{}, err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		return &oauth2.Token{}, err
	}

	return token, nil
}
