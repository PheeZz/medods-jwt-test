## Пример .env файла
```ini
JWT_SECRET_KEY=SECRET
MONGO_URI=mongodb://localhost:27017/
ACCESS_TOKEN_DURATION_HOURS=24
REFRESH_TOKEN_DURATION_HOURS=72
DATABASE_NAME=medods_jwt_test
```

## Запуск

```bash
go run github.com/pheezz/medods-jwt-test/cmd/medods-jwt-test
```

## Доступные эндпоинты

- GET 127.0.0.1:8080/keyPair?GUID={GUID}
- POST 127.0.0.1:8080/refreshPair

## Связь
- [Telegram](https://t.me/pheezz)