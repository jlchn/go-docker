package command

import (
	"fmt"
	"go-docker/container"

	"github.com/urfave/cli"
)

/**
go-docker run [-it] <command>
*/
var RunCommand = cli.Command{
	Name: "run",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "it",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("command missing")
		}
		tty := context.Bool("it")
		container.StartContainer(tty, context.Args())
		return nil
	},
}

var RunCommandInContainerCommand = cli.Command{
	Name: "runCommandInContainer",
	Action: func(context *cli.Context) error {

		container.Init()
		container.RunCommand()
		return nil
	},
}
