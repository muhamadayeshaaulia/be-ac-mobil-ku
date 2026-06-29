package usecase

import (
	"context"
	"errors"

	"be-ac-mobil-ku/domain"
	"gorm.io/gorm"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) GetProfile(ctx context.Context, uid string) (*domain.User, error) {
	return u.userRepo.GetByUID(ctx, uid)
}

func (u *userUsecase) RegisterOrUpdate(ctx context.Context, user *domain.User) error {
	existing, err := u.userRepo.GetByUID(ctx, user.UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Register
			if user.Role == "" {
				user.Role = "pelanggan" // default role
			}
			return u.userRepo.Create(ctx, user)
		}
		return err
	}

	// Update existing record
	existing.Nama = user.Nama
	existing.Email = user.Email
	if user.Role != "" {
		existing.Role = user.Role
	}
	if user.FotoURL != "" {
		existing.FotoURL = user.FotoURL
	}
	if user.Telepon != "" {
		existing.Telepon = user.Telepon
	}
	return u.userRepo.Update(ctx, existing)
}

func (u *userUsecase) UpdateLocation(ctx context.Context, uid string, lat, lng float64) error {
	existing, err := u.userRepo.GetByUID(ctx, uid)
	if err != nil {
		return err
	}
	existing.Latitude = lat
	existing.Longitude = lng
	return u.userRepo.Update(ctx, existing)
}
