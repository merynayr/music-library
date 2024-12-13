# music-library

## Реализация онлайн библиотеки песен

**Задание:**
Необходимо реализовать следующее:

1. Выставить rest методы:

Получение данных библиотеки с фильтрацией по всем полям и пагинацией;

Получение текста песни с пагинацией по куплетам;

Удаление песни;

Изменение данных песни;

Добавление новой песни в формате
      JSON
{
 "group": "Muse",
 "song": "Supermassive Black Hole"
};

2. При добавлении сделать запрос в АПИ, описанного сваггером
3. Обогащенную информацию положить в БД postgres (структура БД должна быть создана путем миграций при старте сервиса)
4. Покрыть код debug- и info-логами
5. Вынести конфигурационные данные в .env-файл
6. Сгенерировать сваггер на реализованное АПИ


## Запуск
1. Запустить Docker Desktop
2. Прописать в CLI: docker compose up -d

## Swagger
Swagger с описанием API доступен после запуска сервиса по адресу: http://localhost:8088/swagger/index.html
