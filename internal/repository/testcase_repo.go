package repository

import (
	"context"

	"gorm.io/gorm"

	"online-judge/internal/model"
)

type TestCaseRepository interface {
	Create(ctx context.Context, tc *model.TestCase) error
	GetByProblemID(ctx context.Context, problemID uint) ([]model.TestCase, error)
	DeleteByProblemID(ctx context.Context, problemID uint) error
}

type testCaseRepo struct {
	db *gorm.DB
}

func NewTestCaseRepository(db *gorm.DB) TestCaseRepository {
	return &testCaseRepo{db: db}
}

func (r *testCaseRepo) Create(ctx context.Context, tc *model.TestCase) error {
	return r.db.WithContext(ctx).Create(tc).Error
}

func (r *testCaseRepo) GetByProblemID(ctx context.Context, problemID uint) ([]model.TestCase, error) {
	var tcs []model.TestCase
	err := r.db.WithContext(ctx).Where("problem_id = ?", problemID).Order("seq asc").Find(&tcs).Error
	return tcs, err
}

func (r *testCaseRepo) DeleteByProblemID(ctx context.Context, problemID uint) error {
	return r.db.WithContext(ctx).Where("problem_id = ?", problemID).Delete(&model.TestCase{}).Error
}
