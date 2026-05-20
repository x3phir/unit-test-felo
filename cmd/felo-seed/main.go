package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/felo/felo-backend/internal/demo/app"
	"github.com/felo/felo-backend/internal/demo/config"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}

	seedRoot := filepath.Join("tests", "e2e", "testdata", "seeds")
	if len(os.Args) > 1 {
		seedRoot = os.Args[1]
	}

	if err := application.Seed(ctx, seedRoot); err != nil {
		log.Fatalf("seed data: %v", err)
	}

	log.Printf("seed completed from %s", seedRoot)
}
