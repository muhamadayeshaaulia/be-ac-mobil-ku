package delivery

import (
	"net/http"
	"strconv"

	"be-ac-mobil-ku/domain"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	recUsecase domain.RecommendationUsecase
}

func NewRecommendationHandler(r *gin.Engine, ru domain.RecommendationUsecase, authMiddleware gin.HandlerFunc) {
	handler := &RecommendationHandler{recUsecase: ru}

	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		api.GET("/recommendations", handler.GetRecommendations)
	}
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	uid := c.MustGet("user_uid").(string)

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	recs, err := h.recUsecase.GetRecommendations(c.Request.Context(), uid, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": recs})
}
