# Сервис обработки заказов

## Запуск
Для запуска необходимо запустить контейнер с базой данных, а затем запустить сам сервис:
```
dpcker-compose build
docker-compose up -d
go run cmd/app/main.go
```

## Контрольные точки
- `/order`  
В данной контрольной точке производится заказ блюд. Для этого необходимо, чтобы пользователь был
авторизован (в заголовке `Authorization` находился валидный JWT-токен), а также в тело запроса
передавались данные заказа: id блюд и их количество.  
Прмиер тела запроса:
```json
[
  {
    "id": 1,
    "quantity": 5
  },
  {
    "id": 3,
    "quantity": 2
  }
]
```
В качестве ответа сервер вернет статус заказа (успешно или не успешно)
- `/get_order_info`  
Данная контрольная точка возвращает информацию о заказе. Пользователь должен передать в
заголовок `Authorization` свой JWT-токен, а также в query-параметры запроса id заказа. И, если
пользователь имеет доступ к заказу, ему вернется информация о нем.
- `/menu`  
Данная контрольная точка возвращает меню: все доступные для заказа блюда и информацию о них
  (пользователь также должен быть авторизован, передавая JWT-токен в заголовок `Authorization`)
- `/dish/add`  
Данная контрольная точка предназначена для менеджера, и позволяет добавить новое блюдо в систему.  
Для этого менеджер должен передать свой JWT-токен в заголовок `Authorization`, а также в теле запроса
указать данные нового блюда.  
Пример тела запроса:
```json
{
  "name": "potato",
  "description": "tasty potato",
  "price": 100,
  "quantity": 7
}
```
- `/dish/get`  
  Данная контрольная точка предназначена для менеджера, и позволяет получить информацию о блюде.  
  Для этого менеджер должен передать свой JWT-токен в заголовок `Authorization`, а также в 
  query-параметрах запроса указать id блюда.
- `/dish/update`
  Данная контрольная точка предназначена для менеджера, и позволяет обновить информацию о блюде.  
  Для этого менеджер должен передать свой JWT-токен в заголовок `Authorization`, а также в теле запроса
  передать всю информацию о блюде, и для блюда с переданным id информация будет заменена.  
  Пример тела запроса:
```json
{
  "id": 1,
  "name": "potato",
  "description": "very very tasty potato",
  "price": 80,
  "quantity": 120
}
```
- `/dish/delete`  
Данная контрольная точка предназначена для менеджера, и позволяет удалить блюдо из системы.  
Для этого менеджер должен передать свой JWT-токен в заголовок `Authorization`, а также в
query-параметрах запроса указать id блюда.