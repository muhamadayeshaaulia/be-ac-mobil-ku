package usecase

import (
	"context"

	"be-ac-mobil-ku/domain"
)

type bengkelUsecase struct {
	bengkelRepo domain.BengkelRepository
	ratingRepo  domain.RatingRepository
}

func NewBengkelUsecase(bengkelRepo domain.BengkelRepository, ratingRepo domain.RatingRepository) domain.BengkelUsecase {
	return &bengkelUsecase{
		bengkelRepo: bengkelRepo,
		ratingRepo:  ratingRepo,
	}
}

func (u *bengkelUsecase) ListBengkel(ctx context.Context) ([]domain.Bengkel, error) {
	list, err := u.bengkelRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	// Populate average ratings
	for i := range list {
		avgKualitas, avgHarga, _ := u.ratingRepo.GetAverageRatingByBengkelID(ctx, list[i].ID)
		list[i].AvgRatingKualitas = avgKualitas
		list[i].AvgRatingHarga = avgHarga
		if avgKualitas > 0 || avgHarga > 0 {
			list[i].AvgRatingKeseluruhan = (avgKualitas + avgHarga) / 2.0
		}
	}
	return list, nil
}

func (u *bengkelUsecase) GetBengkelByID(ctx context.Context, id uint) (*domain.Bengkel, error) {
	b, err := u.bengkelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	avgKualitas, avgHarga, _ := u.ratingRepo.GetAverageRatingByBengkelID(ctx, b.ID)
	b.AvgRatingKualitas = avgKualitas
	b.AvgRatingHarga = avgHarga
	if avgKualitas > 0 || avgHarga > 0 {
		b.AvgRatingKeseluruhan = (avgKualitas + avgHarga) / 2.0
	}
	return b, nil
}

func (u *bengkelUsecase) GetBengkelByPengelola(ctx context.Context, pengelolaID string) (*domain.Bengkel, error) {
	b, err := u.bengkelRepo.GetByPengelolaID(ctx, pengelolaID)
	if err != nil {
		return nil, err
	}
	avgKualitas, avgHarga, _ := u.ratingRepo.GetAverageRatingByBengkelID(ctx, b.ID)
	b.AvgRatingKualitas = avgKualitas
	b.AvgRatingHarga = avgHarga
	if avgKualitas > 0 || avgHarga > 0 {
		b.AvgRatingKeseluruhan = (avgKualitas + avgHarga) / 2.0
	}
	return b, nil
}

func (u *bengkelUsecase) CreateBengkel(ctx context.Context, b *domain.Bengkel) error {
	if b.Status == "" {
		b.Status = "aktif"
	}
	return u.bengkelRepo.Create(ctx, b)
}

func (u *bengkelUsecase) UpdateBengkel(ctx context.Context, b *domain.Bengkel) error {
	return u.bengkelRepo.Update(ctx, b)
}
