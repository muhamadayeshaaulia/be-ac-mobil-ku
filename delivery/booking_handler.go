package delivery

import (
	"net/http"
	"strconv"
	"time"

	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUsecase domain.BookingUsecase
}

func NewBookingHandler(bu domain.BookingUsecase) *BookingHandler {
	return &BookingHandler{bookingUsecase: bu}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	type BookingRequest struct {
		BengkelID      uint   `json:"bengkel_id" binding:"required"`
		LayananID      uint   `json:"layanan_id" binding:"required"`
		TanggalBooking string `json:"tanggal_booking" binding:"required"` // ISO string format or "2006-01-02 15:04"
		Catatan        string `json:"catatan"`
	}

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Try standard formats
	var tBooking time.Time
	var err error
	tBooking, err = time.Parse(time.RFC3339, req.TanggalBooking)
	if err != nil {
		tBooking, err = time.Parse("2006-01-02 15:04:00", req.TanggalBooking)
		if err != nil {
			tBooking, err = time.Parse("2006-01-02 15:04", req.TanggalBooking)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 or 'YYYY-MM-DD HH:MM'"})
				return
			}
		}
	}

	booking := &domain.Booking{
		PelangganID:    uid,
		BengkelID:      req.BengkelID,
		LayananID:      req.LayananID,
		TanggalBooking: tBooking,
		Catatan:        req.Catatan,
	}

	err = h.bookingUsecase.CreateBooking(c.Request.Context(), booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": booking})
}

func (h *BookingHandler) GetHistory(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	history, err := h.bookingUsecase.GetBookingHistory(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": history})
}

func (h *BookingHandler) GetQueue(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	queue, err := h.bookingUsecase.GetBengkelQueue(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": queue})
}

func (h *BookingHandler) UpdateStatus(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	type StatusRequest struct {
		Status string `json:"status" binding:"required"` // "dikonfirmasi", "selesai", "dibatalkan"
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.bookingUsecase.UpdateBookingStatus(c.Request.Context(), uid, uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Booking status updated successfully"})
}
