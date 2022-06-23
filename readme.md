# Курсовой проект по курсу Backend разработка на Go уровень 1

__ЗАДАЧА__: Реализовать REST-сервис для сокращение url со статистикой.

### Проектирование
1. [спецификация по стандарту `openapi`](./api/api.yml).

2. библиотка [Echo](https://echo.labstack.com/) для реализации сервера.

3. Postgres в качестве хранилища данных

### Реализация

#### Оновная логика 
В основе представления данных лежат [модели](./internal/pkg/models/), которые используются в каждом слое приложения для представления данных.

Работа с данными осуществляется на 3х уровнях:
 * Repository - работа с хранилищем
 * Usecase - бизнес логика
 * Deliver - хэндлеры для доставки через api

Связь между уровнями осуществляется через интерфейсы.

#### Дополнительные инструменты
 * Для логгирования на каждом этапе используется [отдельный пакет](./internal/logger/logger.go), построенный на основе библиотеки [`zap`](https://github.com/uber-go/zap)
 * Конфигурации можно задавать в `.yaml` файле или через переменные окружения. За работу с конфигом отвечает [соответствующий пакет](./internal/config/)

## Запуск
1. Развертка базы в докере - `make init_db`
2. Запуск юнит-тестов `make test-unit`
3 Запуск интеграционных тестов  ` make test-integration `
4. Запуск проекта с конфигом в переменных окружения (подготовлено для Heroku) - `make run`
5. Запуск приложения в контейнере c пересборкой всех контейнеров без кэша - `make run-docker-full`
6. Запуск созданных контейнеров - `make run-docker`



### Примерные запросы к api:

```bash
curl -X POST localhost:8080/user/create -H "Content-Type: application/json" -d '{
   "name": "login",
  "password": "123",
  "email": "login@mail.ru"
}'

curl -X GET "http://localhost:8080/user/login?name=login&password=123"


curl -X POST http://localhost:8080/url -d '{
                "rawurl": "http://url.ru",    
                "userid": 1
          }' -H "Authorization: Bearer <JWT>'

curl -X GET http://localhost:8080/<shortened>

curl -X GET http://localhost:8080/redirects/<shortened>
```

## Веб-интерфейс по [ссылке](https://shrtnergb.herokuapp.com/web/generate) 
