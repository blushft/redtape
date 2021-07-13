package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "redtape"
	app.Usage = "redtape cli"
	app.Commands = []*cli.Command{
		roleBuildCmd(),
		policyCmd(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
