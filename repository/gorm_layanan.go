package repository

import (
	"context"

	"be-ac-mobil-ku/domain"
	"gorm.io/gorm"
)

type gormLayananRepository struct {
	db *gorm.DB
}

func NewGormLayananRepository(db *gorm.DB) domain.LayananRepository {
	return &gormLayananRepository{db: db}
}

func (r *gormLayananRepository) GetByBengkelID(ctx context.Context, bengkelID uint) ([]domain.Layanan, error) {
	var list []domain.Layanan
	err := r.db.WithContext(ctx).Where("bengkel_id = ?", bengkelID).Find(&list).Error
	return list, err
}

func (r *gormLayananRepository) GetByID(ctx context.Context, id uint) (*domain.Layanan, error) {
	var l domain.Layanan
	err := r.db.WithContext(ctx).First(&l, id).Error
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *gormLayananRepository) Create(ctx context.Context, l *domain.Layanan) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *gormLayananRepository) Update(ctx context.Context, l *domain.Layanan) error {
	return r.db.WithContext(ctx).Save(l).Error
}

func (r *gormLayananRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Layanan{}, id).Error
}
