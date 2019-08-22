package container

import (
	"go-docker/utils"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func StartContainer(tty bool, command []string) {
	cmd, writePipe := prepareCommand(tty)
	folkContainerProcess(cmd, writePipe, command)
}

func prepareCommand(tty bool) (*exec.Cmd, *os.File) {
	cmd := exec.Command("/proc/self/exe", "runCommandInContainer")
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

	readPipe, writePipe, err := utils.NewPipe()

	if err != nil {
		panic("cannot create pipe for docker process and container process")
	}

	cmd.ExtraFiles = []*os.File{readPipe}

	return cmd, writePipe
}

func folkContainerProcess(cmd *exec.Cmd, writePipe *os.File, command []string) {
	utils.WriteToPipe(writePipe, strings.Join(command, " "))

	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}
