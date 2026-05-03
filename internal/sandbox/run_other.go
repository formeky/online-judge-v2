//go:build !linux

package sandbox

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

func Run(ctx context.Context, cfg *ExecConfig) (*ExecResult, error) {
	wallTime := cfg.WallTime
	if wallTime == 0 {
		wallTime = cfg.TimeLimit * 3
	}
	wallCtx, cancel := context.WithTimeout(ctx, wallTime)
	defer cancel()

	cmd := exec.CommandContext(wallCtx, cfg.Executable, cfg.Args...)
	cmd.Dir = cfg.WorkDir
	if len(cfg.Env) > 0 {
		cmd.Env = cfg.Env
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(cfg.Stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	result := &ExecResult{}

	err := cmd.Run()
	result.TimeUsed = time.Since(startTime)

	if err != nil {
		if wallCtx.Err() != nil {
			result.Killed = true
			result.TimeUsed = cfg.TimeLimit + time.Millisecond
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}

	result.Stdout = stdout.Bytes()
	result.Stderr = stderr.Bytes()
	return result, nil
}
