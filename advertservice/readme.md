## Добавить объявление

Тип объявления передавается в поле ad_type:

- ad_type = 1 - Потерявшеся

- ad_type = 2 - Найденное

#### Добавить объявление потерявшегося животного

REQUEST:

POST: /api/advert/user/add

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body: form-data
   
 ```
json: {
          "ad_type":1,
          "animal_type":"Собака",
          "animal_breed":"Овчарка",
          "geo_latitude":50.0,
          "geo_longitude":50.0,
          "comment_text":"Потерялась, помогите найти"
      }

file: <file_bytes>
```

RESPONSE:

Status: 200

 ```json
{
    "ad_owner_id": 66,
    "ad_owner_name": "Влад Кузнецов",
    "ad_type": 1,
    "ad_id": 3,
    "animal_type": "Собака",
    "animal_breed": "Овчарка",
    "geo_latitude": 50,
    "geo_longitude": 50,
    "comment_text": "Потерялась, помогите найти",
    "image_url": "http://127.0.0.1:4491/download/advert/id/3",
    "date_create": "2020-12-13T23:57:09.146475+03:00",
    "date_close": "2020-12-27T23:57:09.146475+03:00"
}
```

#### Добавить объявление найденного животного

REQUEST:

POST: /api/advert/user/add

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05OTktOTk5LTk5OTkiLCJzZWNvbmQiOiJEYW5pbCJ9LCJleHAiOjE4Mjc2Njk4MTN9.egHjkwVQKLwQ5KK3pr6M-8C3YilzrZIMdLFJ87PsBHM
```

Body: form-data
   
 ```
json: {
          "ad_type":2,
          "animal_type":"Собака",
          "animal_breed":"Какая-то порода",
          "geo_latitude":50.00001,
          "geo_longitude":50.00002,
          "comment_text":"Найдена, хозяин, отзовись"
      }

file: <file_bytes>
```

RESPONSE:

Status: 200

 ```json
{
    "ad_owner_id": 68,
    "ad_owner_name": "Danil",
    "ad_type": 2,
    "ad_id": 4,
    "animal_type": "Собака",
    "animal_breed": "Какая-то порода",
    "geo_latitude": 50.00001,
    "geo_longitude": 50.00001,
    "comment_text": "Найдена, хозяин, отзовись",
    "image_url": "http://127.0.0.1:4491/download/advert/id/4",
    "date_create": "2020-12-14T00:24:26.747065+03:00",
    "date_close": "2020-12-17T00:24:26.747065+03:00"
}
```

## Загрузить список объявлений пользователя

Отправленный объявления разделяются на:

- lost - потерявшиеся животные
    
    - list[] - список активных объявлений 
    
    - expire[] - список просрочных объявлений 

- found - найденные животные
    
    - list[] - список активных объявлений 
    
    - expire[] - список просрочных объявлений 


REQUEST:

POST: /api/advert/user/add

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

RESPONSE:

Status: 200

 ```json
{
    "lost": {
        "list": [
            {
                "ad_owner_id": 66,
                "ad_owner_name": "Влад Кузнецов",
                "ad_type": 1,
                "ad_id": 3,
                "animal_type": "Собака",
                "animal_breed": "Овчарка",
                "geo_latitude": 50,
                "geo_longitude": 50,
                "comment_text": "Потерялась, помогите найти",
                "image_url": "http://127.0.0.1:4491/download/advert/id/3",
                "date_create": "2020-12-13T23:57:09.146475+03:00",
                "date_close": "2020-12-27T23:57:09.146475+03:00"
            }
        ],
        "expire": []
    },
    "found": {
        "list": [],
        "expire": []
    }
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 401, не валидные данные пользователя

```json
{
  "error": "Non valid data. "
}
```

Status: 400, не валидные данные, переданные сервису
 
```json
{
  "error": "Non valid data. "
}
```


## Закрыть объявление

Закрытым объявлением считается объявление,

дата закрытия которого равна null

REQUEST:

POST: /api/advert/user/close

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body: raw

```json
{
    "ad_id":3,
    "ad_type":1
}
```

RESPONSE:

Status: 200

 ```json
{
    "ad_id": 3,
    "date_create": null,
    "date_close": null
}
```
#### Загрузить объявения


REQUEST:

POST: /api/advert/user/list

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

RESPONSE:

Status: 200

 ```json
{
    "lost": {
        "list": [],
        "expire": [
            {
                "ad_owner_id": 66,
                "ad_owner_name": "Влад Кузнецов",
                "ad_type": 1,
                "ad_id": 3,
                "animal_type": "Собака",
                "animal_breed": "Овчарка",
                "geo_latitude": 50,
                "geo_longitude": 50,
                "comment_text": "Потерялась, помогите найти",
                "image_url": "http://127.0.0.1:4491/download/advert/id/3",
                "date_create": "2020-12-13T23:57:09.146475+03:00",
                "date_close": null
            }
        ]
    },
    "found": {
        "list": [],
        "expire": []
    }
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 401, не валидные данные пользователя

```json
{
  "error": "Non valid data. "
}
```

Status: 400, не валидные данные, переданные сервису
 
```json
{
  "error": "Non valid data. "
}
```

Status: 400, ошибка обновления статуса объявления
 
```json
{
  "error": "Some problems with data operation. "
}
```


## Продлить объявление

Закрытым объявлением считается объявление,

дата закрытия которого равна null

REQUEST:

POST: /api/advert/user/refresh

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body: raw

```json
{
    "ad_id":3,
    "ad_type":1
}
```

RESPONSE:

Status: 200

 ```json
{
    "ad_id": 3,
    "date_create": "2020-12-14T01:15:23.9562128+03:00",
    "date_close": "2020-12-28T01:15:23.9562128+03:00"
}
```
#### Загрузить объявения

REQUEST:

POST: /api/advert/user/list

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

RESPONSE:

Status: 200

 ```json
{
    "lost": {
        "list": [
            {
                "ad_owner_id": 66,
                "ad_owner_name": "Влад Кузнецов",
                "ad_type": 1,
                "ad_id": 3,
                "animal_type": "Собака",
                "animal_breed": "Овчарка",
                "geo_latitude": 50,
                "geo_longitude": 50,
                "comment_text": "Потерялась, помогите найти",
                "image_url": "http://127.0.0.1:4491/download/advert/id/3",
                "date_create": "2020-12-14T01:15:23.956212+03:00",
                "date_close": "2020-12-28T01:15:23.956212+03:00"
            }
        ],
        "expire": []
    },
    "found": {
        "list": [],
        "expire": []
    }
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 401, не валидные данные пользователя

```json
{
  "error": "Non valid data. "
}
```

Status: 400, не валидные данные, переданные сервису
 
```json
{
  "error": "Non valid data. "
}
```

Status: 400, ошибка обновления статуса объявления
 
```json
{
  "error": "Some problems with data operation. "
}
```

## Обновить объявление

При смене типа объявления, для данного объявления,

буду пересчитанны сроки "годности" в соответсвие с 

тем типом который был выбран.

#### Обновление без смены типа объявлеия

REQUEST:

POST: /api/advert/user/update

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body: form-data
   
 ```
json: {
          "ad_id":3,
          "animal_type":"Кошка",
          "animal_breed":"чб",
          "geo_latitude":50.0003,
          "geo_longitude":50.0003
      }
```

RESPONSE:

Status: 200

 ```json
{
    "ad_owner_id": 66,
    "ad_owner_name": "Влад Кузнецов",
    "ad_type": 1,
    "ad_id": 3,
    "animal_type": "Кошка",
    "animal_breed": "чб",
    "geo_latitude": 50.0003,
    "geo_longitude": 50.0003,
    "comment_text": "Потерялась, помогите найти",
    "image_url": "http://127.0.0.1:4491/download/advert/id/3",
    "date_create": "2020-12-14T01:15:23.956212+03:00",
    "date_close": "2020-12-28T01:15:23.956212+03:00"
}
```

#### Обновление со сменой типа объявлеия

REQUEST:

POST: /api/advert/user/update

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body: form-data
   
 ```
json: {
          "ad_id":3,
          "ad_type":2,
          "animal_type":"Кошка",
          "animal_breed":"Сибирская порода",
          "comment_text": "Найдена, кошка, живет в подвале 63 дома",
      }
```

RESPONSE:

Status: 200

 ```json
{
    "ad_owner_id": 66,
    "ad_owner_name": "Влад Кузнецов",
    "ad_type": 2,
    "ad_id": 3,
    "animal_type": "Кошка",
    "animal_breed": "Сибирская порода",
    "geo_latitude": 50.0003,
    "geo_longitude": 50.0003,
    "comment_text": "Найдена, кошка, живет в подвале 63 дома",
    "image_url": "http://127.0.0.1:4491/download/advert/id/3",
    "date_create": "2020-12-14T01:31:39.595468+03:00",
    "date_close": "2020-12-17T01:31:39.595468+03:00"
}
```

#### Обновление со сменой обложки (файла)


REQUEST:

POST: /api/advert/user/update

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body: form-data
   
 ```
json: {
          "ad_id":3        
      }

file: <file_bytes>
```

RESPONSE:

Status: 200

 ```json
{
    "ad_owner_id": 66,
    "ad_owner_name": "Влад Кузнецов",
    "ad_type": 2,
    "ad_id": 3,
    "animal_type": "Кошка",
    "animal_breed": "Сибирская порода",
    "geo_latitude": 50.0003,
    "geo_longitude": 50.0003,
    "comment_text": "Найдена, кошка, живет в подвале 63 дома",
    "image_url": "http://127.0.0.1:4491/download/advert/id/3",
    "date_create": "2020-12-14T01:42:51.775011+03:00",
    "date_close": "2020-12-14T01:42:51.775011+03:00"
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 401, не валидные данные пользователя

```json
{
  "error": "Non valid data. "
}
```

Status: 400, не валидные данные, переданные сервису
 
```json
{
  "error": "Non valid data. "
}
```

## Загрузить объявления в радиусе от пользователя

ad_owner_id поле - id получателя,

его объявления в ответе не будут отражены.

Отправленный объявления разделяются на:

- lost - потерявшиеся животные
    
    - list[] - список активных объявлений 
    
    - expire[] - список просрочных объявлений 

- found - найденные животные
    
    - list[] - список активных объявлений 
    
    - expire[] - список просрочных объявлений 

REQUEST:

POST: /api/advert/get/in/area

Body: raw
   
```json
{
    "ad_owner_id":66,
    "only_not_closed":true,
    "geo_longitude":50.0,
    "geo_latitude":50.0
}
```

RESPONSE:

Status: 200

 ```json
{
    "lost": {
        "list": [],
        "expire": []
    },
    "found": {
        "list": [
            {
                "ad_owner_id": 68,
                "ad_owner_name": "Danil",
                "ad_type": 2,
                "ad_id": 4,
                "animal_type": "Собака",
                "animal_breed": "Какая-то порода",
                "geo_latitude": 50.00001,
                "geo_longitude": 50.00001,
                "comment_text": "Найдена, хозяин, отзовись",
                "image_url": "http://127.0.0.1:4491/download/advert/id/4",
                "date_create": "2020-12-14T00:24:26.747065+03:00",
                "date_close": "2020-12-17T00:24:26.747065+03:00"
            }
        ],
        "expire": []
    }
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 401, не валидные данные пользователя

```json
{
  "error": "Non valid data. "
}
```

Status: 400, не валидные данные, переданные сервису
 
```json
{
  "error": "Non valid data. "
}
```