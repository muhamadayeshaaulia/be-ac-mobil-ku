package domain

import "context"

type Layanan struct {
	ID            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	BengkelID     uint    `gorm:"index;not null" json:"bengkel_id"`
	Nama          string  `gorm:"not null" json:"nama"`
	Deskripsi     string  `gorm:"type:text" json:"deskripsi"`
	EstimasiHarga float64 `gorm:"not null" json:"estimasi_harga"`
	Status        string  `json:"status"` // "tersedia" atau "tidak_tersedia"
	FotoURL       string  `gorm:"type:text" json:"foto_url"`
}

type LayananRepository interface {
	GetByBengkelID(ctx context.Context, bengkelID uint) ([]Layanan, error)
	GetByID(ctx context.Context, id uint) (*Layanan, error)
	Create(ctx context.Context, layanan *Layanan) error
	Update(ctx context.Context, layanan *Layanan) error
	Delete(ctx context.Context, id uint) error
}

type LayananUsecase interface {
	GetLayananByBengkel(ctx context.Context, bengkelID uint) ([]Layanan, error)
	GetLayananByID(ctx context.Context, id uint) (*Layanan, error)
	CreateLayanan(ctx context.Context, layanan *Layanan) error
	UpdateLayanan(ctx context.Context, layanan *Layanan) error
	DeleteLayanan(ctx context.Context, id uint) error
}
