package delivery

import (
	"net/http"
	"strconv"

	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type BengkelHandler struct {
	bengkelUsecase domain.BengkelUsecase
}

func NewBengkelHandler(bu domain.BengkelUsecase) *BengkelHandler {
	return &BengkelHandler{bengkelUsecase: bu}
}

func (h *BengkelHandler) ListBengkel(c *gin.Context) {
	list, err := h.bengkelUsecase.ListBengkel(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": list})
}

func (h *BengkelHandler) GetBengkelByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bengkel ID"})
		return
	}

	b, err := h.bengkelUsecase.GetBengkelByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bengkel not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": b})
}

func (h *BengkelHandler) GetMyBengkel(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	b, err := h.bengkelUsecase.GetBengkelByPengelola(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bengkel not found for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": b})
}

func (h *BengkelHandler) CreateBengkel(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	var b domain.Bengkel
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b.PengelolaID = uid

	err := h.bengkelUsecase.CreateBengkel(c.Request.Context(), &b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": b})
}

func (h *BengkelHandler) UpdateBengkel(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	existing, err := h.bengkelUsecase.GetBengkelByPengelola(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bengkel not found"})
		return
	}

	var input domain.Bengkel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.Nama = input.Nama
	existing.Alamat = input.Alamat
	existing.Latitude = input.Latitude
	existing.Longitude = input.Longitude
	existing.Deskripsi = input.Deskripsi
	existing.JamBuka = input.JamBuka
	existing.JamTutup = input.JamTutup
	existing.Telepon = input.Telepon
	existing.Status = input.Status
	existing.FotoURL = input.FotoURL

	err = h.bengkelUsecase.UpdateBengkel(c.Request.Context(), existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": existing})
}
