package utils

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()

		fields := strings.Split(txt, " ")

		if fields[9] != "cgroup" {
			continue
		}

		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	subSystemSCgroupPath := path.Join(FindCgroupMountpoint(subsystem), cgroupPath)
	if _, err := os.Stat(subSystemSCgroupPath); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(subSystemSCgroupPath, 0755); err == nil {
				log.Infof("cgroup: %s created", subSystemSCgroupPath)
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		log.Infof("cgroup: %s created", subSystemSCgroupPath)
		return subSystemSCgroupPath, nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
