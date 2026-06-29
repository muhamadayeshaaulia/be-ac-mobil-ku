package domain

import "context"

type Bengkel struct {
	ID                   uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	PengelolaID          string  `gorm:"index;not null" json:"pengelola_id"`
	Nama                 string  `gorm:"not null" json:"nama"`
	Alamat               string  `gorm:"type:text;not null" json:"alamat"`
	Latitude             float64 `gorm:"not null" json:"latitude"`
	Longitude            float64 `gorm:"not null" json:"longitude"`
	Deskripsi            string  `gorm:"type:text" json:"deskripsi"`
	JamBuka              string  `json:"jam_buka"`  // e.g. "08:00"
	JamTutup             string  `json:"jam_tutup"` // e.g. "17:00"
	Telepon              string  `json:"telepon"`
	Status               string  `json:"status"` // "aktif" atau "nonaktif"
	Distance             float64 `gorm:"-" json:"distance,omitempty"`
	AvgRatingKualitas    float64 `gorm:"-" json:"avg_rating_kualitas,omitempty"`
	AvgRatingHarga       float64 `gorm:"-" json:"avg_rating_harga,omitempty"`
	AvgRatingKeseluruhan float64 `gorm:"-" json:"avg_rating_keseluruhan,omitempty"`
}

type BengkelRepository interface {
	GetAll(ctx context.Context) ([]Bengkel, error)
	GetByID(ctx context.Context, id uint) (*Bengkel, error)
	GetByPengelolaID(ctx context.Context, pengelolaID string) (*Bengkel, error)
	Create(ctx context.Context, bengkel *Bengkel) error
	Update(ctx context.Context, bengkel *Bengkel) error
	Delete(ctx context.Context, id uint) error
}

type BengkelUsecase interface {
	ListBengkel(ctx context.Context) ([]Bengkel, error)
	GetBengkelByID(ctx context.Context, id uint) (*Bengkel, error)
	GetBengkelByPengelola(ctx context.Context, pengelolaID string) (*Bengkel, error)
	CreateBengkel(ctx context.Context, bengkel *Bengkel) error
	UpdateBengkel(ctx context.Context, bengkel *Bengkel) error
}
