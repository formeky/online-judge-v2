package mq

type JudgeMessage struct {
	SubmissionID uint   `json:"submission_id"`
	ProblemID    uint   `json:"problem_id"`
	Language     string `json:"language"`
	TimeLimit    int64  `json:"time_limit"`
	MemoryLimit  int64  `json:"memory_limit"`
}
