package container

import (
	"go-docker/utils"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func Init() {
	setContainerName()
	mountProc()
}

func setContainerName() {

	if err := syscall.Sethostname([]byte("container")); err != nil {
		panic(err)
	}
}

func mountProc() {

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		panic(err)
	}
}

func readCommand() []string {

	return strings.Split(utils.ReadFromPipe(), " ")
}

func RunCommand() {

	cmdArray := readCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		panic("no command provided to start a container")
	}

	binary, err := exec.LookPath(cmdArray[0])
	if err != nil {
		panic(err)
	}
	if err := syscall.Exec(binary, cmdArray[0:], os.Environ()); err != nil {
		log.Error(err)
	}

}
