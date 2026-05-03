//go:build linux

package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Run(ctx context.Context, cfg *ExecConfig) (*ExecResult, error) {
	cgName := fmt.Sprintf("oj_%d", time.Now().UnixNano())
	cg, err := NewCgroupManager(cgName)
	if err != nil {
		return nil, fmt.Errorf("create cgroup: %w", err)
	}
	defer cg.Cleanup()

	if err := cg.SetMemoryLimit(cfg.MemoryLimit); err != nil {
		return nil, err
	}
	if err := cg.SetMemorySwapLimit(); err != nil {
		return nil, err
	}

	cmd := buildSandboxedCmd(cfg.Executable, cfg.Args)
	cmd.Dir = cfg.WorkDir
	if len(cfg.Env) > 0 {
		cmd.Env = cfg.Env
	} else {
		cmd.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(cfg.Stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	wallTime := cfg.WallTime
	if wallTime == 0 {
		wallTime = cfg.TimeLimit * 3
	}
	wallCtx, cancel := context.WithTimeout(ctx, wallTime)
	defer cancel()

	startTime := time.Now()

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start process: %w", err)
	}

	if err := cg.AddProcess(cmd.Process.Pid); err != nil {
		_ = cmd.Process.Kill()
		return nil, fmt.Errorf("add to cgroup: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	result := &ExecResult{}

	select {
	case <-wallCtx.Done():
		_ = cmd.Process.Kill()
		<-done
		result.Killed = true
		result.TimeUsed = cfg.TimeLimit + time.Millisecond
	case err := <-done:
		result.TimeUsed = time.Since(startTime)
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				result.ExitCode = exitErr.ExitCode()
			}
		}
	}

	memBytes, _ := cg.GetMemoryUsed()
	result.MemoryUsed = memBytes

	oomKilled, _ := cg.IsOOMKilled()
	result.OOMKilled = oomKilled

	result.Stdout = stdout.Bytes()
	result.Stderr = stderr.Bytes()

	if result.OOMKilled {
		if p, err := os.FindProcess(cmd.Process.Pid); err == nil {
			_ = p.Kill()
		}
	}

	return result, nil
}
