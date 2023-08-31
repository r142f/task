# Сервис динамического сегментирования пользователей

### Запуск сервиса:
```bash
docker-compose up
```

### Запуск тестов:
```bash
docker-compose --file docker-compose_test.yml up
```

### Примеры запросов:
1. Метод создания сегмента.
    - Запрос: 
        ```bash
        curl --header "Content-Type: application/json" \
             --request POST \
             --data '{"name":"AVITO_VOICE_MESSAGES"}' \
             http://localhost:8080/createSegment
        ```
    - Ответ:
        ```bash
        HTTP/1.1 201 Created
        ```
2. Метод удаления сегмента. Принимает название сегмента (`name`). 
    - Запрос: 
        ```bash
        curl --header "Content-Type: application/json" \
             --request DELETE \
             --data '{"name":"AVITO_VOICE_MESSAGES"}' \
             http://localhost:8080/deleteSegment
        ```
    - Ответ:
        ```bash
        HTTP/1.1 200 OK
        ```
3. Метод добавления пользователя в сегмент. Принимает список названий сегментов, которые нужно добавить пользователю (`segmentsToAdd`), список названий сегментов, которые нужно удалить у пользователя (`segmentsToDelete`), id пользователя (`userId`).
    - Запрос: 
        ```bash
        curl --header "Content-Type: application/json" \
             --request POST \
             --data '{"segmentsToAdd":["AVITO_VOICE_MESSAGES"], "segmentsToDelete": ["AVITO_DISCOUNT_30"], "userId": 1}' \
             http://localhost:8080/updateUserSegments
        ```
    - Ответ:
        ```bash
        HTTP/1.1 201 CREATED
        ```
4. Метод получения активных сегментов пользователя. Принимает на вход id пользователя.
    - Запрос: 
        ```bash
        curl http://localhost:8080/userSegments?userId=1
        ```
    - Ответ:
        ```bash
        HTTP/1.1 200 OK
        Content-Type: application/json
        ...
        ["AVITO_VOICE_MESSAGES"]
        ```