package storage

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"backup-keeper/internal/domain"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveStorage struct {
	service  *drive.Service
	folderID string
}

func NewGoogleDriveStorage(credentialsFile, tokenFile string) (domain.Storage, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client, err := getClient(config, tokenFile)
	if err != nil {
		return nil, fmt.Errorf("unable to create OAuth client: %v", err)
	}

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive service: %v", err)
	}

	return &GoogleDriveStorage{service: srv}, nil
}

func (s *GoogleDriveStorage) Save(fileName string, data []byte) error {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(data); err != nil {
		return fmt.Errorf("compression failed: %v", err)
	}

	if err := gz.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %v", err)
	}

	compressedName := fileName + ".gz"
	return s.Save(compressedName, buf.Bytes())
}

// Retrieve a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read token file: %v", err)
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	var authCode string
	fmt.Print("Enter authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		panic(fmt.Sprintf("Unable to read authorization code: %v", err))
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		panic(fmt.Sprintf("Unable to retrieve token from web: %v", err))
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

// Returns an HTTP client that automatically refreshes tokens as needed
func getClient(config *oauth2.Config, tokenFile string) (*http.Client, error) {
	// Try to load cached token
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		// Token not found, get new one
		tok = getTokenFromWeb(config)
		if err := saveToken(tokenFile, tok); err != nil {
			return nil, fmt.Errorf("failed to save token: %v", err)
		}
	}

	// Check if token needs refresh
	if !tok.Valid() {
		newTok, err := config.TokenSource(context.Background(), tok).Token()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}

		if newTok.AccessToken != tok.AccessToken {
			if err := saveToken(tokenFile, newTok); err != nil {
				return nil, fmt.Errorf("failed to save refreshed token: %v", err)
			}
		}
		tok = newTok
	}

	return config.Client(context.Background(), tok), nil
}
