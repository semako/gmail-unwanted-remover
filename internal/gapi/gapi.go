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
	client                             *http.Client
	credentialsFilePath, tokenFilePath string
}

func New(credentialsFilePath, tokenFilePath string) (*GAPI, error) {
	g := &GAPI{
		credentialsFilePath: credentialsFilePath,
		tokenFilePath:       tokenFilePath,
	}

	cfg, err := g.getConfig()
	if err != nil {
		return nil, fmt.Errorf("gapi: failed to init config: %w", err)
	}

	cli, err := g.getClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("gapi: failed to init client: %w", err)
	}

	g.client = cli

	return g, nil
}

func NewSimple(credentialsFilePath, tokenFilePath string) *GAPI {
	return &GAPI{
		credentialsFilePath: credentialsFilePath,
		tokenFilePath:       tokenFilePath,
	}
}

func (g *GAPI) GetClient() *http.Client {
	return g.client
}

func (g *GAPI) getConfig() (*oauth2.Config, error) {
	cred, err := g.readCredentials()
	if err != nil {
		return nil, fmt.Errorf("gapi: failed to read credentials: %w", err)
	}

	return google.ConfigFromJSON(cred, gmail.MailGoogleComScope)
}

func (g *GAPI) readCredentials() ([]byte, error) {
	return os.ReadFile(g.credentialsFilePath)
}

func (g *GAPI) GenerateToken() error {
	cfg, err := g.getConfig()
	if err != nil {
		return fmt.Errorf("gapi: GenerateToken: failed to init config: %w", err)
	}

	tok, err := g.getTokenFromWeb(cfg)
	if err != nil {
		return fmt.Errorf("gapi: GenerateToken: failed to get token from web: %w", err)
	}
	if err = g.saveToken(g.tokenFilePath, tok); err != nil {
		return fmt.Errorf("gapi: GenerateToken: failed to save token: %w", err)
	}
	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func (g *GAPI) getClient(config *oauth2.Config) (*http.Client, error) {
	tok, err := g.tokenFromFile(g.tokenFilePath)
	if err != nil {
		return nil, err
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
