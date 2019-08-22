package cgroup

import (
	"go-docker/cgroup/subsystem"

	log "github.com/Sirupsen/logrus"
)

type CgroupManager struct {
	Path     string
	Resource *subsystem.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystem.SubsystemsIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystem.ResourceConfig) error {
	for _, subSysIns := range subsystem.SubsystemsIns {
		subSysIns.CreateCgroup(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystem.SubsystemsIns {
		if err := subSysIns.RemoveCgroup(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
