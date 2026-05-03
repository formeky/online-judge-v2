package sandbox

import "time"

type PolicyType string

const (
	PolicyCCPP   PolicyType = "c_cpp"
	PolicyJava   PolicyType = "java"
	PolicyPython PolicyType = "python"
	PolicyGo     PolicyType = "golang"
)

type ExecConfig struct {
	Executable    string
	Args          []string
	Env           []string
	WorkDir       string
	TimeLimit     time.Duration
	WallTime      time.Duration
	MemoryLimit   int64
	Stdin         []byte
	SeccompPolicy PolicyType
}

type ExecResult struct {
	ExitCode   int
	Signal     int
	TimeUsed   time.Duration
	MemoryUsed int64
	Stdout     []byte
	Stderr     []byte
	Killed     bool
	OOMKilled  bool
}
