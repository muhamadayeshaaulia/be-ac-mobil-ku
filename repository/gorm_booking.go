package repository

import (
	"context"
	"time"

	"be-ac-mobil-ku/domain"
	"gorm.io/gorm"
)

type gormBookingRepository struct {
	db *gorm.DB
}

func NewGormBookingRepository(db *gorm.DB) domain.BookingRepository {
	return &gormBookingRepository{db: db}
}

func (r *gormBookingRepository) Create(ctx context.Context, b *domain.Booking) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *gormBookingRepository) GetByID(ctx context.Context, id uint) (*domain.Booking, error) {
	var b domain.Booking
	err := r.db.WithContext(ctx).
		Preload("Bengkel").
		Preload("Layanan").
		Preload("Pelanggan").
		First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *gormBookingRepository) GetByPelangganID(ctx context.Context, pelangganID string) ([]domain.Booking, error) {
	var list []domain.Booking
	err := r.db.WithContext(ctx).
		Preload("Bengkel").
		Preload("Layanan").
		Preload("Pelanggan").
		Preload("Rating").
		Where("pelanggan_id = ?", pelangganID).
		Order("tanggal_booking desc").
		Find(&list).Error
	return list, err
}

func (r *gormBookingRepository) GetByBengkelID(ctx context.Context, bengkelID uint) ([]domain.Booking, error) {
	var list []domain.Booking
	err := r.db.WithContext(ctx).
		Preload("Layanan").
		Preload("Pelanggan").
		Where("bengkel_id = ?", bengkelID).
		Order("tanggal_booking desc").
		Find(&list).Error
	return list, err
}

func (r *gormBookingRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Booking{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *gormBookingRepository) CountActiveByBengkelAndTime(ctx context.Context, bengkelID uint, t time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Booking{}).
		Where("bengkel_id = ? AND tanggal_booking = ? AND status IN ('menunggu', 'dikonfirmasi')", bengkelID, t).
		Count(&count).Error
	return count, err
}
