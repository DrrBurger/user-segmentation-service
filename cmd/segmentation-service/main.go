package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"user-segmentation-service/config"
	"user-segmentation-service/internal/db"
	"user-segmentation-service/internal/server"
)

func main() {
	// Путь до файла конфигурации
	configPath := "config/config.yml"

	// Инициализация конфигурации
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err) // Завершение программы, если не удается загрузить конфигурацию
	}

	// Подключение к базе данных
	sqlDB, err := connectToDB(cfg)
	if err != nil {
		log.Fatal(err) // Завершение программы, если не удается подключиться к БД
	}

	myDB := db.NewDB(sqlDB)

	// Запуск приложения
	srv := server.NewApp(myDB).Run(cfg)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")
}

// Функция для подключения к базе данных
func connectToDB(cfg *config.Config) (*sql.DB, error) {
	var err error

	sqlDb, err := sql.Open("postgres", cfg.PG.URL) // для запуска локально использовать cfg.PG.URLLocal
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удается создать соединение
	}

	// Проверка соединения с базой данных
	if err := sqlDb.Ping(); err != nil {
		return nil, err // Возвращаем ошибку, если не удается установить соединение
	}

	return sqlDb, nil // Возвращаем подключение к базе данных
}
