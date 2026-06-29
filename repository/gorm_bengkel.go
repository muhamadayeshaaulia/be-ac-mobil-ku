package repository

import (
	"context"

	"be-ac-mobil-ku/domain"
	"gorm.io/gorm"
)

type gormBengkelRepository struct {
	db *gorm.DB
}

func NewGormBengkelRepository(db *gorm.DB) domain.BengkelRepository {
	return &gormBengkelRepository{db: db}
}

func (r *gormBengkelRepository) GetAll(ctx context.Context) ([]domain.Bengkel, error) {
	var list []domain.Bengkel
	err := r.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (r *gormBengkelRepository) GetByID(ctx context.Context, id uint) (*domain.Bengkel, error) {
	var b domain.Bengkel
	err := r.db.WithContext(ctx).First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *gormBengkelRepository) GetByPengelolaID(ctx context.Context, pengelolaID string) (*domain.Bengkel, error) {
	var b domain.Bengkel
	err := r.db.WithContext(ctx).Where("pengelola_id = ?", pengelolaID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *gormBengkelRepository) Create(ctx context.Context, b *domain.Bengkel) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *gormBengkelRepository) Update(ctx context.Context, b *domain.Bengkel) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *gormBengkelRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Bengkel{}, id).Error
}
