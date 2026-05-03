package judge

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"online-judge/internal/model"
	"online-judge/internal/mq"
	"online-judge/internal/repository"
)

type Worker struct {
	engine      *Engine
	subRepo     repository.SubmissionRepository
	problemRepo repository.ProblemRepository
	tcRepo      repository.TestCaseRepository
	logger      *zap.Logger
}

func NewWorker(
	engine *Engine,
	subRepo repository.SubmissionRepository,
	problemRepo repository.ProblemRepository,
	tcRepo repository.TestCaseRepository,
	logger *zap.Logger,
) *Worker {
	return &Worker{
		engine:      engine,
		subRepo:     subRepo,
		problemRepo: problemRepo,
		tcRepo:      tcRepo,
		logger:      logger,
	}
}

func (w *Worker) Handle(ctx context.Context, msg *mq.JudgeMessage) error {
	log := w.logger.With(zap.Uint("submission_id", msg.SubmissionID))
	log.Info("start judging")

	sub, err := w.subRepo.GetByID(ctx, msg.SubmissionID)
	if err != nil {
		return err
	}
	if sub.Status != model.StatusPending {
		log.Info("skip: already processed", zap.String("status", string(sub.Status)))
		return nil
	}

	if err := w.subRepo.UpdateStatus(ctx, msg.SubmissionID, model.StatusCompiling); err != nil {
		return err
	}

	problem, err := w.problemRepo.GetByID(ctx, msg.ProblemID)
	if err != nil {
		return err
	}

	testCases, err := w.tcRepo.GetByProblemID(ctx, msg.ProblemID)
	if err != nil {
		return err
	}

	_ = w.subRepo.UpdateStatus(ctx, msg.SubmissionID, model.StatusRunning)

	judgeResult, err := w.engine.Judge(ctx, sub, problem, testCases)
	if err != nil {
		log.Error("judge system error", zap.Error(err))
		_ = w.subRepo.UpdateStatus(ctx, msg.SubmissionID, model.StatusSystemError)
		return err
	}

	detailsJSON, _ := json.Marshal(judgeResult.Points)
	updates := map[string]interface{}{
		"status":         judgeResult.FinalStatus,
		"score":          judgeResult.TotalScore,
		"compile_output": judgeResult.CompileOutput,
		"judge_details":  string(detailsJSON),
		"time_used":      judgeResult.TimeUsed,
		"memory_used":    judgeResult.MemoryUsed,
	}
	if err := w.subRepo.UpdateResult(ctx, msg.SubmissionID, updates); err != nil {
		return err
	}

	log.Info("judge completed",
		zap.String("status", string(judgeResult.FinalStatus)),
		zap.Int("score", judgeResult.TotalScore),
	)
	return nil
}
