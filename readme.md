# Курсовой проект по курсу Backend разработка на Go уровень 1

## Запуск
1. Развертка базы в докере - `make init_db`
2. Запуск юнит-тестов `make test-unit`
3 Запуск интеграционных тестов  ` make test-integration `
4. Запуск проекта с конфигом в переменных окружения (подготовлено для Heroku) - `make run`
5. Запуск приложения в контейнере c пересборкой всех контейнеров без кэша - `make run-docker-full`
6. Запуск созданных контейнеров - `make run-docker`


## API
Описание всех эндпоинтов в формате openapi в каталоге `api`

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
