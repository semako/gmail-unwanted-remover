package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"net/http"
	"os"
)

const credentialsFilePath = "credentials.json"
const tokenFilePath = "token.json"

type GAPI struct {
	client *http.Client
}

func New() (*GAPI, error) {
	g := &GAPI{}

	cred, err := g.readCredentials()
	if err != nil {
		return nil, fmt.Errorf("gapi: failed to read credentials: %w", err)
	}

	config, err := google.ConfigFromJSON(cred, gmail.MailGoogleComScope)
	if err != nil {
		return nil, fmt.Errorf("gapi: unable to parse client secret file to config: %w", err)
	}

	cli, err := g.getClient(config)
	if err != nil {
		return nil, fmt.Errorf("gapi: failed to init client: %w", err)
	}

	g.client = cli

	return g, nil
}

func (g *GAPI) GetClient() *http.Client {
	return g.client
}

func (g *GAPI) readCredentials() ([]byte, error) {
	return os.ReadFile(credentialsFilePath)
}

// Retrieve a token, saves the token, then returns the generated client.
func (g *GAPI) getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := g.tokenFromFile(tokenFilePath)
	if err != nil {
		tok, err = g.getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err = g.saveToken(tokenFilePath, tok); err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func (g *GAPI) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %w", err)
	}

	return tok, nil
}

// Retrieves a token from a local file.
func (g *GAPI) tokenFromFile(file string) (*oauth2.Token, error) {
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
func (g *GAPI) saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}
	defer f.Close()
	_ = json.NewEncoder(f).Encode(token)
	return nil
}
