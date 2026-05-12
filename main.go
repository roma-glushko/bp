package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/roma-glushko/bp/cmd"
	"github.com/roma-glushko/bp/internal/version"
)

const Copyright = `-Present, Roma Hlushko (c)`

func main() {
	cliApp := cli.App{
		Name:                 "bp",
		Usage:                "A blood pressure journal",
		Version:              version.FullVersion,
		Copyright:            Copyright,
		Suggest:              true,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable debug logging",
			},
		},
		Before: func(c *cli.Context) error {
			level := slog.LevelInfo
			if c.Bool("debug") {
				level = slog.LevelDebug
			}

			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: level,
			}))
			slog.SetDefault(logger)

			return nil
		},
		Commands: []*cli.Command{
			cmd.HelloCommand,
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}
