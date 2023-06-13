package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sequelie"
)

func main() {
	app := &cli.App{
		Name:        "sequelie",
		Description: "Off-loading your SQL queries from your Golang code, now includes codegen!",
		Commands: []*cli.Command{
			{
				Name:        "generate",
				Aliases:     []string{"g"},
				Description: "Generates Golang files containing SQL queries from Sequelie.",
				Action: func(context *cli.Context) error {
					dir := context.String("directory")
					if dir == "" {
						dir = "./"
					}
					if err := sequelie.ReadDirectory(dir); err != nil {
						return err
					}
					if err := sequelie.Generate(); err != nil {
						return err
					}
					log.Println("[SQL] Generated all Golang files under the \".sequelie\" directory.")
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "directory",
						Aliases: []string{"d"},
						Usage:   "Reads and generates all the Sequelie files from the specific directory (recursive).",
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
