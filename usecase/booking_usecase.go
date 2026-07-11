package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

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

func generateOrderNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	dateStr := time.Now().Format("20060102")
	randomNum := r.Intn(9000) + 1000 // 1000-9999
	return fmt.Sprintf("ORD-TRX-%s-%d", dateStr, randomNum)
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, b *domain.Booking) error {
	count, err := u.bookingRepo.CountActiveByBengkelAndTime(ctx, b.BengkelID, b.TanggalBooking)
	if err != nil {
		return err
	}
	if count >= 5 {
		return errors.New("bengkel penuh pada jadwal tersebut (maksimal 5 pelanggan)")
	}

	b.Status = "menunggu"
	b.OrderNumber = generateOrderNumber()
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

func (u *bookingUsecase) GetFullTimeSlots(ctx context.Context, bengkelID uint, dateStr string) ([]string, error) {
	// dateStr is expected to be "YYYY-MM-DD"
	loc, _ := time.LoadLocation("Local") // or time.Now().Location()
	start, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return nil, err
	}
	end := start.Add(24 * time.Hour)

	bookings, err := u.bookingRepo.GetActiveBookingsByBengkelAndDate(ctx, bengkelID, start, end)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, b := range bookings {
		timeStr := b.TanggalBooking.Format("15:04")
		counts[timeStr]++
	}

	var fullSlots []string
	for t, c := range counts {
		if c >= 5 {
			fullSlots = append(fullSlots, t)
		}
	}
	return fullSlots, nil
}
