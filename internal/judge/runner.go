package judge

import (
	"context"

	"online-judge/internal/sandbox"
)

type LanguageRunner interface {
	Name() string
	NeedsCompile() bool
	SourceFileName() string
	SeccompPolicy() sandbox.PolicyType
	Compile(ctx context.Context, workDir, srcFile string) (output string, err error)
	ExecutableArgs(workDir string) (executable string, args []string)
}

var LangRegistry = map[string]LanguageRunner{}

func RegisterRunner(r LanguageRunner) {
	LangRegistry[r.Name()] = r
}
