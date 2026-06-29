package usecase

import (
	"context"
	"errors"

	"be-ac-mobil-ku/domain"
)

type bookingUsecase struct {
	bookingRepo domain.BookingRepository
	bengkelRepo domain.BengkelRepository
}

func NewBookingUsecase(bookingRepo domain.BookingRepository, bengkelRepo domain.BengkelRepository) domain.BookingUsecase {
	return &bookingUsecase{
		bookingRepo: bookingRepo,
		bengkelRepo: bengkelRepo,
	}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, b *domain.Booking) error {
	b.Status = "menunggu"
	return u.bookingRepo.Create(ctx, b)
}

func (u *bookingUsecase) GetBookingByID(ctx context.Context, id uint) (*domain.Booking, error) {
	return u.bookingRepo.GetByID(ctx, id)
}

func (u *bookingUsecase) GetBookingHistory(ctx context.Context, pelangganID string) ([]domain.Booking, error) {
	return u.bookingRepo.GetByPelangganID(ctx, pelangganID)
}

func (u *bookingUsecase) GetBengkelQueue(ctx context.Context, pengelolaID string) ([]domain.Booking, error) {
	bengkel, err := u.bengkelRepo.GetByPengelolaID(ctx, pengelolaID)
	if err != nil {
		return nil, err
	}
	return u.bookingRepo.GetByBengkelID(ctx, bengkel.ID)
}

func (u *bookingUsecase) UpdateBookingStatus(ctx context.Context, pengelolaID string, bookingID uint, status string) error {
	booking, err := u.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	bengkel, err := u.bengkelRepo.GetByPengelolaID(ctx, pengelolaID)
	if err != nil {
		return err
	}

	if booking.BengkelID != bengkel.ID {
		return errors.New("unauthorized: booking does not belong to your bengkel")
	}

	return u.bookingRepo.UpdateStatus(ctx, bookingID, status)
}
