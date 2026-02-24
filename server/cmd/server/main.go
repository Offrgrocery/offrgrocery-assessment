package main

import (
	"context"
	"log"
	"os"

	"offgrocery-assessment/internal/app"
	"offgrocery-assessment/internal/config"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "server",
		Usage: "offrgrocery server",
		Commands: []*cli.Command{
			{
				Name:  "web",
				Usage: "start the web server",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg := config.Load()
					return app.NewWeb(cfg)
				},
			},
			{
				Name:  "import",
				Usage: "import the database with grocery data",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg := config.Load()
					return app.NewImporter(cfg, "internal/seed/data/store_a.json")
				},
			},
			{
				Name:  "seed",
				Usage: "seed the database with dev data",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg := config.Load()
					return app.NewSeed(cfg)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("main: %v", err)
	}
}
