package delivery

import (
	"net/http"

	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(r *gin.Engine, uu domain.UserUsecase, authMiddleware gin.HandlerFunc) {
	handler := &UserHandler{userUsecase: uu}

	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		api.GET("/user/profile", handler.GetProfile)
		api.POST("/user/profile", handler.RegisterOrUpdate)
		api.POST("/user/location", handler.UpdateLocation)
	}
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

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Profile saved successfully", "data": user})
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
