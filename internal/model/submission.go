package model

import "time"

type Language string

const (
	LangC      Language = "c"
	LangCPP    Language = "cpp"
	LangJava   Language = "java"
	LangPython Language = "python"
	LangGo     Language = "go"
)

type SubmissionStatus string

const (
	StatusPending      SubmissionStatus = "pending"
	StatusCompiling    SubmissionStatus = "compiling"
	StatusRunning      SubmissionStatus = "running"
	StatusAccepted     SubmissionStatus = "accepted"
	StatusWrongAnswer  SubmissionStatus = "wrong_answer"
	StatusTimeLimitExc SubmissionStatus = "time_limit_exceeded"
	StatusMemLimitExc  SubmissionStatus = "memory_limit_exceeded"
	StatusRuntimeError SubmissionStatus = "runtime_error"
	StatusCompileError SubmissionStatus = "compile_error"
	StatusSystemError  SubmissionStatus = "system_error"
)

type Submission struct {
	ID            uint             `gorm:"primarykey"                json:"id"`
	CreatedAt     time.Time        `                                 json:"created_at"`
	UpdatedAt     time.Time        `                                 json:"updated_at"`
	ProblemID     uint             `gorm:"not null;index"            json:"problem_id"`
	UserID        uint             `gorm:"not null;index"            json:"user_id"`
	Language      Language         `gorm:"size:20;not null"          json:"language"`
	Code          string           `gorm:"type:text;not null"        json:"code"`
	Status        SubmissionStatus `gorm:"size:50;default:'pending'" json:"status"`
	Score         int              `gorm:"default:0"                 json:"score"`
	CompileOutput string           `gorm:"type:text"                 json:"compile_output,omitempty"`
	JudgeDetails  string           `gorm:"type:text;default:'[]'"   json:"judge_details,omitempty"`
	TimeUsed      int64            `gorm:"default:0"                 json:"time_used"`
	MemoryUsed    int64            `gorm:"default:0"                 json:"memory_used"`
}
