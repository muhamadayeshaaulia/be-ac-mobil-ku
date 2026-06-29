package domain

import "context"

type User struct {
	UID       string  `gorm:"primaryKey" json:"uid"`
	Email     string  `gorm:"uniqueIndex;not null" json:"email"`
	Nama      string  `gorm:"not null" json:"nama"`
	Role      string  `gorm:"not null" json:"role"` // "pelanggan" atau "pengelola_bengkel"
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Telepon   string  `json:"telepon"`
	FotoURL   string  `json:"foto_url"`
}

type UserRepository interface {
	GetByUID(ctx context.Context, uid string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type UserUsecase interface {
	GetProfile(ctx context.Context, uid string) (*User, error)
	RegisterOrUpdate(ctx context.Context, user *User) error
	UpdateLocation(ctx context.Context, uid string, lat, lng float64) error
}
