asset_storage
====================
Тестовое задание - см [task.pdf](task.pdf)
Ветка для реализации дополнительных вопросов - development

Системные требования
====================

Код проверялся на Fedora 39 Linux, возможно он будет работать и на других дистрибутивах.
Требования
- GNU Make `dnf install make`
- Golang `dnf install golang`
- Podman - `dnf install podman podman-compose podman-plugins containernetworking-plugins`
- Curl `dnf install curl`
- Goose - https://github.com/pressly/goose - инструмент для миграций базы данных.

Теоритически podman можно заменить `docker`+ `docker-compose-plugin` но такой подход не 
рекомендуется для моего дистрибутива.

Как подготовить базу данных?
======================

1. Убедитесь, что установлен [goose](https://github.com/pressly/goose)

```shell
vodolaz095@steel:~/projects/asset_storage$ which goose
~/bin/goose
vodolaz095@steel:~/projects/asset_storage$ goose 
Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND
```

2. Запустите базу данных PostgreSQL в контейнере (или докером, или подманом, а может быть, и с хоста)

```shell

make podman/resource
make docker/resource

```

3. Убедитесь, что `goose` может соединиться с базой данных. Возможно, придётся отредактировать строку соединения в 
   [migrate.mk](make%2Fmigrate.mk) если база данных Postgres запущена по другому...
```shell

$ make migrate/info

goose --dir ./migrations/ postgres "user=assets password=assets dbname=assets host=localhost port=5432 sslmode=disable" status
2024/11/18 21:52:08     Applied At                  Migration
2024/11/18 21:52:08     =======================================
2024/11/18 21:52:08     Pending                  -- 00001_users.sql
2024/11/18 21:52:08     Pending                  -- 00002_sessions.sql
2024/11/18 21:52:08     Pending                  -- 00003_assets.sql
2024/11/18 21:52:08     Pending                  -- 00004_test_user.sql


```
4. Примените миграции к базе данных - запустите `make migrate/up` 4 раза
```shell

$ make migrate/up
goose --dir ./migrations/ postgres "user=assets password=assets dbname=assets host=localhost port=5432 sslmode=disable" up
2024/11/18 21:53:52 OK    00001_users.sql
2024/11/18 21:53:52 OK    00002_sessions.sql
2024/11/18 21:53:52 OK    00003_assets.sql
2024/11/18 21:53:52 OK    00004_test_user.sql
2024/11/18 21:53:52 goose: no migrations to run. current version: 4
```

  Если вызвать `make migrate/info` то покажет такое

```shell
$ make migrate/info 
goose --dir ./migrations/ postgres "user=assets password=assets dbname=assets host=localhost port=5432 sslmode=disable" status
2024/11/18 21:53:58     Applied At                  Migration
2024/11/18 21:53:58     =======================================
2024/11/18 21:53:58     Mon Nov 18 18:53:52 2024 -- 00001_users.sql
2024/11/18 21:53:58     Mon Nov 18 18:53:52 2024 -- 00002_sessions.sql
2024/11/18 21:53:58     Mon Nov 18 18:53:52 2024 -- 00003_assets.sql
2024/11/18 21:53:58     Mon Nov 18 18:53:52 2024 -- 00004_test_user.sql

```

5. Убедитесь, что в базе данных есть 4 таблицы
```
   assets
   goose_db_version
   sessions
   users
```

6. Убедитесь, что в таблице `users` есть тестовый пользователь `alice`/`secret`. Запрос
```sql
select * from users;
```

Возвращает как минимум одного пользователя, у которого будет 
```
login=alice
password_hash=5ebe2294ecd0e0f08eab7690d2a6ee69
```

7. База данных готова!

Как запустить приложение, если на хост машине есть компилятор Go?
======================================

Убедитесь, что версия компилятора не ниже 1.22.8 (на других не проверял)
```shell
vodolaz095@steel:~/projects/asset_storage$ go version
go version go1.22.8 linux/amd64
```

Запустите приложение: `make start` или `go run main.go`.
Конфигурацию можно задать через переменные окружения POSIX - см. [config.go](config%2Fconfig.go)
По умолчанию приложение соединяется с базой данных на хост машине и включает HTTP сервер на http://localhost:3000
Приложение пишет логи в STDOUT и они должны быть видны на экране.

Как запустить приложение с помощью docker/podman?
========================================
Возможно, придётся отредактировать [docker-compose.yml](docker-compose.yml)

```shell
$ make podman/up
$ make docker/up

```
По умолчанию приложение соединяется с базой данных на в контейнере и включает HTTP сервер на 3000 порту хост машины, то есть,
на http://localhost:3000. Приложение пишет логи в STDOUT и их можно увидеть стандартными средствами docker/podman.

Как протестировать приложение?
======================================
Приложение тестируется в полуавтоматическом режиме путём отправки HTTP запросов на http://localhost:3000 с помощью curl.
Вызовы curl можно упростить, используя рецепты из [integration.mk](make%2Fintegration.mk).

Примерный ход тестирования приложения

1. Успешно создаём сессию как пользователь `alice` с паролем `secret`
```shell
$ make integration/auth_ok

curl -v --data '{"login":"alice","password":"secret"}' http://localhost:3000/api/auth
* processing: http://localhost:3000/api/auth
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
> POST /api/auth HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/8.2.1
> Accept: */*
> Content-Length: 37
> Content-Type: application/x-www-form-urlencoded
> 
< HTTP/1.1 200 OK
< Date: Mon, 18 Nov 2024 18:36:23 GMT
< Content-Length: 44
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
{"token":"d3074facf8386b374cb74098071583d9"}

```

В выводе важна строка с токеном.
Токен должен быть создан в таблице sessions c UID пользователя `alice`.
После создания токена можно обновить переменную `session_good_token` в [integration.mk](make%2Fintegration.mk).

В логе приложение должно написать
```
22:25:47 transport_login.go:35: ASSET: Поступил запрос на авторизацию с [::1]:36652
22:25:47 authentication.go:18: ASSET: Пользователь alice пытается авторизоваться...
22:25:47 authentication.go:24: ASSET: Пользователь alice создал сессию.
```

2. Успешно создаём новый объект

```shell

$ make integration/create_ok 
curl -v -H "Authorization: Bearer "c38fcc2c923d818b1ba765eaf5053d18"" --data body_46 "http://localhost:3000/"api/upload-asset/key_46
* processing: http://localhost:3000/api/upload-asset/key_46
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
> POST /api/upload-asset/key_46 HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/8.2.1
> Accept: */*
> Authorization: Bearer c38fcc2c923d818b1ba765eaf5053d18
> Content-Length: 7
> Content-Type: application/x-www-form-urlencoded
> 
< HTTP/1.1 201 Created
< Location: /api/asset/key_46
< Date: Mon, 18 Nov 2024 20:09:46 GMT
< Content-Length: 19
< Content-Type: text/plain; charset=utf-8
< 
{
 "status": "ok"
* Connection #0 to host localhost left intact
}

```

В логе приложения будет написано
```

22:28:42 authentication.go:34: ASSET: Сессия востановлена для пользователя: alice
22:28:42 transport_upload.go:45: ASSET: Пользователь alice восстановлен из сессии
22:28:42 transport_upload.go:53: ASSET: Пользователь alice пытается загрузить 7 байт в ключ key_42
22:28:42 assets.go:32: ASSET: Пользователь alice пытается создать данные (7 байт) по ключу key_42
22:28:42 assets.go:41: ASSET: Пользователь alice создал данные по ключу key_42 (7 байт)
22:28:42 transport_upload.go:62: ASSET: Пользователь alice успешно загрузил 7 байт в ключ key_42

```

В базе данных в таблице `assets` будет создана запись c ключом, в данном случае `key_42`.
В файле [integration.mk](make%2Fintegration.mk) надо обновить переменную `asset_key_good`

3. Пробуем получить объект по ключу
```shell
$ make integration/get_ok 

curl -v -H "Authorization: Bearer "c38fcc2c923d818b1ba765eaf5053d18"" "http://localhost:3000/"api/asset/"key_42"
* processing: http://localhost:3000/api/asset/key_42
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
> GET /api/asset/key_42 HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/8.2.1
> Accept: */*
> Authorization: Bearer c38fcc2c923d818b1ba765eaf5053d18
> 
< HTTP/1.1 200 OK
< Date: Mon, 18 Nov 2024 19:32:13 GMT
< Content-Length: 7
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
body_42

```

В лог программы будет написано:

```
22:32:13 authentication.go:34: ASSET: Сессия востановлена для пользователя: alice
22:32:13 transport_get.go:43: ASSET: Пользователь alice восстановлен из сессии
22:32:13 assets.go:17: ASSET: Пользователь alice пытается получить данные по ключу key_42
22:32:13 assets.go:25: ASSET: Пользователь alice получил данные по ключу key_42 (7 байт)
22:32:13 transport_get.go:70: ASSET: Пользователь alice запросил данные по ключу key_42 и получил 7 байт
```


Ответы на дополнительные вопросы
==========================================

**Что можно улучшить в схеме БД?**
1. Хешировать пароль нормальным алгоритмом [argon2id](https://github.com/alexedwards/argon2id) на стороне приложения, а не в базе данных
2. Добавить [constraints](https://www.postgresql.org/docs/current/ddl-constraints.html), чтобы не было сессий и ассетов без пользователя
   Реализовано в https://github.com/vodolaz095/asset_storage/commit/8684d7e5d46ccd8aa554674056376e0ef3bb6243

**Доработайте механизм авторизации таким образом, что бы в каждый момент времени у пользователя была активна только одна (последняя) сессия.**

Теоритически можно сделать поле `uid` таблицы `sessions` уникальным индексом - см
https://github.com/vodolaz095/asset_storage/commit/98c1763259ae67ff9566cbde2d23c9504cd3c864

**Ограничьте максимальное время пользовательской сессии до 24-х часов.**

сделал

**Добавьте в БД данные об IP адресе авторизованного пользователя.**
1. надо менять структуру базы данных - добавлять поле типа https://www.postgresql.org/docs/current/datatype-net-types.html#DATATYPE-INET в таблицу sessions
2. Если мы используем какие-либо reverse-proxy или балансировщики/CDN перед API, то реальный адрес клиента может быть или в
   заголовке `X-Forwarded-For` или в `CF-Connecting-IP` для Cloudflare. Также бы было неплохо указать список IP адресов reverse-proxy 
   запросам с которых мы доверяем см. https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies

реализовано в коммите
https://github.com/vodolaz095/asset_storage/commit/7fb5ec52578a2e8c59879290f7f3a1bc948f4174

**Реализуйте методы API для получения списка всех закаченных файлов**

Реализован в коммите https://github.com/vodolaz095/asset_storage/commit/2e8914f99aacce4bbca64242afa65453f4f6113f
```shell

$ make integration/list 

curl -v "http://localhost:3000/"api/list
* processing: http://localhost:3000/api/list
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
> GET /api/list HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/8.2.1
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Tue, 19 Nov 2024 17:13:41 GMT
< Content-Length: 338
< Content-Type: text/plain; charset=utf-8
< 
[
 {
  "name": "key_13",
  "author": "alice",
  "size": 7,
  "created_at": "2024-11-19T19:54:13.327996+03:00"
 },
 {
  "name": "key_16",
  "author": "alice",
  "size": 7,
  "created_at": "2024-11-19T19:54:16.46712+03:00"
 },
 {
  "name": "key_17",
  "author": "alice",
  "size": 7,
  "created_at": "2024-11-19T19:54:17.397095+03:00"
 }
]





```


**Реализуйте методы API для удаления файлов.**
долго

**Реализуйте работу сервера по протоколу HTTPS.**
См. [transport.go](internal%2Ftransport%2Fhttp%2Ftransport.go) - функция ListenHTTPS








