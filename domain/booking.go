package domain

import (
	"context"
	"time"
)

type Booking struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PelangganID    string    `gorm:"index;not null" json:"pelanggan_id"`
	BengkelID      uint      `gorm:"index;not null" json:"bengkel_id"`
	LayananID      uint      `gorm:"index;not null" json:"layanan_id"`
	TanggalBooking time.Time `gorm:"not null" json:"tanggal_booking"`
	Status         string    `gorm:"default:'menunggu'" json:"status"` // "menunggu", "dikonfirmasi", "selesai", "dibatalkan"
	Catatan        string    `gorm:"type:text" json:"catatan"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	Bengkel  *Bengkel `gorm:"foreignKey:BengkelID" json:"bengkel,omitempty"`
	Layanan  *Layanan `gorm:"foreignKey:LayananID" json:"layanan,omitempty"`
	Pelanggan *User    `gorm:"foreignKey:PelangganID;references:UID" json:"pelanggan,omitempty"`
}

type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	GetByID(ctx context.Context, id uint) (*Booking, error)
	GetByPelangganID(ctx context.Context, pelangganID string) ([]Booking, error)
	GetByBengkelID(ctx context.Context, bengkelID uint) ([]Booking, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, booking *Booking) error
	GetBookingByID(ctx context.Context, id uint) (*Booking, error)
	GetBookingHistory(ctx context.Context, pelangganID string) ([]Booking, error)
	GetBengkelQueue(ctx context.Context, pengelolaID string) ([]Booking, error)
	UpdateBookingStatus(ctx context.Context, pengelolaID string, bookingID uint, status string) error
}
