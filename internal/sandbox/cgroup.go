//go:build linux

package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type CgroupManager struct {
	name   string
	cgPath string
}

func NewCgroupManager(name string) (*CgroupManager, error) {
	cgBase := "/sys/fs/cgroup/oj"
	cgPath := filepath.Join(cgBase, name)
	if err := os.MkdirAll(cgPath, 0755); err != nil {
		return nil, fmt.Errorf("create cgroup dir: %w", err)
	}
	return &CgroupManager{name: name, cgPath: cgPath}, nil
}

func (c *CgroupManager) SetMemoryLimit(bytes int64) error {
	return writeFile(filepath.Join(c.cgPath, "memory.max"), strconv.FormatInt(bytes, 10))
}

func (c *CgroupManager) SetMemorySwapLimit() error {
	return writeFile(filepath.Join(c.cgPath, "memory.swap.max"), "0")
}

func (c *CgroupManager) AddProcess(pid int) error {
	return writeFile(filepath.Join(c.cgPath, "cgroup.procs"), strconv.Itoa(pid))
}

func (c *CgroupManager) GetMemoryUsed() (int64, error) {
	data, err := os.ReadFile(filepath.Join(c.cgPath, "memory.current"))
	if err != nil {
		return 0, err
	}
	s := string(data)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return strconv.ParseInt(s, 10, 64)
}

func (c *CgroupManager) IsOOMKilled() (bool, error) {
	data, err := os.ReadFile(filepath.Join(c.cgPath, "memory.events"))
	if err != nil {
		return false, err
	}
	var oomKill int
	for _, line := range splitLines(string(data)) {
		var key string
		var val int
		if _, err := fmt.Sscanf(line, "%s %d", &key, &val); err == nil && key == "oom_kill" {
			oomKill = val
		}
	}
	return oomKill > 0, nil
}

func (c *CgroupManager) Cleanup() error {
	return os.Remove(c.cgPath)
}

func writeFile(path, value string) error {
	return os.WriteFile(path, []byte(value), 0644)
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
