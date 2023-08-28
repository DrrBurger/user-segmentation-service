package app

import (
	"github.com/gin-gonic/gin"
	"user-segmentation-service/config"
	"user-segmentation-service/internal/db"
)

// App структура для приложения
type App struct {
	db db.InterfaceDB
}

// NewApp создаёт новый экземпляр приложения
func NewApp(db db.InterfaceDB) *App {
	return &App{db: db}
}

// Run запускает приложение
func (a *App) Run(cfg *config.Config) error {
	r := a.setupRouter()        // Настройка маршрутизации
	err := r.Run(cfg.HTTP.Port) // Запуск сервера на определенном порту
	if err != nil {
		return err // Завершение программы, если не удаётся запустить сервер T
	}
	return nil
}

// setupRouter настраивает маршрутизацию для приложения
func (a *App) setupRouter() *gin.Engine {
	r := gin.Default()
	// Определение обработчиков маршрутов
	r.POST("/user", a.createUserHandler)
	r.DELETE("/user", a.deleteUserHandler)
	r.POST("/segment", a.createSegmentHandler)
	r.DELETE("/segment", a.deleteSegmentHandler)
	r.POST("/user/segments", a.updateUserSegmentsHandler)
	r.GET("/user/segments", a.getUserSegmentsHandler)
	r.GET("/user/report", a.getUserReportHandler)
	r.Static("/user/report", "reports")

	return r
}

// respondWithError отправляет ошибку клиенту
func respondWithError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"error": message}) // Отправка JSON ответа с кодом ошибки и сообщением
	ctx.Abort()                               // Прекращение обработки текущего запроса
}
