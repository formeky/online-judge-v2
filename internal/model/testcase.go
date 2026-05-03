package model

import "time"

type TestCase struct {
	ID         uint    `gorm:"primarykey"         json:"id"`
	CreatedAt  time.Time `                        json:"created_at"`
	ProblemID  uint    `gorm:"not null;index"     json:"problem_id"`
	Seq        int     `gorm:"not null;default:1" json:"seq"`
	InputFile  string  `gorm:"size:512"           json:"input_file"`
	OutputFile string  `gorm:"size:512"           json:"output_file"`
	InputData  *string `gorm:"type:text"          json:"-"`
	OutputData *string `gorm:"type:text"          json:"-"`
	Score      int     `gorm:"default:10"         json:"score"`
}
