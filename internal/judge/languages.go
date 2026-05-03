package judge

import (
	"context"
	"os/exec"
	"path/filepath"

	"online-judge/internal/sandbox"
)

type CRunner struct{}

func (r *CRunner) Name() string                        { return "c" }
func (r *CRunner) NeedsCompile() bool                  { return true }
func (r *CRunner) SourceFileName() string              { return "solution.c" }
func (r *CRunner) SeccompPolicy() sandbox.PolicyType   { return sandbox.PolicyCCPP }

func (r *CRunner) Compile(ctx context.Context, workDir, srcFile string) (string, error) {
	outFile := filepath.Join(workDir, "solution")
	cmd := exec.CommandContext(ctx, "gcc", "-O2", "-std=c11", "-lm", "-o", outFile, srcFile)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (r *CRunner) ExecutableArgs(workDir string) (string, []string) {
	return filepath.Join(workDir, "solution"), []string{}
}

type CppRunner struct{}

func (r *CppRunner) Name() string                        { return "cpp" }
func (r *CppRunner) NeedsCompile() bool                  { return true }
func (r *CppRunner) SourceFileName() string              { return "solution.cpp" }
func (r *CppRunner) SeccompPolicy() sandbox.PolicyType   { return sandbox.PolicyCCPP }

func (r *CppRunner) Compile(ctx context.Context, workDir, srcFile string) (string, error) {
	outFile := filepath.Join(workDir, "solution")
	cmd := exec.CommandContext(ctx, "g++", "-O2", "-std=c++17", "-lm", "-o", outFile, srcFile)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (r *CppRunner) ExecutableArgs(workDir string) (string, []string) {
	return filepath.Join(workDir, "solution"), []string{}
}

type JavaRunner struct{}

func (r *JavaRunner) Name() string                        { return "java" }
func (r *JavaRunner) NeedsCompile() bool                  { return true }
func (r *JavaRunner) SourceFileName() string              { return "Main.java" }
func (r *JavaRunner) SeccompPolicy() sandbox.PolicyType   { return sandbox.PolicyJava }

func (r *JavaRunner) Compile(ctx context.Context, workDir, srcFile string) (string, error) {
	cmd := exec.CommandContext(ctx, "javac", "-encoding", "UTF-8", srcFile)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (r *JavaRunner) ExecutableArgs(workDir string) (string, []string) {
	return "java", []string{
		"-cp", workDir,
		"-Xss8m",
		"-Dfile.encoding=UTF-8",
		"-XX:+UseSerialGC",
		"-XX:MaxRAMPercentage=75.0",
		"Main",
	}
}

type PythonRunner struct{}

func (r *PythonRunner) Name() string                        { return "python" }
func (r *PythonRunner) NeedsCompile() bool                  { return false }
func (r *PythonRunner) SourceFileName() string              { return "solution.py" }
func (r *PythonRunner) SeccompPolicy() sandbox.PolicyType   { return sandbox.PolicyPython }

func (r *PythonRunner) Compile(_ context.Context, _, _ string) (string, error) {
	return "", nil
}

func (r *PythonRunner) ExecutableArgs(workDir string) (string, []string) {
	return "python3", []string{"-u", filepath.Join(workDir, "solution.py")}
}

type GoRunner struct{}

func (r *GoRunner) Name() string                        { return "go" }
func (r *GoRunner) NeedsCompile() bool                  { return true }
func (r *GoRunner) SourceFileName() string              { return "solution.go" }
func (r *GoRunner) SeccompPolicy() sandbox.PolicyType   { return sandbox.PolicyGo }

func (r *GoRunner) Compile(ctx context.Context, workDir, srcFile string) (string, error) {
	outFile := filepath.Join(workDir, "solution")
	cmd := exec.CommandContext(ctx, "go", "build", "-o", outFile, srcFile)
	cmd.Dir = workDir
	cmd.Env = append(cmd.Environ(), "GOPATH="+filepath.Join(workDir, "gopath"), "GOCACHE="+filepath.Join(workDir, "gocache"))
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (r *GoRunner) ExecutableArgs(workDir string) (string, []string) {
	return filepath.Join(workDir, "solution"), []string{}
}

func init() {
	RegisterRunner(&CRunner{})
	RegisterRunner(&CppRunner{})
	RegisterRunner(&JavaRunner{})
	RegisterRunner(&PythonRunner{})
	RegisterRunner(&GoRunner{})
}
