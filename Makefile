.PHONY: build dc run clean test fmt vet lint help coverage-html coverage

## migrate-up: миграции up
migrate-up:
	migrate -path migrations -database "postgres://localhost:5432/segmentation?sslmode=disable" up

## migrate-down: миграции down
migrate-down:
	migrate -path migrations -database "postgres://localhost:5432/segmentation?sslmode=disable" down

## build: Билдит бинарный файл
build:
	go build -o bin/app -v cmd/resttest/main.go

## dc: Докер компоуз билд
dc:
	docker-compose up --remove-orphans --build

## run: Запускает app нативно
run:
	go run cmd/resttest/main.go

## clean: Очищяет и удаляет бинарный файл
clean:
	go clean
	rm -f bin/app

## test: Запускает все тесты
test:
	go test -v -race ./...

## cover-html: запускает тесты с получением отчёта в html формате
cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

 ## cover: Запускает тесты с покрытием
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

## fmt: Форматирование кода для соответствия стандартному стилю Go
fmt:
	go fmt ./...

## vet: Статический анализ кода на поиск подозрительных конструкций
vet:
	go vet ./...

## lint: Запускает линтер
lint:
	golangci-lint run

help: Makefile
	@echo " Choose a command run in "RESTtest":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'



