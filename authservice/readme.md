## Регистрация

REQUEST:

POST: /api/user/register

Body raw:
   
 ```json
{
    "telephone" : "8-953-983-0807",
    "password" : "toster123",
    "email" : "walkmanmail19@gmail.com",
    "name" : "Vlad"
}
```

RESPONSE:

Status: 200

 ```json
{
    "user_id": 66,
    "telephone": "8-953-983-0807",
    "email": "walkmanmail19@gmail.com",
    "name": "Vlad",
    "avatar_url": "http://127.0.0.1:4491/download/avatar/id/base",
    "access_token": ""
}
```

Status: 400, попытка зарегистировать, уже зарегистрированного пользователя

 ```json
{
    "error": "Some problems with data operation. "
}
```

Status: 400, получен неправильный json

 ```json
{
    "error": "Non valid data. "
}
```

## Авторизация

REQUEST:

POST: /api/user/authorized

Body raw:
   
 ```json
{
    "telephone" : "8-953-983-0807",
    "password" : "toster123"
}
```

RESPONSE:

Status: 200

 ```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTYwNzg3OTI4M30.i4WpzufRWSuFLhJ-q5nuBQjm2T7YSjoqFE771-tHVuc",
    "user": {
        "user_id": 66,
        "telephone": "8-953-983-0807",
        "email": "walkmanmail19@gmail.com",
        "name": "Vlad",
        "avatar_url": "http://127.0.0.1:4491/download/avatar/id/base",
        "access_token": ""
    }
}
```

Status: 400, попытка авторизивать пользователя, который уже авторизован, то есть пользователя, у которого не прогорел refresh токен

 ```json
{
    "error": "Non valid refresh token. "
}
```

Status: 400, получен неправильный json

 ```json
{
    "error": "Non valid data. "
}
```


## Проверка авторизации

REQUEST:

GET: /api/user/access

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTYwNzg3OTI4M30.i4WpzufRWSuFLhJ-q5nuBQjm2T7YSjoqFE771-tHVuc
```

RESPONSE:

Status: 200

 ```json
{
    "user_id": 66,
    "telephone": "8-953-983-0807",
    "email": "walkmanmail19@gmail.com",
    "name": "Vlad",
    "avatar_url": "http://127.0.0.1:4491/download/avatar/id/base",
    "access_token": ""
}
```

Status: 401, access токен прогорел.

 ```json
{
    "error": "Non valid access token. "
}
```

Status: 401, refresh токен прогорел.

 ```json
{
    "error": "Non valid refresh token. "
}
```

Status: 401, попытка проверить авторизацию, не существующего пользователя.

 ```json
{
    "error": "User isn't exist. "
}
```

Status: 401, при сбросе прогоревшего refresh токена, произошла ошибка.

 ```json
{
    "error": "Some problems with data operation. "
}
```

Status: 500, некорректный токен.

## Обновить токен

REQUEST:

GET: /api/user/access/update

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTYwNzg3OTI4M30.i4WpzufRWSuFLhJ-q5nuBQjm2T7YSjoqFE771-tHVuc
```

RESPONSE:

Status: 200

 ```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTYwNzg4MTM3Mn0.rL5kj4ryl2HORdZlcdEIUTKVn-x26ZqtWpHX-s6Uq9s",
    "user": {
        "user_id": 66,
        "telephone": "8-953-983-0807",
        "email": "walkmanmail19@gmail.com",
        "name": "Vlad",
        "avatar_url": "http://127.0.0.1:4491/download/avatar/id/base",
        "access_token": ""
    }
}
```

Status: 401, ошибка расшифровки access токена; неправильный json; ошибка создании пары токенов.

 ```json
{
    "error": "Non valid access token. "
}
```

Status: 401, refresh токен прогорел; refresh токен пустой.

 ```json
{
    "error": "Non valid refresh token. "
}
```

Status: 401, не существующий пользователь.

 ```json
{
    "error": "User isn't exist. "
}
```

Status: 401, при сбросе прогоревшего refresh токена, произошла ошибка.

 ```json
{
    "error": "Some problems with data operation. "
}
```

Status: 500, некорректный токен.

## Обновить пользователя

REQUEST:

POST: /api/user/update

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTYwNzg4MTM3Mn0.rL5kj4ryl2HORdZlcdEIUTKVn-x26ZqtWpHX-s6Uq9s
```

Body: raw

```json
{
    "user_id":66,
    "email":"vladislav.kuznetsovRTN1@yandex.ru"
}
```

RESPONSE:

Status: 200, при смене данных пользователя не являющимися именем или телефоном, возвращается прошлый токен.

 ```json
{
    "user_id": 66,
    "telephone": "8-953-983-0807",
    "email": "vladislav.kuznetsovRTN1@yandex.ru",
    "name": "Vlad",
    "avatar_url": "http://127.0.0.1:4491/download/avatar/id/base",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTYwNzg4MTM3Mn0.rL5kj4ryl2HORdZlcdEIUTKVn-x26ZqtWpHX-s6Uq9s"
}
```
REQUEST:

POST: /api/user/update

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiJWbGFkIn0sImV4cCI6MTYwNzg4MTM3Mn0.rL5kj4ryl2HORdZlcdEIUTKVn-x26ZqtWpHX-s6Uq9s
```

Body: raw

```json
{
    "user_id":66,
    "name":"Влад Кузнецов"
}
```

RESPONSE:

Status: 200, при смене данных пользователя являющимися именем или телефоном, возвращается новый токен.

 ```json
{
    "user_id": 66,
    "telephone": "8-953-983-0807",
    "email": "vladislav.kuznetsovRTN1@yandex.ru",
    "name": "Влад Кузнецов",
    "avatar_url": "http://127.0.0.1:4491/download/avatar/id/base",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTYwNzg4MjA2OH0.t9LWRayZbrBvO4WL_u9se3KJthYH9Y6Ns_HtpLp3C8I"
}
```

Status: 404, неправильные данные.

 ```json
{
    "error": "Non valid data. "
}
```

Status: 404, пользователя не найден.

 ```json
{
    "error": "User isn't exist. "
}
```

Status: 404, ошибка обновления пользователя или сброса refresh токена.

 ```json
{
    "error": "Some problems with data operation. "
}
```

Status: 401, access токен прогорел.

 ```json
{
    "error": "Non valid access token. "
}
```

Status: 401, refresh токен прогорел.

 ```json
{
    "error": "Non valid refresh token. "
}
```

## Вернуть пользователя

REQUEST:

POST: /api/user/get

Body: raw

```json
{
    "user_id":66
}
```

RESPONSE:

Status: 200.

```json
{
    "users": [
        {
            "user_id": 66,
            "telephone": "8-953-983-0807",
            "email": "vladislav.kuznetsovRTN1@yandex.ru",
            "name": "Влад Кузнецов",
            "avatar_url": "http://127.0.0.1:4491/download/avatar/id/base",
            "access_token": ""
        }
    ]
}
```

## Сбросить пароль

Получить товен сброса:

REQUEST:

POST: /api/user/password/token

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTYwNzg4MjA2OH0.t9LWRayZbrBvO4WL_u9se3KJthYH9Y6Ns_HtpLp3C8I
```

Body: raw

```json
{
    "telephone": "8-953-983-0807",
    "email": "vladislav.kuznetsovRTN1@yandex.ru"
}
```

RESPONSE:

Status: 200, 30-минутный токен отправлен на почту.

Dsgjkybnm товен сброса:

REQUEST:

POST: /api/user/password/reset

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTYwNzg4MjA2OH0.t9LWRayZbrBvO4WL_u9se3KJthYH9Y6Ns_HtpLp3C8I
```

Body: raw

```json
{
    "reset_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiT0MwNU5UTXRPVGd6TFRBNE1EYz0iLCJzZWNvbmQiOiJ0b3N0ZXIxMjMifSwiZXhwIjoxNjA3ODgzMTA2fQ.w-5ZxVQlTwnI2VQNG94hNODusfkM8ecxcX8tY7gRWXc",
    "telephone" : "8-953-983-0807",
    "email" : "walkmanmail19@gmail.com",
    "password" :"toster12345"
}
```

RESPONSE:

Status: 200, новый access токен.

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTYwNzg4MzQ3OH0.yBx8E-lsvxKVL7jxRKP7Z0coi9yyTR1nXl5OalpnFZk"
}
```

## Обновить аватарку

Обновить:

REQUEST:

POST: /api/user/update/avatar

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTYwNzg4MzQ3OH0.yBx8E-lsvxKVL7jxRKP7Z0coi9yyTR1nXl5OalpnFZk
```

Body: form-data

```
file: <file_bytes>
```

RESPONSE:

Status: 200.

#### Вернуть пользователя:

REQUEST:

POST: /api/user/get

Body: raw

```json
{
    "user_id":66
}
```

RESPONSE:

Status: 200.

```json
{
    "users": [
        {
            "user_id": 66,
            "telephone": "8-953-983-0807",
            "email": "vladislav.kuznetsovRTN1@yandex.ru",
            "name": "Влад Кузнецов",
            "avatar_url": "http://127.0.0.1:4491/download/avatar/id/66",
            "access_token": ""
        }
    ]
}
```