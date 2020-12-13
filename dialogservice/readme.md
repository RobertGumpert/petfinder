## Создать диалог

Имеется два пользователя:

- id 66, Влад Кузнецов
- id 68, Danil

#### 'Влад Кузнецов', создаёт диалог с 'Danil', передав в Header свой access токен, и в json информацию о 'Danil'

REQUEST:

POST: /api/user/dialog/create

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body raw:
   
 ```json
{
    "user_id": 68,
    "name": "Danil"
}
```

RESPONSE:

Status: 200, получен id нового диалога.

 ```json
{
    "id": 14
}
```

#### 'Danil', попытается создать с 'Влад Кузнецов' новый диалог

REQUEST:

POST: /api/user/dialog/create

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05OTktOTk5LTk5OTkiLCJzZWNvbmQiOiJEYW5pbCJ9LCJleHAiOjE4Mjg2MzU2Nzl9.4gAn3nW7IjeHUVgeGocVmmcYZet8DjMWUg87co9anSM
```

Body raw:
   
 ```json
{
    "user_id": 66,
    "name": "Влад Кузнецов"
}
```

RESPONSE:

Status: 400, ошибка создания существующего диалога.

 ```json
{
    "error": "Some problems with data operation. "
}
```

Status: 400, неправильноый json.

```json
{
    "error": "Non valid data. "
}
```

Status: 401, не авторизованный польз./сгорел access токен.


## Отправить сообщение

#### 'Влад Кузнецов', отправлет сообщение 'Danil'
REQUEST:

POST: /api/user/message/send

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body raw:
   
 ```json
{
    "dialog_id" : 14,
    "text":"Короче собираюсь"
}
```

RESPONSE:

Status: 200, TO DO: "user_id": 0, - убрать поле.

 ```json
{
    "dialog_id": 14,
    "user_receiver_id": 66,
    "message_id": 54,
    "user_id": 0,
    "user_name": "Влад Кузнецов",
    "text": "Короче собираюсь",
    "date_create": "2020-12-13T22:04:48.624158+03:00"
}
```

#### 'Danil', отправлет сообщение 'Влад Кузнецов'

REQUEST:

POST: /api/user/message/send

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05OTktOTk5LTk5OTkiLCJzZWNvbmQiOiJEYW5pbCJ9LCJleHAiOjE4Mjg2MzU2Nzl9.4gAn3nW7IjeHUVgeGocVmmcYZet8DjMWUg87co9anSM
```

Body raw:
   
 ```json
{
    "dialog_id" : 14,
    "text":"давай в 11 10 там"
}
```

RESPONSE:

Status: 200, TO DO: "user_id": 0, - убрать поле.

 ```json
{
    "dialog_id": 14,
    "user_receiver_id": 68,
    "message_id": 50,
    "user_id": 0,
    "user_name": "Danil",
    "text": "давай в 11 10 там",
    "date_create": "2020-12-13T22:00:49.050434+03:00"
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 400, неправильноый json.

```json
{
    "error": "Non valid data. "
}
```

Status: 400, ошибка записи в таблицу сообщений.

```json
{
    "error": "Some problems with data operation. "
}
```

## Загрузить диалоги

Для каждого диалога, сервер формирует поле skip_messages,
 
значение которого надо передать, для того чтобы получить

следующие 20 сообщений в диалоге.

При первой загрузке, skip_messages всегда равно 0.

Поле dialog_name содержит имя диалога, которое отобразит у пользоваеля.

#### 'Влад Кузнецов' загружает диалоги.

REQUEST:

GET: /api/user/dialog/get

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

RESPONSE:

Status: 200, получен список диалогов и сообщений.

 ```json
{
    "dialogs": [
        {
            "dialog_id": 14,
            "dialog_name": "Danil",
            "user_receiver_id": 66,
            "skip_messages": 0,
            "messages": [
                {
                    "message_id": 54,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Короче собираюсь",
                    "date_create": "2020-12-13T22:04:48.624158+03:00"
                },
                {
                    "message_id": 53,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Поэтому и спрашиваю",
                    "date_create": "2020-12-13T22:04:38.395116+03:00"
                },
                {
                    "message_id": 52,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Только что покушал",
                    "date_create": "2020-12-13T22:04:30.190108+03:00"
                },
                {
                    "message_id": 51,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Я вообще в трусах сижу",
                    "date_create": "2020-12-13T22:04:19.362091+03:00"
                },
                {
                    "message_id": 50,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "давай в 11 10 там",
                    "date_create": "2020-12-13T22:00:49.050434+03:00"
                },
                {
                    "message_id": 49,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Я еще не одевался",
                    "date_create": "2020-12-13T22:00:39.612707+03:00"
                },
                {
                    "message_id": 48,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Погоди",
                    "date_create": "2020-12-13T22:00:32.694187+03:00"
                },
                {
                    "message_id": 47,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Мне собираться?",
                    "date_create": "2020-12-13T22:00:08.86826+03:00"
                },
                {
                    "message_id": 46,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "занят щас",
                    "date_create": "2020-12-13T21:59:39.315027+03:00"
                },
                {
                    "message_id": 45,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Потом отпишу",
                    "date_create": "2020-12-13T21:59:30.542242+03:00"
                },
                {
                    "message_id": 44,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "не",
                    "date_create": "2020-12-13T21:59:21.792244+03:00"
                },
                {
                    "message_id": 43,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Сдал?",
                    "date_create": "2020-12-13T21:58:59.027711+03:00"
                },
                {
                    "message_id": 42,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Влад мытая подмышка",
                    "date_create": "2020-12-13T21:58:32.354381+03:00"
                },
                {
                    "message_id": 41,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Ну ок",
                    "date_create": "2020-12-13T21:57:47.713692+03:00"
                },
                {
                    "message_id": 40,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Ля я только помылся",
                    "date_create": "2020-12-13T21:57:36.06567+03:00"
                },
                {
                    "message_id": 39,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "У энтерры?",
                    "date_create": "2020-12-13T21:56:47.756205+03:00"
                },
                {
                    "message_id": 38,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "в 11 00?",
                    "date_create": "2020-12-13T21:56:36.898897+03:00"
                },
                {
                    "message_id": 37,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "я пока ИИС делаю",
                    "date_create": "2020-12-13T21:56:28.729797+03:00"
                },
                {
                    "message_id": 36,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "хз",
                    "date_create": "2020-12-13T21:56:21.757691+03:00"
                },
                {
                    "message_id": 35,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Во сколько встречаемся?",
                    "date_create": "2020-12-13T21:55:51.260809+03:00"
                }
            ]
        }
    ]
}
```

#### 'Danil' загружает свои диалоги

REQUEST:

GET: /api/user/dialog/get

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05OTktOTk5LTk5OTkiLCJzZWNvbmQiOiJEYW5pbCJ9LCJleHAiOjE4Mjc2Njk4MTN9.egHjkwVQKLwQ5KK3pr6M-8C3YilzrZIMdLFJ87PsBHM
```

RESPONSE:

Status: 200, получен список диалогов и сообщений.

 ```json
{
    "dialogs": [
        {
            "dialog_id": 14,
            "dialog_name": "Влад Кузнецов",
            "user_receiver_id": 68,
            "skip_messages": 0,
            "messages": [
                {
                    "message_id": 54,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Короче собираюсь",
                    "date_create": "2020-12-13T22:04:48.624158+03:00"
                },
                {
                    "message_id": 53,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Поэтому и спрашиваю",
                    "date_create": "2020-12-13T22:04:38.395116+03:00"
                },
                {
                    "message_id": 52,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Только что покушал",
                    "date_create": "2020-12-13T22:04:30.190108+03:00"
                },
                {
                    "message_id": 51,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Я вообще в трусах сижу",
                    "date_create": "2020-12-13T22:04:19.362091+03:00"
                },
                {
                    "message_id": 50,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "давай в 11 10 там",
                    "date_create": "2020-12-13T22:00:49.050434+03:00"
                },
                {
                    "message_id": 49,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Я еще не одевался",
                    "date_create": "2020-12-13T22:00:39.612707+03:00"
                },
                {
                    "message_id": 48,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Погоди",
                    "date_create": "2020-12-13T22:00:32.694187+03:00"
                },
                {
                    "message_id": 47,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Мне собираться?",
                    "date_create": "2020-12-13T22:00:08.86826+03:00"
                },
                {
                    "message_id": 46,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "занят щас",
                    "date_create": "2020-12-13T21:59:39.315027+03:00"
                },
                {
                    "message_id": 45,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Потом отпишу",
                    "date_create": "2020-12-13T21:59:30.542242+03:00"
                },
                {
                    "message_id": 44,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "не",
                    "date_create": "2020-12-13T21:59:21.792244+03:00"
                },
                {
                    "message_id": 43,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Сдал?",
                    "date_create": "2020-12-13T21:58:59.027711+03:00"
                },
                {
                    "message_id": 42,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "Влад мытая подмышка",
                    "date_create": "2020-12-13T21:58:32.354381+03:00"
                },
                {
                    "message_id": 41,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Ну ок",
                    "date_create": "2020-12-13T21:57:47.713692+03:00"
                },
                {
                    "message_id": 40,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Ля я только помылся",
                    "date_create": "2020-12-13T21:57:36.06567+03:00"
                },
                {
                    "message_id": 39,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "У энтерры?",
                    "date_create": "2020-12-13T21:56:47.756205+03:00"
                },
                {
                    "message_id": 38,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "в 11 00?",
                    "date_create": "2020-12-13T21:56:36.898897+03:00"
                },
                {
                    "message_id": 37,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "я пока ИИС делаю",
                    "date_create": "2020-12-13T21:56:28.729797+03:00"
                },
                {
                    "message_id": 36,
                    "dialog_id": 14,
                    "user_id": 68,
                    "user_name": "Danil",
                    "text": "хз",
                    "date_create": "2020-12-13T21:56:21.757691+03:00"
                },
                {
                    "message_id": 35,
                    "dialog_id": 14,
                    "user_id": 66,
                    "user_name": "Влад Кузнецов",
                    "text": "Во сколько встречаемся?",
                    "date_create": "2020-12-13T21:55:51.260809+03:00"
                }
            ]
        }
    ]
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 400, ошибка чтения из бд.

```json
{
    "error" : "Some problems with data operation. "
}
```

## Загрузить следующие сообщения из диалога

Для загрузки следующих 20 сообщений,

необходимо передать поле last_skip.

Сервер вернет сообщения и следующее значение

next_skip, которое надо передать в качестве last_skip,

для получения очередной партии из 20 собщений.

#### 'Влад Кузнецов' загружает следующие 20 сообщений.

REQUEST:

POST: /api/user/message/batching/next

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body, raw

```json
{
    "dialog_id" : 14,
    "last_skip":0
}
```

RESPONSE:

Status: 200, получен список сообщений.

```json
{
    "dialog_id": 14,
    "next_skip": 20,
    "user_receiver_id": 66,
    "messages": [
        {
            "message_id": 34,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "mb",
            "date_create": "2020-12-13T21:55:28.544242+03:00"
        },
        {
            "message_id": 33,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "da",
            "date_create": "2020-12-13T21:55:20.720985+03:00"
        },
        {
            "message_id": 32,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Ты то едешь?",
            "date_create": "2020-12-13T21:54:52.835343+03:00"
        },
        {
            "message_id": 31,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Закрывать по мии",
            "date_create": "2020-12-13T21:54:43.305742+03:00"
        },
        {
            "message_id": 30,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Мне 5,6",
            "date_create": "2020-12-13T21:54:34.465647+03:00"
        },
        {
            "message_id": 29,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "kk",
            "date_create": "2020-12-13T21:54:09.704135+03:00"
        },
        {
            "message_id": 28,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Еду",
            "date_create": "2020-12-13T21:53:19.011566+03:00"
        },
        {
            "message_id": 27,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Мы же вчера говорили",
            "date_create": "2020-12-13T21:53:11.249427+03:00"
        },
        {
            "message_id": 26,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Ты шо",
            "date_create": "2020-12-13T21:53:03.409186+03:00"
        },
        {
            "message_id": 25,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Дядь",
            "date_create": "2020-12-13T21:52:55.202224+03:00"
        },
        {
            "message_id": 24,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "Т.е. в уник сегодня не едешь?",
            "date_create": "2020-12-13T21:52:31.240912+03:00"
        },
        {
            "message_id": 23,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "Ты мии закрыл?",
            "date_create": "2020-12-13T21:52:20.999606+03:00"
        },
        {
            "message_id": 22,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "Сорри",
            "date_create": "2020-12-13T21:52:14.167923+03:00"
        },
        {
            "message_id": 21,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "Нечаянно набрал",
            "date_create": "2020-12-13T21:52:01.781521+03:00"
        },
        {
            "message_id": 20,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "ponyal",
            "date_create": "2020-12-13T21:51:52.262153+03:00"
        },
        {
            "message_id": 19,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "И быстрее было на строках сделать",
            "date_create": "2020-12-13T21:51:19.010563+03:00"
        },
        {
            "message_id": 18,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Я же говорил Что мне в падлу разбираться",
            "date_create": "2020-12-13T21:51:08.064649+03:00"
        },
        {
            "message_id": 17,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Я хотел быстро сделать",
            "date_create": "2020-12-13T21:50:53.507967+03:00"
        },
        {
            "message_id": 16,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Похуй",
            "date_create": "2020-12-13T21:50:45.788135+03:00"
        },
        {
            "message_id": 15,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Мне было",
            "date_create": "2020-12-13T21:50:38.49751+03:00"
        }
    ]
}
```


#### 'Влад Кузнецов' загружает следующие 20 сообщений.

REQUEST:

POST: /api/user/message/batching/next

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body, raw

```json
{
    "dialog_id" : 14,
    "last_skip":20
}
```

RESPONSE:

Status: 200, получен список сообщений.

```json
{
    "dialog_id": 14,
    "next_skip": 40,
    "user_receiver_id": 66,
    "messages": [
        {
            "message_id": 14,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Ибо",
            "date_create": "2020-12-13T21:50:31.309524+03:00"
        },
        {
            "message_id": 13,
            "dialog_id": 14,
            "user_id": 66,
            "user_name": "Влад Кузнецов",
            "text": "Ну я говорил только что не знаю есть ли на шарп",
            "date_create": "2020-12-13T21:50:23.561383+03:00"
        },
        {
            "message_id": 12,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "Простейшая битовая коллекция bool[]",
            "date_create": "2020-12-13T21:47:28.941552+03:00"
        },
        {
            "message_id": 11,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "типо ты строками делал",
            "date_create": "2020-12-13T21:47:11.920255+03:00"
        },
        {
            "message_id": 10,
            "dialog_id": 14,
            "user_id": 68,
            "user_name": "Danil",
            "text": "2 лаба",
            "date_create": "2020-12-13T21:47:01.885231+03:00"
        }
    ]
}
```

#### 'Влад Кузнецов' загрузил все сообщения.

REQUEST:

POST: /api/user/message/batching/next

HEADER:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjp7ImZpcnN0IjoiOC05NTMtOTgzLTA4MDciLCJzZWNvbmQiOiLQktC70LDQtCDQmtGD0LfQvdC10YbQvtCyIn0sImV4cCI6MTgyODYzMzY3OH0.VefsviYJ_0lRBlzcK_Hj9XhXk5-Tq40Omkmnja6vQ8U
```

Body, raw

```json
{
    "dialog_id" : 14,
    "last_skip":40
}
```

RESPONSE:

Status: 200, получен пустой список сообщений.

```json
{
    "dialog_id": 14,
    "next_skip": 60,
    "user_receiver_id": 66,
    "messages": []
}
```

Status: 401, не авторизованный польз./сгорел access токен.

Status: 400, ошибка чтения из бд.

```json
{
    "error" : "Some problems with data operation. "
}
```