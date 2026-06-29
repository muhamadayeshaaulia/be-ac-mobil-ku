package usecase

import (
	"context"

	"be-ac-mobil-ku/domain"
)

type layananUsecase struct {
	layananRepo domain.LayananRepository
}

func NewLayananUsecase(layananRepo domain.LayananRepository) domain.LayananUsecase {
	return &layananUsecase{layananRepo: layananRepo}
}

func (u *layananUsecase) GetLayananByBengkel(ctx context.Context, bengkelID uint) ([]domain.Layanan, error) {
	return u.layananRepo.GetByBengkelID(ctx, bengkelID)
}

func (u *layananUsecase) GetLayananByID(ctx context.Context, id uint) (*domain.Layanan, error) {
	return u.layananRepo.GetByID(ctx, id)
}

func (u *layananUsecase) CreateLayanan(ctx context.Context, l *domain.Layanan) error {
	if l.Status == "" {
		l.Status = "tersedia"
	}
	return u.layananRepo.Create(ctx, l)
}

func (u *layananUsecase) UpdateLayanan(ctx context.Context, l *domain.Layanan) error {
	return u.layananRepo.Update(ctx, l)
}

func (u *layananUsecase) DeleteLayanan(ctx context.Context, id uint) error {
	return u.layananRepo.Delete(ctx, id)
}
