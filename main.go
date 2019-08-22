package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

		cmdArray := ReadCommand()
		log.Info(cmdArray)
		if cmdArray == nil || len(cmdArray) == 0 {
			panic("no command provided to start a container")
		}

		//set a new hostname
		if err := syscall.Sethostname([]byte("container")); err != nil {
			panic(err)
		}

		//mount proc folder
		defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
		if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
			panic(err)
		}

		binary, err := exec.LookPath(cmdArray[0])
		if err != nil {
			panic(err)
		}
		log.Info(cmdArray[0])
		log.Info(cmdArray[0:])
		if err := syscall.Exec(binary, cmdArray[0:], os.Environ()); err != nil {
			log.Error(err)
		}

		return nil
	},
}

func SetNamespace(tty bool, command []string) {
	cmd := exec.Command("/proc/self/exe", "startContainer")
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

	readPipe, writePipe, err := NewPipe()
	if err != nil {
		panic("cannot create pipe for docker process and container process")
	}

	cmd.ExtraFiles = []*os.File{readPipe}
	writePipe.WriteString(strings.Join(command, " "))
	writePipe.Close()
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	return read, write, nil
}

func ReadCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	log.Info(msgStr)
	return strings.Split(msgStr, " ")
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
