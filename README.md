# Practice0 - Демонстрационный сервис заказов с Kafka, PostgreSQL и кешем

Демонстрационный микросервис на Go, который получает заказы через Kafka, сохраняет их в PostgreSQL, кеширует в памяти и предоставляет простой HTTP API для получения данных о заказах.

---

## Функционал

- Kafka consumer и producer для сообщений о заказах
- Интеграция с PostgreSQL для хранения данных
- Кеш в памяти с временем
- HTTP API:
    - `/orders` — возвращает все кешированные заказы
    - `/order/{id}` — возвращает заказ по ID
    - `/` — корневой эндпоинт с базовой информацией

---


## Переменные окружения

Создайте файл `.env` в корне проекта:

```env
DB_URL=postgres://user:password@localhost:5433/orders?sslmode=disable
KAFKA_URL=localhost:9092
KAFKA_TOPIC=orders
KAFKA_GROUP_ID=order-consumers
PORT=8081
DEBUG=true
```
## Поднятие Docker окружения
```
docker compose up
docker ps 
```
После запуска должны быть подняты три контейнера: PostgreSQL, Kafka и Zookeeper
должно быть поднято три сервера

## Миграции базы данныхи
Для применения миграций используйте:
```
migrate -path migrations -database "postgres://user:password@localhost:5433/orders?sslmode=disable" up
```

## Добавление топика в кафку
```
docker exec -it kafka kafka-topics --create \
  --topic orders \
  --bootstrap-server localhost:9092 \
  --partitions 1 \
  --replication-factor 1
```
