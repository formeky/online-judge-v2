package repository

import (
	"context"

	"gorm.io/gorm"

	"online-judge/internal/model"
)

type ProblemRepository interface {
	Create(ctx context.Context, p *model.Problem) error
	GetByID(ctx context.Context, id uint) (*model.Problem, error)
	List(ctx context.Context, page, size int) ([]model.Problem, int64, error)
	Update(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
}

type problemRepo struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) ProblemRepository {
	return &problemRepo{db: db}
}

func (r *problemRepo) Create(ctx context.Context, p *model.Problem) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *problemRepo) GetByID(ctx context.Context, id uint) (*model.Problem, error) {
	var p model.Problem
	err := r.db.WithContext(ctx).First(&p, id).Error
	return &p, err
}

func (r *problemRepo) List(ctx context.Context, page, size int) ([]model.Problem, int64, error) {
	var problems []model.Problem
	var total int64
	offset := (page - 1) * size
	r.db.WithContext(ctx).Model(&model.Problem{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(size).Find(&problems).Error
	return problems, total, err
}

func (r *problemRepo) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Problem{}).Where("id = ?", id).Updates(updates).Error
}

func (r *problemRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Problem{}, id).Error
}
