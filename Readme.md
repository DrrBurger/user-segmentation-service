# User Segmentation service

Микросервис для динамической работы с сегментами пользователей, присвоением/удалением сегмента/ов пользователю с возможностью установки TTL, историей присвоения того или иного сегмента, а так же
получением отчета по пользователю за определенный период.

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Gin (веб фреймворк)
- golang-migrate/migrate (для миграций БД)
- pq (драйвер для работы с PostgreSQL)
- golang/mock, testify (для тестирования)
- в качестве логера использовал builtin log (при необходимости заменить на желаемый log/slog, logrus, и т.д)

# Getting Started

Для подготовки сервиса необходимо:

- Клонировать репозиторий
```
    git clone https://github.com/DrrBurger/user-segmentation-service
```
- Заполнить .env файл,
- Опционально, настроить `congig/config.yaml` под себя


# Usage

Запустить сервис можно с помощью команды `make dc` (запуск через докер)

Для запуска локально в файле main.go поменять аргумент "cfg.PG.URL" на "cfg.PG.URLLocal" в функции connectToDB и выполнить две команды последовательно:
1) `make migrate-up` 
2) `make run`

Для запуска тестов необходимо выполнить одну из команд:
1. `make test` запуск всех тестов
2. `make cover` для запуска тестов с покрытием
3. `make cover-html` для запуска тестов с покрытием и получения отчёта в html формате

Для запуска линтера необходимо выполнить команду `make lint`

Остальные команды можно получить выполнив команду `make help`

## Examples

Некоторые примеры запросов
- [Создание пользователя](#create-user)
- [Удаление пользователя](#del-user)
- [Создание сегмента](#create-seg)
- [Удаление сегмента](#del-seg)
- [Добавление/Удаление сегментов](#add-remove)
- [Получение списка сегментов](#seg-list)
- [Получение истории пользователя](#user-history)
- [Вопросы во время разработки](#decisions)


### Создание пользователя <a name="create-user"></a>

Создание пользователя с указанным именем:
```curl
curl --location --request POST 'http://localhost:8080/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Maks"
}'
```
Пример ответа:
```json
{
  "message": "User created successfully",
  "user_id": 1
}
```

### Удаление пользователя <a name="del-user"></a>

Удаление пользователя по указанному user_id:
```curl
curl --location --request DELETE 'http://localhost:8080/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": 1
}'
```
Пример ответа:
```json
{
   "message": "User deleted successfully",
   "user_id": 1
}
```

### Создание сегмента <a name="create-seg"></a>

При создании сегмента реализована опция указания процента пользователей (из общего колличества), которые попадут в этот сегмент автоматически, а так же есть возможность установить TTL:
```curl
curl --location --request POST 'http://localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data-raw '{
    "slug": "AVITO_SALE_60",
    "expiration_date": "2023-12-31T23:59:59Z",
    "random_percentage": 0.0
}'
```
Пример ответа:
```json
{
   "message": "Segment and user assignments created successfully"
}
```

### Удаление сегмента <a name="del-seg"></a>

Удаление сегмента по указанному slug:
```curl
curl --location --request DELETE 'http://localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data-raw '{
    "slug": "AVITO_SALE_60"
}'
```
Пример ответа:
```json
{
   "message": "Segment deleted successfully",
   "segment_id": 1
}
```

### Добавление/Удаление сегментов <a name="#add-remove"></a>

Добавление / удаление сегментов пользователя списком без перетирания существующих сегментов с возможностью установить TTL.
```curl
curl --location --request POST 'http://localhost:8080/user/segments' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id": 1,
    "add": [
    {
      "slug": "AVITO_SALE_10",
      "expiration_date": "2023-12-31T23:59:59Z"
    },
    {
      "slug": "AVITO_SALE_30",
      "expiration_date": "2023-11-30T23:59:59Z"
    },
    {
      "slug": "AVITO_SALE_20",
      "expiration_date": "2023-11-30T23:59:59Z"
    }
    ], 
    "remove": ["AVITO_SALE_40"]
}'
```
Пример ответа:
```json
{
   "message": "User segments updated successfully",
   "user_id": 1
}

```

### Получение списка сегментов <a name="seg-list"></a>

Получение списка сегментов пользователя по id:
```curl
curl --location --request GET 'http://localhost:8080/user/segments' \
--header 'Content-Type: application/json' \
--data-raw '{
   "user_id": 1
}'
```
Пример ответа:
```json
{
   "segments": ["AVITO_SALE_10","AVITO_SALE_30"],
   "user_id": 1
}
```

### Получение истории пользователя <a name="user-history"></a>

Получение отчета по указанным (user_id и период) в формате CSV.
```curl
curl --location --request GET 'http://localhost:8080/user/segments' \
--header 'Content-Type: application/json' \
--data-raw '{
   "user_id": 3,
   "yearMonth":"2023-08"
}'
```
Пример ответа:
```json
{
  "download_link": "http://localhost:8080/user/report/user_3_report_2023-08.csv",
  "message": "Report generated successfully"
}
```

# Decisions <a name="decisions"></a>

В ходе разработки были сомнения по тем или иным вопросам, которые были решены следующим образом:

1. При создании пользователя стоит ли автоматически присваивать ему случайные сегменты?
> Решил, что не стоит, т.к. возможно случайное попадание пользоватлей не в те сегменты. Но возможно в будущем стоит добавить эту возможность с дополнительными проверками
2. При обновлении сегментов пользователя стоит ли разделить операции добавления и удаления сегментов?
> Решил не разделять так как в одном запросе операция происходит быстрее
3. Стоит ли добавлять отдельные сервисы для хранения отчетов и как реализовать возврат отчёта по ссылке??
> Решил что возврат csv в http будет быстрее и исключит лишние зависимости, но реализовал сохранение файла локально для
> дальнейшей интеграции с сервисами/базами для хранения файлов
4. Как реализовать TTL?
> Для полноценной реализации функционала с TTL необходимо (реализовать планировщик, использовать индексы для оптимизации или иной известный способ)
5. Реализовать чистую архитектуру полностью?
> Решил полностью не реализовывать так как сервис маленький и решил оставить два слоя для сохранения читаемости кода. 