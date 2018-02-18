# Тестовое задание для IQ Option

Реализация in-memory key-value cache storage: сервера и простого клиента с web-интерфейсом. Поддерживаемые форматы данных: строка, список строк, ассоциативный массив строк.
Rest-API хранилища поддерживает:
  - Добавление пары ключ-значение
  - Получение значения по ключу
  - Удаление ключа
  - Получение списка ключей, соответствующих паттерну
  - Установка времени жизни ключа (в секундах)
  - Получение значения по индексу (для списков и массивов)
  - Автосохранение кеша в файл и загрузка из файла
  - Авторизация

### API

Спецификация описана в формате OpenAPI, ссылка на файл спецификации <https://github.com/andreipimenov/kvstore/blob/master/doc/swagger.yaml>

Запуск локального SwaggerUI. Доступ через браузер на 127.0.0.1:3000
```
docker run -p 3000:8080 -e "API_URL=swagger.yaml" -v $(pwd)/doc/swagger.yaml:/usr/share/nginx/html/swagger.yaml swaggerapi/swagger-ui
```

#### Примеры запросов к API
Создание пары ключ-значение
```
curl -X POST -d '{"key":"name", "value": "John Doe"}' 127.0.0.1:8080/api/v1/keys

{"message":"OK"}
```
Установка времени "жизни" ключа и последующее получение этого времени
```
curl -X POST -d '{"expires":100}' 127.0.0.1:8080/api/v1/keys/name/expires

{"message":"OK"}

curl -X GET 127.0.0.1:8080/api/v1/keys/name/expires

{"expires":49}
```
Получение всех ключей, соответствующих паттерну
```
curl -X GET 127.0.0.1:8080/api/v1/keys/h*

{"keys":["hello","hell"]}
```
Получение значения ключа
```
curl -X GET 127.0.0.1:8080/api/v1/keys/hell/values

{"value":"Java"}
```

### Сборка и запуск

Билд сервера и клиента
```
go build -o ./bin/server ./cmd/server
go build -o ./bin/client ./cmd/client
```
Сервер и клиент могут быть сконфигурированы при запуске с помощью флагов:
 - port — порт приложения (по умолчанию, 8080 для сервера и 8090 для клиента)
 - config — путь к файлу конфигурации (например, файл конфигурации сервера <https://github.com/andreipimenov/kvstore/blob/master/etc/server.conf.json>)
 Дополнительно при запуске клиента можно указать флаг:
 - server — строка в формате host:port для связи с сервером (например, 127.0.0.1:8080)

Пример создания ключа через веб-интерфейс клиента
![](https://github.com/andreipimenov/kvstore/blob/master/asset/client.example.jpg)

Запуск сервера в виде systemd сервиса.
Например, запуск сервера от пользователя user:
```
nano /usr/lib/systemd/system/kvserver.service
```
Содержимое файла
```
[Unit]
Description=KVServer

[Service]
Restart=always
RestartSec=10
WorkingDirectory=/home/u/
ExecStart=/home/u/server
LimitNOFILE=524576
LimitNPROC=524576
User=u
Group=u
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=KVServer

[Install]
WantedBy=multi-user.target
```
Добавление приложения в автозапуск и старт
```
systemctl start kvserver
systemctl enable kvserver
```