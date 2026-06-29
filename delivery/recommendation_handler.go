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

func NewRecommendationHandler(ru domain.RecommendationUsecase) *RecommendationHandler {
	return &RecommendationHandler{recUsecase: ru}
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
