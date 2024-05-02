FROM alpine:latest

# Установка SQLite
RUN apk --no-cache add sqlite

# Создание рабочей директории
WORKDIR /app

# Копирование файла базы данных
COPY ./storage/sso.db /app/storage/sso.db

# Запуск SQLite с файлом базы данных
CMD ["sqlite3", "/app/storage/sso.db"]
