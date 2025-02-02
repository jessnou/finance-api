# Проект Finance App

Это проект для управления финансами. Он использует Docker для контейнеризации, а также управляется через `Makefile` для удобного выполнения различных команд, таких как запуск приложения и миграции базы данных.

## Клонирование репозитория

Для начала клонируйте репозиторий:

```bash
git clone https://github.com/jessnou/finance-api.git
cd finance-api
```

Для того чтобы настроить переменные окружения, вам нужно создать файл .env в корне проекта.
Вы можете создать файл вручную или скопировать его из примера:
```
cp .env.example .env
```
Для запуска проекта с использованием Docker и Makefile выполните следующую команду:

```
make run
```
