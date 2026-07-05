package delivery

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

func (h *LayananHandler) cleanUpOrphanedImages(oldCSV, newCSV string) {
	if oldCSV == "" {
		return
	}
	oldURLs := strings.Split(oldCSV, ",")
	newURLsMap := make(map[string]bool)
	if newCSV != "" {
		for _, url := range strings.Split(newCSV, ",") {
			newURLsMap[url] = true
		}
	}

	for _, oldURL := range oldURLs {
		if oldURL != "" && !newURLsMap[oldURL] {
			if strings.Contains(oldURL, "/uploads/") {
				parts := strings.Split(oldURL, "/uploads/")
				if len(parts) > 1 {
					filename := parts[len(parts)-1]
					localPath := filepath.Join("uploads", filename)
					_ = os.Remove(localPath)
				}
			}
		}
	}
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

	// Clean up any removed photos from disk
	h.cleanUpOrphanedImages(existing.FotoURL, req.FotoURL)

	existing.Nama = req.Nama
	existing.Deskripsi = req.Deskripsi
	existing.EstimasiHarga = req.EstimasiHarga
	existing.Status = req.Status
	existing.FotoURL = req.FotoURL

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

	bengkel, err := h.bengkelUsecase.GetBengkelByPengelola(c.Request.Context(), uid)
	if err != nil || bengkel.ID != existing.BengkelID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this service's bengkel"})
		return
	}

	// Delete all photos from server disk when deleting service
	h.cleanUpOrphanedImages(existing.FotoURL, "")

	err = h.layananUsecase.DeleteLayanan(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Service deleted successfully"})
}
