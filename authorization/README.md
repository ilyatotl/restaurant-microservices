# Сервис авторизации

## Запуск
Для запуска необходимо запустить контейнер с базой данных, а затем запустить сам сервис:  
```
dpcker-compose build
docker-compose up -d
go run cmd/app/main.go
```

## Контрольные точки
- `/register/{role}`  
В данной контрольной точке происходит регистрация пользователя. В качестве role указывается роль пользователя в
ресторане (customer, chef, manager).
В теле запроса необходимо передать имя пользователя, почту и пароль. Пример:
```json
{
  "username": "user",
  "email": "aaa@mail.ru",
  "password": "12345"
}
```
- `/authorize`  
В данной ручке происходит авторизация пользователя. В теле запроса необходимо передать пароль и почту. Приемр:  
```json
{
  "email": "aaa@mail.ru",
  "password": "12345"
}
```  
В качестве ответа, при успешной авторизации пользователь получит JWT-токен.
- `/get_user_info`  
Данная контрольная точка возвращает информацию о пользователе (id, имя пользователя, почта, роль). Для получения данных
в заголовок `Authorization` необходимо передать JWT-токен