package main

import (
	"gmail-unwanted-remover/internal/config"
	"gmail-unwanted-remover/internal/gapi"
	"log"
)

func main() {
	cfg := initCfg()
	ga := gapi.NewSimple(cfg.CredentialsFilePath, cfg.TokenFilePath)
	if err := ga.GenerateToken(); err != nil {
		log.Fatalf("failed to generate token: %v\n", err)
	}
}

func initCfg() *config.Config {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	return cfg
}
