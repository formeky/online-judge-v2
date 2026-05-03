package model

type TestPointResult struct {
	Seq        int              `json:"seq"`
	Status     SubmissionStatus `json:"status"`
	TimeUsed   int64            `json:"time_used"`
	MemoryUsed int64            `json:"memory_used"`
	Score      int              `json:"score"`
	Message    string           `json:"message,omitempty"`
}

type JudgeResult struct {
	SubmissionID  uint
	FinalStatus   SubmissionStatus
	TotalScore    int
	TimeUsed      int64
	MemoryUsed    int64
	CompileOutput string
	Points        []TestPointResult
}
