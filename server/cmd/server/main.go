package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
	"offgrocery-assessment/internal/app"
	"offgrocery-assessment/internal/config"
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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("main: %v", err)
	}
}
