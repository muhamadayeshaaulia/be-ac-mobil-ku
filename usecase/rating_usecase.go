package usecase

import (
	"context"
	"errors"

	"be-ac-mobil-ku/domain"
)

type ratingUsecase struct {
	ratingRepo  domain.RatingRepository
	bookingRepo domain.BookingRepository
}

func NewRatingUsecase(ratingRepo domain.RatingRepository, bookingRepo domain.BookingRepository) domain.RatingUsecase {
	return &ratingUsecase{
		ratingRepo:  ratingRepo,
		bookingRepo: bookingRepo,
	}
}

func (u *ratingUsecase) AddRating(ctx context.Context, rating *domain.Rating) error {
	// Validate booking exists and is completed
	booking, err := u.bookingRepo.GetByID(ctx, rating.BookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	if booking.PelangganID != rating.PelangganID {
		return errors.New("unauthorized: this booking is not yours to rate")
	}

	if booking.Status != "selesai" {
		return errors.New("booking must be completed (selesai) before rating")
	}

	rating.BengkelID = booking.BengkelID
	return u.ratingRepo.Create(ctx, rating)
}

func (u *ratingUsecase) GetRatingsByBengkel(ctx context.Context, bengkelID uint) ([]domain.Rating, error) {
	return u.ratingRepo.GetByBengkelID(ctx, bengkelID)
}
