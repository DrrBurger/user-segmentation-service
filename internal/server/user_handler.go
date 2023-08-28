package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"

	"user-segmentation-service/internal/models"
)

// createUserHandler создает нового пользователя.
func (a *App) createUserHandler(ctx *gin.Context) {
	var user models.User

	// Пробуем привязать JSON к структуре User.
	if err := ctx.BindJSON(&user); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := a.db.CreateUser(user.Name)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": userID})
}

// deleteUserHandler удаляет пользователя по ID, полученному из JSON.
func (a *App) deleteUserHandler(ctx *gin.Context) {
	var req models.DeleteUserRequest

	// Привязываем входящий JSON к структуре DeleteUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	userID, err := a.db.DeleteUser(req.UserId)
	if err != nil {
		respondWithError(ctx, http.StatusNotFound, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "user_id": userID})
}

// updateUserSegmentsHandler обновляет сегменты пользователя.
func (a *App) updateUserSegmentsHandler(ctx *gin.Context) {
	var req models.UpdateSegmentsRequest

	// Привязываем входящий JSON к структуре UpdateSegmentsRequest.
	if err := ctx.BindJSON(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := a.db.UpdateUserSegments(req.UserId, req.Add, req.Remove)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User segments updated successfully", "user_id": userID})
}

// getUserSegmentsHandler возвращает сегменты пользователя.
func (a *App) getUserSegmentsHandler(ctx *gin.Context) {
	var req models.UserSegmentsRequest

	// Привязываем входящий JSON к структуре UserSegmentsRequest.
	if err := ctx.BindJSON(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userID, segments, err := a.db.GetUserSegments(req.UserId)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID, "segments": segments})
}

// getUserReportHandler создает CSV отчет по истории сегментов пользователя.
func (a *App) getUserReportHandler(ctx *gin.Context) {
	var req models.ReportRequest

	// Привязываем входящий JSON к структуре ReportRequest.
	if err := ctx.BindJSON(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	fileName, err := a.db.GetUserReport(req.UserId, req.YearMonth)
	if err != nil {
		respondWithError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reportHost := os.Getenv("HTTP_REPORT_HOST")

	ctx.JSON(http.StatusOK, gin.H{"message": "Report generated successfully", "download_link": reportHost + fileName})
}
