# Тестовое задание на позицию стажера-бекендера
# Сервис динамического сегментирования пользователей

### Проблема:

В Авито часто проводятся различные эксперименты — тесты новых продуктов, тесты интерфейса, скидочные и многие другие.
На архитектурном комитете приняли решение централизовать работу с проводимыми экспериментами и вынести этот функционал в отдельный сервис.

### Задача:

Требуется реализовать сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)

### Решение:

В ходе работы, был создан сервис, позволяющий выполнять следующие операции:

1. Метод создания сегмента. Принимает название сегмента.
2. Метод удаления сегмента. Принимает название сегмента.
3. Метод добавления пользователя в сегмент. Принимает список названий сегментов, которые нужно добавить пользователю и список названий сегментов, которые нужно удалить у пользователя, а также id пользователя.
4. Метод получения активных сегментов пользователя. Принимает на вход id пользователя.

А также решены все опциональные задачи:

1. Реализован механизм сохранения истории (попадания/выбывания) пользователя из сегмента и возможность получения отчета по пользователю за определенный период. На вход: месяц и год. На выходе ссылка на CSV файл.
2. Реализована возможность задавать TTL (время автоматического удаления пользователя из сегмента).
3. Автоматизирована задача по добавлению заданного процента пользователей в создаваемый сегмент.

При разработке данного сервиса также соблюдены необходимые технические требования, а именно:
1. Сервис предоставляет HTTP API с форматом JSON как при отправке запроса, так и при получении результата.
2. Язык разработки: Golang.
3. Реляционная СУБД: PostgreSQL.
4. Покрытие кода unit тестами
5. Наличие Swagger документации
6. Использование docker и docker-compose для поднятия и развертывания dev-среды.
7. Весь код выложен на Github с Readme файлом и с инструкцией по запуску и примерами взаимодействия с API.
8. Все возникшие вопросы и варианты решения оставлены в списке, который расположен конце данного README файла.

### Используемые библиотеки и технологии при разработке
* [httprouter](https://github.com/julienschmidt/httprouter) - Httprouter by Julienschmidt
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [cleanenv](https://github.com/ilyakaznacheev/cleanenv) - Go configuration
* [logrus](https://github.com/sirupsen/logrus) - Logger
* [swag](https://github.com/swaggo/swag) - Swagger
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [mockery](https://github.com/vektra/mockery) - Mocking framework
* [csvutil](https://github.com/jszwec/csvutil) - Csv utility
* [Docker](https://www.docker.com/) - Docker
* [cron](https://en.wikipedia.org/wiki/Cron) - Cron job scheduler

## Как использовать?

- Запуск тестов производится командой `make test`
- Запуск тестов с отчетом о покрытии в формате html `make cover`
- Запуск сервиса производится с помощью выполнения команды `make compose-up`
- По окончанию работы с сервисом, остановить и удалить созданные контейнеры можно командой `make compose-down`

Для выполенения последующих запросов предлагается использовать приведенные ниже Curl запросы или 
на странице Swagger документации, доступной по ссылке: http://127.0.0.1:8080/swagger/index.html/

Более подробная работа с API приведена ниже

### Примеры использования API сервиса:
Некоторые примеры запросов
- [Создать сегмент](#create-segment)
- [Удалить сегмент](#delete-segment)
- [Добавить сегмент пользователю. Удалить сегмент пользователю. Задать TTL для сегмента](#edit-segment-to-user)
- [Получить активные сегменты у пользователя](#get-active-segments-from-user)
- [Получить отчет в формате CSV по сегментам пользователя](#segment-report-link)
- [Автоматическое удаление пользователя из сегмента (TTl)](#segment-ttl)

### Процесс создания сегмента <a name="create-segment"></a>

***Внимание***: Опциональное задание №3: " ***В методе создания сегмента***, добавить опцию указания процента пользователей, которые будут попадать в сегмент автоматически. 
В методе получения сегментов пользователя, добавленный сегмент должен отдаваться у заданного процента пользователей."
Таким образом, расширяем уже существующий функционал ***опциональным*** полем **"percent"**. 
- В случае отсутствия необходимости автоматически добавлять сегмент заданному проценту пользователей - поле опускается.
- В случае, когда при заданном проценте, количество пользователей, входящих в этот процент будет <1 - сегмент добавится 1му случайному пользователю. 

Пример создания сегмента:
```curl
curl --location --request POST 'http://localhost:8080/api/segments/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"discount90",
    "percent":80
}'
```
Пример ответа:
```json
{
  "description":"Content created",
  "developer_msg":"201 Created",
  "code":"Avito_Segment_Service-000201"
}
```
### Процесс удаления сегмента <a name="delete-segment"></a>

Пример удаление сегмента:
```curl
curl --location --request DELETE 'http://localhost:8080/api/segments/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"discount90"
}'
```
Пример ответа:
```json
{
    "description": "Content deleted",
    "developer_msg": "200 OK",
    "code": "Avito_Segment_Service-000204"
}
```

### Процесс добавления/удаления сегмента пользователю <a name="edit-segment-to-user"></a>

***Внимание***: Опциональное задание №2: "Для реализации TTL В ***метод добавления сегментов*** пользователю передаём время удаления пользователя из сегмента отдельным полем".
Таким образом, расширяем уже существующий функционал ***опциональным*** полем **"ttl"**.

- В случае отсутствия необходимости добавлять пользователю время действия(жизни) сегмента - поле опускается.
- В случае отсутствия необходимости удалять пользователю сегмент - поле опускается.

Пример добавления/ удаления сегментов пользователю/ у пользователя.
```curl
curl --location --request PUT 'http://localhost:8080/api/segments/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "userID":"2",
    "add":[
        {
            "name":"discount50",
            "ttl_days":3
        },
        {
            "name":"discount30" 
        }
    ],
    "delete":["discount80"]
}'
```
Пример ответа:
```json
{
    "description": "Content updated",
    "developer_msg": "200 OK",
    "code": "Avito_Segment_Service-000200"
}
```

### Процесс получения активных сегментов у пользователя <a name="get-active-segments-from-user"></a>

Пример получения активных сегментов у пользователя с id=2
```curl 
curl --location --request GET 'http://localhost:8080/api/segments/2' \
--header 'Content-Type: application/json'
```

Пример ответа:
```json
[
    {
        "ID": "12",
        "name": "discount30"
    },
    {
        "ID": "13",
        "name": "discount50"
    }
]
```

### Процесс получения отчета в формате CSV о истории взаимодействия сегментов и пользователя <a name="segment-report-link"></a>

В данном задании представленно 2 варианта отчета. 

Пример получения CSV очета образца №1 (согласно ТЗ) по ID пользователя и query параметрам month,year.
```curl
curl --location --request GET 'http://127.0.0.1:8080/api/reports/download/2?month=august&year=2023' \
  --header 'accept: application/octet-stream'
```

Пример отчета:

| UserID | SegmentName | Action  | Date                        |
|--------|-------------|---------|-----------------------------|
| 2      | discount50  | created | 2023-08-30T15:27:24.633856Z |
| 2      | discount30  | created | 2023-08-30T15:27:24.637647Z |
| 2      | discount80  | created | 2023-08-30T19:07:15.915439Z |
| 2      | discount90  | created | 2023-08-30T19:07:37.009223Z |
| 2      | discount100 | created | 2023-08-30T19:07:58.741777Z |
| 2      | discount50  | deleted | 2023-08-30T19:09:03.353874Z |

Проблема данного отчета состоит в том, что отслеживание добавления сегмента пользователю и его удаление достаточно затруднено с точки зрения клиента, для которого и составляется отчет.
В данном случае между созданием сегмента ***discount50*** и его удалением **4** записи. Однако, при активном взаимодействии сегментов и пользователей это число может сильно возрасти.
В таком случае, искать когда сегмент был удален крайне неудобно. 

В связи с этим недостатком было принято решение - создать другой формат отчета. Новый формат не должен терять информативность и при этом увеличить качество восприятия информации.

Пример получения CSV очета образца №2 по ID пользователя и query параметрам month,year.
```curl
curl --location --request GET 'http://127.0.0.1:8080/api/reports/optimized/download/2?month=august&year=2023' \
  --header 'accept: application/octet-stream'
```

Пример отчета:

| UserID | SegmentName | Active | CreatedAt                   | DeletedAt                   |
|--------|-------------|--------|-----------------------------|-----------------------------|
| 2      | discount50  | false  | 2023-08-30T15:27:24.633856Z | 2023-08-30T19:09:03.353874Z |
| 2      | discount30  | true   | 2023-08-30T15:27:24.637647Z |                             |
| 2      | discount80  | true   | 2023-08-30T19:07:15.915439Z |                             |
| 2      | discount90  | true   | 2023-08-30T19:07:37.009223Z |                             |
| 2      | discount100 | true   | 2023-08-30T19:07:58.741777Z |                             |

В данном отчете можно наблюдать в каком состоянии сейчас находится сегмент. 
Активен -true, иначе false, также дата создания и удаления сегмента у пользователя расположена рядом, что позволяет не искать желаемый сегмент по всему отчету.

### Процесс автоматического удаления пользователя из сегмента (TTl) <a name="segment-ttl"></a>

Для реализации данного функционала было принято решение использовать Cron - инструмент для планирования и выполнения задач на Unix системах.
Для этого был создан bash скрипт cr.sh, содержащий следующую команду:

```curl
curl -X POST 127.0.0.1:8080/api/segments/ttl
```

Данный запрос выполняет работу по автоматическому удалению пользователя из сегмента по истечению времени жизни сегмента. 
Для реализации выполнения данного скрипта по расписанию - описаны следующие строчки в Dockerfile (более подробно см. полный Dockerfile):
```dockerfile
COPY cr.sh /cr.sh
RUN echo '* * * * * bash /cr.sh' >> /etc/crontabs/root
CMD usr/sbin/crond && /segment-service
```
Согласно документации cron сочетание ```0 0 * * *``` будет осуществлять запуск скрипта ежедневно в 00:00, что решает проблему автоматического контроля времени жизни сегментов у пользователей. 
В настоящем случае указаны ```* * * * *``` - активация скрипта ежеминутно. Данная реализаци служит лишь для демонстрации работоспособности данной функции сервиса и экономии времени проверяющим :)

## Возникшие вопросы и выбранные решения

1. Согласно основному заданию: "Метод добавления пользователя в сегмент. 
Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, 
id пользователя". Хорошая ли это идея?
>Slug - это уникальная строка идентификатор, понятная человеку. В случае URL передача списка названий как slug список - странная затея. 
То есть запрос на удаление имел бы примерный вид /api/segments/slug1/slug2/slu3/slug4 и.т.д. При этом по ТЗ все запросы требуется выполнять в JSON формате.
Мною было выбрано решение отказаться от паровозика названий в url и использовать JSON в теле запроса для передачи данных. 
Такой способ передачи наиболее удобен и легко поддается изменению в случае необходимости. 
2. Для формирования CSV отчета необходимо подавать на вход месяц и год как параметры, однако о промежутке ничего не сказано. Какой выбрать промежуток времени? 
>Я решил, что входящие месяц и год будут стартовой точкой фомирования отчета до времени настоящего запроса.
3. Достаточно ли удобен пример отчета? Как можно лучше сформировать отчет?
>Пример требуемого отчета мне показался неудобным и в примерах запроса на скачивания CSV отчета я решил добавить новый формат, но требуемый также сохранил.
4. Какое должно быть поведение если при создании сегмента я задаю такой процент, при котором в выборку попадет меньше 1 пользователя.
>Мне показалось правильным в таком случае брать 1 пользователя, которому случайно достанется сегмент. Данная логика легко поддается изменению в случае промаха по поведению. 
5. Если сегмент уже был добавлен пользователям, а затем сам сегмент удален из таблицы сегментов, то какое должно быть поведение сервиса?
>В таком случае я посчитал, что сегмент должен быть удален у пользователей тоже. Допустим пользователю случайно выдали сегмент "скидка100", а не "скидка10".
Тогда при удалении сегмена "скидка100" из таблицы, пользователь не сможет воспользоваться этим сегментом.
6. Если сервис существует отдельно, то как реализовать взаимодействие с сущностью пользователя?
>Так как сервис существует для взаимодействия с сегментами и никак не изменяет поведение пользователей было принято решение считать, что база данных пользователей
уже существует. В моем случае для тестирования функционала сегментов, я храню 5 пользователей в отдельной таблице, которая заполняется при инициализации
схемы базы данных.
7. Как реализовать удаление сегмента так, чтобы не терять важную информацию, например, историю взаимодействия пользователя и сегмента?
>Удаление сегмента выполнено по принципу Soft Delete, то есть существует дополнительное поле, содержащее информацию об активности сегмента.
Такой способ позволяет не терять важную информацию при взаимодействии с сегментами. 
Например, в случае добавления сегмента, который ранее был удален - изменится поле активности этого сегмента, а новая запись не будет создана.
Дополнительная польза такого подхода в том, что любое удаление обратимо.