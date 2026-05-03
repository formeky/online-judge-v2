package repository

import (
	"context"

	"gorm.io/gorm"

	"online-judge/internal/model"
)

type SubmissionRepository interface {
	Create(ctx context.Context, s *model.Submission) error
	GetByID(ctx context.Context, id uint) (*model.Submission, error)
	List(ctx context.Context, problemID, userID uint, page, size int) ([]model.Submission, int64, error)
	UpdateStatus(ctx context.Context, id uint, status model.SubmissionStatus) error
	UpdateResult(ctx context.Context, id uint, updates map[string]interface{}) error
}

type submissionRepo struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepo{db: db}
}

func (r *submissionRepo) Create(ctx context.Context, s *model.Submission) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *submissionRepo) GetByID(ctx context.Context, id uint) (*model.Submission, error) {
	var s model.Submission
	err := r.db.WithContext(ctx).First(&s, id).Error
	return &s, err
}

func (r *submissionRepo) List(ctx context.Context, problemID, userID uint, page, size int) ([]model.Submission, int64, error) {
	var submissions []model.Submission
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Submission{})
	if problemID > 0 {
		q = q.Where("problem_id = ?", problemID)
	}
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	q.Count(&total)
	offset := (page - 1) * size
	err := q.Order("id desc").Offset(offset).Limit(size).Find(&submissions).Error
	return submissions, total, err
}

func (r *submissionRepo) UpdateStatus(ctx context.Context, id uint, status model.SubmissionStatus) error {
	return r.db.WithContext(ctx).Model(&model.Submission{}).Where("id = ?", id).Update("status", status).Error
}

func (r *submissionRepo) UpdateResult(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Submission{}).Where("id = ?", id).Updates(updates).Error
}
