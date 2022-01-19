# Курсовой проект по курсу Backend разработка на Go уровень 1

1. Запуск проекта
```bash
make run
```
2. Запуск тестирования
```bash
make test
```
3. Запуск интеграционного тестирования
```bash
make test-integration
```

Работает один endpoint `user/create/`. Пример запроса:
```bash
curl -X 'POST'  http://localhost:8000/user/create -H "Content-Type: application/json" -d '{
  "name": "string",    
  "password": "string",
  "email": "string"
}' 
```
