package app

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"user-segmentation-service/internal/models"
)

// createSegmentHandler создает сегмент и добавляет в него установленный % случайных пользователей
func (a *App) createSegmentHandler(ctx *gin.Context) {
	var segment models.Segment

	// Парсинг JSON-запроса в структуру "Segment"
	if err := ctx.BindJSON(&segment); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Проверка допустимости значения поля "RandomPercentage"
	if segment.RandomPercentage < 0 || segment.RandomPercentage > 100 {
		respondWithError(ctx, http.StatusBadRequest, "RandomPercentage should be between 0 and 100")
		return
	}

	err := a.db.CreateSegment(segment.Slug, segment.RandomPercentage, segment.ExpirationDate)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Segment and user assignments created successfully"})
}

// deleteSegmentHandler обрабатывает удаление сегмента
func (a *App) deleteSegmentHandler(ctx *gin.Context) {
	var segment models.Segment

	// Десериализация JSON-запроса в структуру Segment
	if err := ctx.BindJSON(&segment); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	segmentID, err := a.db.DeleteSegment(segment.Slug)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Segment deleted successfully", "segment_id": segmentID})
}
