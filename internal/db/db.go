package db

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"user-segmentation-service/internal/models"
)

type DB struct {
	db *sql.DB
}

func NewDB(sqlDB *sql.DB) *DB {
	return &DB{db: sqlDB}
}

type InterfaceDB interface {
	CreateUser(name string) (int64, error)
	DeleteUser(userID int) (int, error)
	CreateSegment(slug string, randomPercentage float64, expirationDate time.Time) error
	DeleteSegment(slug string) (int, error)
	UpdateUserSegments(userID int, addList []models.Segment, removeList []string) (int, error)
	GetUserSegments(userID int) (int, []string, error)
	GetUserReport(userID int, yearMonth string) (string, error)
}

func (db *DB) CreateUser(name string) (int64, error) {
	var userID int64

	// Вставляем пользователя в базу данных и получаем его ID.
	err := db.db.QueryRow(
		"INSERT INTO users(name) VALUES($1) RETURNING id",
		name,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (db *DB) DeleteUser(userID int) (int, error) {
	// Начало транзакции
	tx, err := db.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("An error occurred while rolling back the transaction: %v\n", err)
		}
	}()

	// Проверка существования пользователя в базе данных
	var existingId int
	err = tx.QueryRow("SELECT id FROM users WHERE id = $1", userID).Scan(&existingId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user with ID %d does not exist", userID)
		}
		return 0, fmt.Errorf("failed to query user with ID %d: %w", userID, err)
	}

	// Удаление сегментов пользователя
	if _, err = tx.Exec(
		"DELETE FROM user_segments WHERE user_id=$1",
		userID,
	); err != nil {
		return 0, fmt.Errorf("failed to delete user_segments with ID %d: %w", userID, err)
	}

	// Удаление самого пользователя
	if _, err = tx.Exec(
		"DELETE FROM users WHERE id=$1",
		userID,
	); err != nil {
		return 0, fmt.Errorf("failed to delete user with ID %d: %w", userID, err)
	}

	// Подтверждение транзакции
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

func (db *DB) CreateSegment(slug string, randomPercentage float64, expirationDate time.Time) error {

	if expirationDate.IsZero() {
		return fmt.Errorf("expirationDate should not be zero")
	}

	currentTime := time.Now()
	if expirationDate.Before(currentTime.Add(1 * time.Hour)) {
		return fmt.Errorf("expirationDate should be at least 1 hours in the future")
	}

	// Начало транзакции
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("An error occurred while rolling back the transaction: %v\n", err)
		}
	}()

	// Проверка на существование сегмента с таким же slug
	var existingId int
	err = tx.QueryRow("SELECT id FROM segments WHERE slug = $1", slug).Scan(&existingId)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return fmt.Errorf("failed to query existing segment: %w", err)
		}

		return fmt.Errorf("segment with slug '%s' already exists", slug)
	}

	// Вставка нового сегмента
	_, err = tx.Exec("INSERT INTO segments(slug) VALUES($1)", slug)
	if err != nil {
		return fmt.Errorf("failed to insert new segment: %w", err)
	}

	// Получение общего числа пользователей
	var totalUsers int
	err = tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		return fmt.Errorf("failed to count total users: %w", err)
	}

	// Вычисление числа пользователей для добавления в сегмент
	numUsersToAdd := int(float64(totalUsers) * (randomPercentage / 100.0))

	// Создание временной таблицы
	_, err = tx.Exec("CREATE TEMP TABLE temp_users AS SELECT id FROM users ORDER BY RANDOM() LIMIT $1", numUsersToAdd)
	if err != nil {
		return fmt.Errorf("failed to create temp table: %w", err)
	}

	// Добавление пользователей в сегмент
	_, err = tx.Exec(
		`INSERT INTO user_segments(user_id, segment_slug, expiration_date)
         SELECT id, $1, $2 FROM temp_users`,
		slug, expirationDate,
	)
	if err != nil {
		return fmt.Errorf("failed to add users to segment: %w", err)
	}

	// Логирование операции добавления
	_, err = tx.Exec(
		`INSERT INTO user_segment_history(user_id, segment_slug, operation)
         SELECT id, $1, 'add' FROM temp_users`,
		slug,
	)
	if err != nil {
		return fmt.Errorf("failed to log segment addition: %w", err)
	}

	// Удаление временной таблицы
	_, err = tx.Exec("DROP TABLE temp_users")
	if err != nil {
		return fmt.Errorf("failed to drop temp table: %w", err)
	}

	// Подтверждение транзакции
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *DB) DeleteSegment(slug string) (int, error) {
	// Начало транзакции
	tx, err := db.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("An error occurred while rolling back the transaction: %v\n", err)
		}
	}()

	// Проверка наличия сегмента в базе данных
	var existingId int
	err = tx.QueryRow("SELECT id FROM segments WHERE slug = $1", slug).Scan(&existingId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("segment with slug '%s' does not exist", slug)
	} else if err != nil {
		return 0, fmt.Errorf("failed to query existing segment: %w", err)
	}

	// Удаление записей о сегменте из таблицы user_segments
	if _, err = tx.Exec("DELETE FROM user_segments WHERE segment_slug = $1", slug); err != nil {
		return 0, fmt.Errorf("failed to delete segment from user_segments: %w", err)
	}

	var segmentId int
	// Удаление сегмента из таблицы segments и возвращение его ID
	err = tx.QueryRow("DELETE FROM segments WHERE slug=$1 RETURNING id", slug).Scan(&segmentId)
	if err != nil {
		return 0, fmt.Errorf("failed to delete segment: %w", err)
	}

	// Подтверждение транзакции
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return segmentId, nil
}

func (db *DB) UpdateUserSegments(userID int, addList []models.Segment, removeList []string) (int, error) {
	// Начинаем транзакцию
	tx, err := db.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("An error occurred while rolling back the transaction: %v\n", err)
		}
	}()

	// Проверка существования пользователя
	var existingUserId int
	err = tx.QueryRow("SELECT id FROM users WHERE id = $1", userID).Scan(&existingUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("user with ID '%d' does not exist", userID)
	} else if err != nil {
		return 0, fmt.Errorf("failed to query existing user: %w", err)
	}

	// Добавляем сегменты
	for _, segment := range addList {
		var existingSlug string
		err = tx.QueryRow("SELECT slug FROM segments WHERE slug = $1", segment.Slug).Scan(&existingSlug)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("segment with slug '%s' does not exist", segment.Slug)
		} else if err != nil {
			return 0, fmt.Errorf("failed to query existing segment: %w", err)
		}

		if _, err = tx.Exec(
			`INSERT INTO user_segments(user_id, segment_slug, expiration_date) VALUES($1, $2, $3)
             ON CONFLICT (user_id, segment_slug) DO NOTHING`,
			userID,
			segment.Slug,
			segment.ExpirationDate,
		); err != nil {
			return 0, fmt.Errorf("failed to add segment '%s': %w", segment.Slug, err)
		}
	}

	// Удаляем сегменты
	for _, slug := range removeList {
		var existingSlug string
		err = tx.QueryRow("SELECT slug FROM segments WHERE slug = $1", slug).Scan(&existingSlug)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("segment with slug '%s' does not exist", slug)
		} else if err != nil {
			return 0, fmt.Errorf("failed to query existing segment: %w", err)
		}

		if _, err = tx.Exec(
			"DELETE FROM user_segments WHERE user_id=$1 AND segment_slug=$2",
			userID,
			slug,
		); err != nil {
			return 0, fmt.Errorf("failed to remove segment '%s': %w", slug, err)
		}
	}

	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

func (db *DB) GetUserSegments(userID int) (int, []string, error) {
	// Начало транзакции
	tx, err := db.db.Begin()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("An error occurred while rolling back the transaction: %v\n", err)
		}
	}()

	// Проверка наличия пользователя в базе данных
	var existingUserId int
	err = tx.QueryRow("SELECT id FROM users WHERE id = $1", userID).Scan(&existingUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil, fmt.Errorf("user with ID '%d' does not exist", userID)
	} else if err != nil {
		return 0, nil, fmt.Errorf("failed to query existing user: %w", err)
	}

	// Запрос на получение сегментов пользователя
	rows, err := tx.Query(
		"SELECT s.slug FROM segments s JOIN user_segments us ON s.slug = us.segment_slug WHERE us.user_id = $1",
		userID,
	)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to query segments for user ID '%d': %w", userID, err)
	}
	defer rows.Close()

	var segments []string
	// Обход результатов запроса и добавление их в массив segments
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return 0, nil, fmt.Errorf("failed to scan row for user ID '%d': %w", userID, err)
		}
		segments = append(segments, slug)
	}

	// Проверка наличия дополнительных ошибок, произошедших при получении всех строк запроса
	if err := rows.Err(); err != nil {
		return 0, nil, fmt.Errorf("error occurred while reading rows: %w", err)
	}

	// Завершение транзакции
	if err = tx.Commit(); err != nil {
		return 0, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, segments, nil
}

func (db *DB) GetUserReport(userID int, yearMonth string) (string, error) {
	// Начало транзакции
	tx, err := db.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("An error occurred while rolling back the transaction: %v\n", err)
		}
	}()

	// Проверка наличия пользователя в базе данных
	var existingUserId int
	err = tx.QueryRow("SELECT id FROM users WHERE id = $1", userID).Scan(&existingUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("user with ID '%d' does not exist", userID)
	} else if err != nil {
		return "", fmt.Errorf("failed to query existing user: %w", err)
	}

	// Создание директории для отчетов, если она не существует
	err = os.MkdirAll("reports", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create reports directory: %w", err)
	}

	// Создание CSV файла
	fileName := fmt.Sprintf("user_%d_report_%s.csv", userID, yearMonth)
	file, err := os.Create(filepath.Join("reports", fileName))
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file for user ID '%d' and year-month '%s': %w", userID, yearMonth, err)
	}
	defer file.Close()

	// Инициализация CSV writer
	w := csv.NewWriter(file)
	defer w.Flush()

	// Запись заголовков CSV файла
	if err := w.Write([]string{"User ID", "Segment Slug", "Operation", "Operation Date"}); err != nil {
		return "", fmt.Errorf("failed to write headers to CSV: %w", err)
	}

	// Выборка данных для отчета из базы данных
	rows, err := tx.Query(
		`SELECT user_id, segment_slug, operation, operation_date 
         FROM user_segment_history 
         WHERE user_id = $1 AND to_char(operation_date, 'YYYY-MM') = $2`,
		userID,
		yearMonth,
	)
	if err != nil {
		return "", fmt.Errorf("failed to query user_segment_history for user ID '%d' and year-month '%s': %w", userID, yearMonth, err)
	}
	defer rows.Close()

	// Запись данных в CSV файл
	for rows.Next() {
		var id int
		var slug, operation, operationDate string
		if err := rows.Scan(&id, &slug, &operation, &operationDate); err != nil {
			return "", fmt.Errorf("failed to scan row for user ID '%d': %w", userID, err)
		}
		if err := w.Write([]string{strconv.Itoa(id), slug, operation, operationDate}); err != nil {
			return "", fmt.Errorf("failed to write row to CSV for user ID '%d': %w", userID, err)
		}
	}

	// Проверка наличия дополнительных ошибок
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error occurred while reading rows: %w", err)
	}

	// Подтверждение транзакции
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return fileName, nil
}
