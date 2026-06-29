package delivery

import (
	"net/http"
	"strconv"

	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type RatingHandler struct {
	ratingUsecase domain.RatingUsecase
}

func NewRatingHandler(ru domain.RatingUsecase) *RatingHandler {
	return &RatingHandler{ratingUsecase: ru}
}

func (h *RatingHandler) AddRating(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	type RatingRequest struct {
		BookingID      uint   `json:"booking_id" binding:"required"`
		RatingKualitas int    `json:"rating_kualitas" binding:"required,min=1,max=5"`
		RatingHarga    int    `json:"rating_harga" binding:"required,min=1,max=5"`
		Ulasan         string `json:"ulasan"`
	}

	var req RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating := &domain.Rating{
		PelangganID:    uid,
		BookingID:      req.BookingID,
		RatingKualitas: req.RatingKualitas,
		RatingHarga:    req.RatingHarga,
		Ulasan:         req.Ulasan,
	}

	err := h.ratingUsecase.AddRating(c.Request.Context(), rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": rating})
}

func (h *RatingHandler) GetRatings(c *gin.Context) {
	bengkelIDStr := c.Param("bengkel_id")
	bengkelID, err := strconv.ParseUint(bengkelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bengkel ID"})
		return
	}

	list, err := h.ratingUsecase.GetRatingsByBengkel(c.Request.Context(), uint(bengkelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": list})
}
