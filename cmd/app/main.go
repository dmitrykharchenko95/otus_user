package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dmitrykharchenko95/otus_user/config"
	"github.com/dmitrykharchenko95/otus_user/internal/database"
	"github.com/dmitrykharchenko95/otus_user/internal/server"
	"github.com/dmitrykharchenko95/otus_user/internal/server/handlers"
)

func main() {
	var (
		cfg     = config.NewFromENVs()
		db, err = database.New(cfg.DB)
	)

	if err != nil {
		log.Fatal(err)
	}

	var migrate string
	flag.StringVar(&migrate, "m", "", "migrate up/down")
	flag.Parse()

	if migrate != "" {
		if err = doMigrate(db, migrate); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	var svc = server.NewServer(cfg.Server, handlers.NewHandler(db, cfg.JWTKey))

	if err = svc.Start(); err != nil {
		log.Fatal(err)
	}
}

func doMigrate(db database.Manager, cmd string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch cmd {
	case "up":
		return db.MigrateUp(ctx)
	case "down":
		return db.MigrateDown(ctx)
	default:
		return fmt.Errorf("unknown migrate command: %s", cmd)
	}
}
