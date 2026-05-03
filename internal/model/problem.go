package model

import (
	"time"

	"gorm.io/gorm"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type Problem struct {
	ID          uint           `gorm:"primarykey"                          json:"id"`
	CreatedAt   time.Time      `                                           json:"created_at"`
	UpdatedAt   time.Time      `                                           json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                               json:"-"`
	Title       string         `gorm:"size:255;not null;uniqueIndex"       json:"title"`
	Description string         `gorm:"type:text;not null"                  json:"description"`
	Difficulty  Difficulty     `gorm:"size:20;default:'easy'"              json:"difficulty"`
	TimeLimit   int64          `gorm:"not null;default:2000"               json:"time_limit"`
	MemoryLimit int64          `gorm:"not null;default:262144"             json:"memory_limit"`
	AllowedLangs string        `gorm:"size:100;default:'c,cpp,java,python,go'" json:"allowed_langs"`
	TestCases   []TestCase     `gorm:"foreignKey:ProblemID"               json:"test_cases,omitempty"`
}
