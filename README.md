# Телеграм бот для получения расписания Череповецкого Государственного Университета

Адрес бота: [@schedulechsubot](https://t.me/schedulechsubot)

## Возможности бота
+ Получение расписания на любой день
+ Запоминание вашей группы
+ Получение расписания на выбранный временной диапазон
+ Ускоренное получение расписание на ближайшие два дня

## Использование

Начните диалог с ботом, перейдя по [ссылке](https://t.me/schedulechsubot)

Общение с ботом происходит за счёт использования клавиатур, которые он вам предоставит во время использования

## Технологии

### Бот
+ Для взаимодействия с Telegram используется [echotron](https://github.com/NicoNex/echotron)
+ Для логирования используется [logrus](https://github.com/sirupsen/logrus)

### Получение расписания
+ Для отправки запросов к серверу ЧГУ используется net/http

### Получение переменных окружения 
+ Для получения переменных среды используется [cleanenv](https://github.com/ilyakaznacheev/cleanenv)

### Хранение данных пользователей(id-пользователя, id-группы)
+ В качестве базы данных используется postgresql c драйвером [pgx](https://github.com/jackc/pgx)

## Установка и запуск

Клонирование репозитория
~~~shell
git clone https://github.com/BobaUbisoft17/chsuBot
~~~

### Добавление переменных среды

Необходимо создать файл .env, затем внести переменные среды
~~~shell
user=пользователь базы данных
password=пароль для подключения к базе данных
DBName=название базы данных
ADMIN=id пользователя с правами администратора
BOTTOKEN=токен бота
DATABASEURL=URL базы данных
~~~



### Запуск
Выполните команду:
~~~shell
docker-compose up --build
~~~