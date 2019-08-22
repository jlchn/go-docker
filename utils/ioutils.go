package utils

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	return read, write, nil
}

func WriteToPipe(writePipe *os.File, content string) {
	writePipe.WriteString(content)
	writePipe.Close()
}

func ReadFromPipe() string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		panic(err)
	}
	msgStr := string(msg)
	log.Info(msgStr)
	return msgStr
}
