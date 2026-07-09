package delivery

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(uu domain.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uu}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	profile, err := h.userUsecase.GetProfile(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": profile})
}

func (h *UserHandler) deleteOldUploadedFile(oldURL string) {
	if oldURL == "" {
		return
	}
	if strings.Contains(oldURL, "/uploads/") {
		parts := strings.Split(oldURL, "/uploads/")
		if len(parts) > 1 {
			filename := parts[len(parts)-1]
			localPath := filepath.Join("uploads", filename)
			_ = os.Remove(localPath)
		}
	}
}

func (h *UserHandler) RegisterOrUpdate(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)
	email := c.MustGet("user_email").(string)
	name := c.MustGet("user_name").(string)

	type RegisterRequest struct {
		Role      string  `json:"role" binding:"required"` // "pelanggan" or "pengelola_bengkel"
		Nama      string  `json:"nama"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Telepon   string  `json:"telepon"`
		FotoURL   string  `json:"foto_url"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Clean up old photo if updated
	if existingUser, err := h.userUsecase.GetProfile(c.Request.Context(), uid); err == nil && existingUser != nil {
		if existingUser.FotoURL != req.FotoURL {
			h.deleteOldUploadedFile(existingUser.FotoURL)
		}
	}

	displayName := name
	if req.Nama != "" {
		displayName = req.Nama
	}

	user := &domain.User{
		UID:       uid,
		Email:     email,
		Nama:      displayName,
		Role:      req.Role,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Telepon:   req.Telepon,
		FotoURL:   req.FotoURL,
	}

	err := h.userUsecase.RegisterOrUpdate(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedUser, _ := h.userUsecase.GetProfile(c.Request.Context(), uid)
	if updatedUser == nil {
		updatedUser = user
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Profile saved successfully", "data": updatedUser})
}

func (h *UserHandler) UpdateLocation(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	type LocationRequest struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	var req LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.UpdateLocation(c.Request.Context(), uid, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Location updated successfully"})
}
