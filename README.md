# Сервис динамического сегментирования пользователей

### Запуск сервиса:
- Склонируйте репозиторий. 
- В репозитории проекта
    ```bash
    docker-compose up --build
    ```
- Серис запущен на `localhost:8080`

### Запуск тестов:
- Склонируйте репозиторий. 
- В репозитории проекта
    ```bash
    docker-compose --file docker-compose_test.yml up --build --abort-on-container-exit
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
        Content-Type: application/json
        ...
        {"name":"AVITO_VOICE_MESSAGES"}
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
4. Метод получения активных сегментов пользователя. Принимает на вход id пользователя (`userId`).
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
5. Метод формирования отчета с историей попадания/выбывания пользователя из сегмента. Принимает на вход id пользователя (`userId`), год (`year`) и месяц (`month`).
    - Запрос:
        ```bash
        curl "http://localhost:8080/generateReport?year=2023&month=9&userId=1"
        ```
    - Ответ:
        ```bash
        HTTP/1.1 200 OK
        Content-Type: application/json
        ...
        {"link":"reports/9849d4bb-cf38-491b-a45a-b83af811046c.csv"}
        ```
6. Метод получения отчета. Принимает на вход `Link` из ответа на предыдущий запрос.
    - Запрос: 
        ```bash
        curl "http://localhost:8080/reports/9849d4bb-cf38-491b-a45a-b83af811046c.csv"
        ```
    - Ответ:
        ```bash
        HTTP/1.1 200 OK
        Content-Type: text/csv
        ...
        1,AVITO_VOICE_MESSAGES,add,2023-09-01 17:49:03
        ```