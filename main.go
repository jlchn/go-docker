package main

import (
	"go-docker/command"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.Name = "go-docker"
	app.Usage = "go-docker run [-it] <command>"

	app.Before = func(context *cli.Context) error {
		return nil
	}

	app.Commands = []cli.Command{
		command.RunCommand,
		command.RunCommandInContainerCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
