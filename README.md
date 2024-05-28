# Сервис динамического сегментирования пользователей
## О проекте
Сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент и тд)

## Api предоставляет следующий функционал
- /create-seg: Метод создания сегмента. Принимает slug (название) сегмента;
- /delete-seg: Метод удаления сегмента. Принимает slug (название) сегмента;
- /add-user-segments: Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя. Опционально принимает дату, когда пользователь должен быть удален из сегмента;
- /show-user-segments: Метод получения активных сегментов пользователя. Принимает на вход id пользователя;
- /report: Метод выгрузки аудита добавления/удаления пользователей из сегмента в формате CSV. Принимает дату вида: n year n month.

## Модель базы данных
Пожалуйста, ознакомьтесь со скриптом модели базы данных в папке materials и примените его к вашей базе данных (вы можете использовать командную строку с psql или просто запустить его через любую IDE, например DataGrip от JetBrains или pgAdmin от сообщества PostgreSQL).

## Примеры работы с api
1. curl --header "Content-Type: application/json" --request POST --data '{"userId":1}' http://localhost:8089/show-user-segments
2. curl --header "Content-Type: application/json" --request POST --data '{"userId":3,"addSlug":["DISCOUNT_50"],"ttl":"2h"}' http://localhost:8089/add-user-segments
3. curl --header "Content-Type: application/json" --request POST --data '{"userId":3,"deleteSlug":["DISCOUNT_30", "DISCOUNT_50"]}' http://localhost:8089/add-user-segments
4. curl -o report.csv --header "Content-Type: application/json" --request POST --data '{"year":2024,"month":5}' http://localhost:8089/report
