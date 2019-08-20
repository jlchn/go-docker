package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

/**
go-docker run [-it] <command>
*/
var runCommand = cli.Command{
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
		SetNamespace(tty, context.Args())

		return nil
	},
}

var startContainerCommand = cli.Command{
	Name: "startContainer",
	Action: func(context *cli.Context) error {

		//set a new hostname
		if err := syscall.Sethostname([]byte("container")); err != nil {
			fmt.Printf("Setting Hostname failed")
		}

		//mount proc folder
		defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
		if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
			panic(err)
		}
		log.Info(context.Args())
		cmd := context.Args().Get(0)

		binary, err := exec.LookPath(cmd)
		if err != nil {
			panic(err)
		}

		if err := syscall.Exec(binary, context.Args(), os.Environ()); err != nil {
			log.Error(err)
		}

		return nil
	},
}

func SetNamespace(tty bool, command []string) {
	args := append([]string{"startContainer"}, command...)
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}

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
		runCommand,
		startContainerCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
