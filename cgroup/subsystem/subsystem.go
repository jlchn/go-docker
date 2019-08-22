package subsystem

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	GetName() string
	CreateCgroup(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	RemoveCgroup(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&MemorySubSystem{},
	}
)
