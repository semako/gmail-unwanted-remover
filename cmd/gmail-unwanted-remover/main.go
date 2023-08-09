package main

import (
	"context"
	"gmail-unwanted-remover/internal/config"
	"gmail-unwanted-remover/internal/gapi"
	"gmail-unwanted-remover/internal/gmail"
	"gmail-unwanted-remover/internal/matcher"
	"gmail-unwanted-remover/internal/remover"
	"log"
	"net/http"
)

func main() {
	cfg := initCfg()
	ga := initGAPI()
	mtchr := matcher.New(cfg.Matcher)
	ctx := context.Background()
	gm := initGmail(ctx, ga.GetClient())
	rm := remover.New(cfg.CheckInterval, gm, mtchr)

	if err := rm.Daemon(); err != nil {
		log.Fatalf("daemon failed with error: %v", err)
	}
}

func initGmail(ctx context.Context, client *http.Client) *gmail.Gmail {
	g, err := gmail.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to init gmail: %v", err)
	}
	return g
}

func initCfg() *config.Config {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	return cfg
}

func initGAPI() *gapi.GAPI {
	g, err := gapi.New()
	if err != nil {
		log.Fatalf("failed to init gapi: %v", err)
	}
	return g
}
