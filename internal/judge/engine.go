package judge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"

	"online-judge/internal/model"
	"online-judge/internal/sandbox"
)

type Engine struct {
	logger      *zap.Logger
	workBaseDir string
	compileTout time.Duration
}

func NewEngine(logger *zap.Logger, workBaseDir string, compileTout int) *Engine {
	return &Engine{
		logger:      logger,
		workBaseDir: workBaseDir,
		compileTout: time.Duration(compileTout) * time.Second,
	}
}

func (e *Engine) Judge(
	ctx context.Context,
	submission *model.Submission,
	problem *model.Problem,
	testCases []model.TestCase,
) (*model.JudgeResult, error) {
	runner, ok := LangRegistry[string(submission.Language)]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", submission.Language)
	}

	workDir, err := e.createWorkDir(submission.ID)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workDir)

	srcFile := filepath.Join(workDir, runner.SourceFileName())
	if err := os.WriteFile(srcFile, []byte(submission.Code), 0644); err != nil {
		return nil, fmt.Errorf("write source file: %w", err)
	}

	result := &model.JudgeResult{SubmissionID: submission.ID}

	if runner.NeedsCompile() {
		compileCtx, cancel := context.WithTimeout(ctx, e.compileTout)
		defer cancel()
		compileOutput, compileErr := runner.Compile(compileCtx, workDir, srcFile)
		result.CompileOutput = compileOutput
		if compileErr != nil {
			result.FinalStatus = model.StatusCompileError
			return result, nil
		}
	}

	var points []model.TestPointResult
	maxTime, maxMem := int64(0), int64(0)
	allAccepted := true

	for _, tc := range testCases {
		inputData, err := e.readTestData(tc.InputData, tc.InputFile)
		if err != nil {
			return nil, fmt.Errorf("read input for testcase %d: %w", tc.Seq, err)
		}
		expectedOutput, err := e.readTestData(tc.OutputData, tc.OutputFile)
		if err != nil {
			return nil, fmt.Errorf("read output for testcase %d: %w", tc.Seq, err)
		}

		executable, args := runner.ExecutableArgs(workDir)
		execCfg := &sandbox.ExecConfig{
			Executable:    executable,
			Args:          args,
			WorkDir:       workDir,
			TimeLimit:     time.Duration(problem.TimeLimit) * time.Millisecond,
			WallTime:      time.Duration(problem.TimeLimit*3) * time.Millisecond,
			MemoryLimit:   problem.MemoryLimit * 1024,
			Stdin:         inputData,
			SeccompPolicy: runner.SeccompPolicy(),
		}

		execResult, execErr := sandbox.Run(ctx, execCfg)
		point := e.evaluatePoint(tc, execCfg, execResult, execErr, expectedOutput)
		points = append(points, point)

		if point.TimeUsed > maxTime {
			maxTime = point.TimeUsed
		}
		if point.MemoryUsed > maxMem {
			maxMem = point.MemoryUsed
		}
		if point.Status != model.StatusAccepted {
			allAccepted = false
		}

		e.logger.Debug("test point judged",
			zap.Uint("submission_id", submission.ID),
			zap.Int("seq", tc.Seq),
			zap.String("status", string(point.Status)),
			zap.Int64("time_ms", point.TimeUsed),
			zap.Int64("memory_kb", point.MemoryUsed),
		)
	}

	result.Points = points
	result.TimeUsed = maxTime
	result.MemoryUsed = maxMem / 1024

	if allAccepted {
		result.FinalStatus = model.StatusAccepted
		for _, p := range points {
			result.TotalScore += p.Score
		}
	} else {
		for _, p := range points {
			if p.Status != model.StatusAccepted {
				result.FinalStatus = p.Status
				break
			}
		}
	}

	return result, nil
}

func (e *Engine) evaluatePoint(
	tc model.TestCase,
	cfg *sandbox.ExecConfig,
	execResult *sandbox.ExecResult,
	execErr error,
	expectedOutput []byte,
) model.TestPointResult {
	point := model.TestPointResult{Seq: tc.Seq, Score: 0}

	if execErr != nil {
		point.Status = model.StatusSystemError
		point.Message = execErr.Error()
		return point
	}

	if execResult.OOMKilled {
		point.TimeUsed = execResult.TimeUsed.Milliseconds()
		point.MemoryUsed = execResult.MemoryUsed / 1024
		point.Status = model.StatusMemLimitExc
		return point
	}

	point.TimeUsed = execResult.TimeUsed.Milliseconds()
	point.MemoryUsed = execResult.MemoryUsed / 1024

	if execResult.Killed || execResult.TimeUsed >= cfg.TimeLimit {
		point.Status = model.StatusTimeLimitExc
		return point
	}

	if execResult.ExitCode != 0 {
		point.Status = model.StatusRuntimeError
		point.Message = truncate(string(execResult.Stderr), 512)
		return point
	}

	if compareOutput(execResult.Stdout, expectedOutput) {
		point.Status = model.StatusAccepted
		point.Score = tc.Score
	} else {
		point.Status = model.StatusWrongAnswer
	}
	return point
}

func (e *Engine) readTestData(inline *string, filePath string) ([]byte, error) {
	if inline != nil {
		return []byte(*inline), nil
	}
	return os.ReadFile(filePath)
}

func (e *Engine) createWorkDir(submissionID uint) (string, error) {
	dir := filepath.Join(e.workBaseDir, fmt.Sprintf("sub_%d_%d", submissionID, time.Now().UnixNano()))
	return dir, os.MkdirAll(dir, 0750)
}

func PointsToJSON(points []model.TestPointResult) string {
	b, _ := json.Marshal(points)
	return string(b)
}
