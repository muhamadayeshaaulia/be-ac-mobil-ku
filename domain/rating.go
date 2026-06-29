package domain

import (
	"context"
	"time"
)

type Rating struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PelangganID    string    `gorm:"index;not null" json:"pelanggan_id"`
	BengkelID      uint      `gorm:"index;not null" json:"bengkel_id"`
	BookingID      uint      `gorm:"uniqueIndex;not null" json:"booking_id"`
	RatingKualitas int       `gorm:"not null" json:"rating_kualitas"` // 1-5
	RatingHarga    int       `gorm:"not null" json:"rating_harga"`    // 1-5
	Ulasan         string    `gorm:"type:text" json:"ulasan"`
	CreatedAt      time.Time `json:"created_at"`

	// Relation
	Pelanggan *User    `gorm:"foreignKey:PelangganID;references:UID" json:"pelanggan,omitempty"`
	Bengkel   *Bengkel `gorm:"foreignKey:BengkelID" json:"bengkel,omitempty"`
}

type RatingRepository interface {
	Create(ctx context.Context, rating *Rating) error
	GetByBengkelID(ctx context.Context, bengkelID uint) ([]Rating, error)
	GetAll(ctx context.Context) ([]Rating, error)
	GetAverageRatingByBengkelID(ctx context.Context, bengkelID uint) (avgKualitas, avgHarga float64, err error)
}

type RatingUsecase interface {
	AddRating(ctx context.Context, rating *Rating) error
	GetRatingsByBengkel(ctx context.Context, bengkelID uint) ([]Rating, error)
}
