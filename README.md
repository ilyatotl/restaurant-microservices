# ДЗ 4. Система обработки заказов ресторана
## Мерочкин Илья БПИ218

### Запуск
Запускать каждый микросервис нужно отедльно: инструкция по запуску каждого находится в соответствующих
README-файлах.

### Архитектура
В данном проекте реализованы два микросервиса: микросервис регистрации/авторизации и микросеврвис обработки заказов.  
Для каждого из них запускается собственная база данных, а реализованы они на основе RESTful API. В
 конечные точки микросервисов (описаны в README для каждого микросервиса) поступают запросы по HTTP. 
Сами микросервисы взаимодействуют друг с другом с помощью gRPC.
