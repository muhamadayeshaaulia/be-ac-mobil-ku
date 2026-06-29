package delivery

import (
	"net/http"
	"strconv"

	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type LayananHandler struct {
	layananUsecase domain.LayananUsecase
	bengkelUsecase domain.BengkelUsecase
}

func NewLayananHandler(lu domain.LayananUsecase, bu domain.BengkelUsecase) *LayananHandler {
	return &LayananHandler{
		layananUsecase: lu,
		bengkelUsecase: bu,
	}
}

func (h *LayananHandler) GetLayananByBengkel(c *gin.Context) {
	bengkelIDStr := c.Param("bengkel_id")
	bengkelID, err := strconv.ParseUint(bengkelIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bengkel ID"})
		return
	}

	list, err := h.layananUsecase.GetLayananByBengkel(c.Request.Context(), uint(bengkelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": list})
}

func (h *LayananHandler) GetLayananByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	l, err := h.layananUsecase.GetLayananByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": l})
}

func (h *LayananHandler) CreateLayanan(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	var req domain.Layanan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the bengkel belongs to this pengelola
	bengkel, err := h.bengkelUsecase.GetBengkelByPengelola(c.Request.Context(), uid)
	if err != nil || bengkel.ID != req.BengkelID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this bengkel"})
		return
	}

	err = h.layananUsecase.CreateLayanan(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": req})
}

func (h *LayananHandler) UpdateLayanan(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	existing, err := h.layananUsecase.GetLayananByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	// Verify the bengkel belongs to this pengelola
	bengkel, err := h.bengkelUsecase.GetBengkelByPengelola(c.Request.Context(), uid)
	if err != nil || bengkel.ID != existing.BengkelID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this service's bengkel"})
		return
	}

	var req domain.Layanan
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing.Nama = req.Nama
	existing.Deskripsi = req.Deskripsi
	existing.EstimasiHarga = req.EstimasiHarga
	existing.Status = req.Status

	err = h.layananUsecase.UpdateLayanan(c.Request.Context(), existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": existing})
}

func (h *LayananHandler) DeleteLayanan(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	existing, err := h.layananUsecase.GetLayananByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	// Verify the bengkel belongs to this pengelola
	bengkel, err := h.bengkelUsecase.GetBengkelByPengelola(c.Request.Context(), uid)
	if err != nil || bengkel.ID != existing.BengkelID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this service's bengkel"})
		return
	}

	err = h.layananUsecase.DeleteLayanan(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Service deleted successfully"})
}
