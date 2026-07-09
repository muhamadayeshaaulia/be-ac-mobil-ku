package repository

import (
	"context"

	"be-ac-mobil-ku/domain"
	"gorm.io/gorm"
)

type gormRatingRepository struct {
	db *gorm.DB
}

func NewGormRatingRepository(db *gorm.DB) domain.RatingRepository {
	return &gormRatingRepository{db: db}
}

func (r *gormRatingRepository) Create(ctx context.Context, rating *domain.Rating) error {
	return r.db.WithContext(ctx).Create(rating).Error
}

func (r *gormRatingRepository) GetByBengkelID(ctx context.Context, bengkelID uint) ([]domain.Rating, error) {
	var list []domain.Rating
	err := r.db.WithContext(ctx).
		Preload("Pelanggan").
		Preload("Booking").
		Preload("Booking.Layanan").
		Where("bengkel_id = ?", bengkelID).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *gormRatingRepository) GetAll(ctx context.Context) ([]domain.Rating, error) {
	var list []domain.Rating
	err := r.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (r *gormRatingRepository) GetAverageRatingByBengkelID(ctx context.Context, bengkelID uint) (avgKualitas, avgHarga float64, err error) {
	type Result struct {
		AvgKualitas float64
		AvgHarga    float64
	}
	var res Result
	err = r.db.WithContext(ctx).
		Model(&domain.Rating{}).
		Select("COALESCE(AVG(rating_kualitas), 0) as avg_kualitas, COALESCE(AVG(rating_harga), 0) as avg_harga").
		Where("bengkel_id = ?", bengkelID).
		Scan(&res).Error
	if err != nil {
		return 0, 0, err
	}
	return res.AvgKualitas, res.AvgHarga, nil
}
