package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"user-segmentation-service/config"
	"user-segmentation-service/internal/app"
	"user-segmentation-service/internal/db"
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
	err = app.NewApp(myDB).Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// Функция для подключения к базе данных
func connectToDB(cfg *config.Config) (*sql.DB, error) {
	var err error

	sqlDb, err := sql.Open("postgres", cfg.PG.URLLocal) // для запуска в docker использовать cfg.PG.URL
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удается создать соединение
	}

	// Проверка соединения с базой данных
	if err := sqlDb.Ping(); err != nil {
		return nil, err // Возвращаем ошибку, если не удается установить соединение
	}

	return sqlDb, nil // Возвращаем подключение к базе данных
}
