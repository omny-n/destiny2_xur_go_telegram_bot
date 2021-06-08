# Xur location tbot for Destiny2

A simple telegram bot for showing the location of Xur writen in Golang with [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) and [MongoDB](https://www.mongodb.com)

you can see how the bot works here: https://t.me/Destiny2_Xur_Bot


![](https://github.com/omny-n/destiny2_xur_go_telegram_bot/blob/master/example.jpg)

## Getting started

build docker container:
`docker build -t xur_bot .`


Set environment variables in `docker-compose.yml`:

```
MONGO_INITDB_ROOT_USERNAME: root_username
MONGO_INITDB_ROOT_PASSWORD: root_password
MONGO_INITDB_DATABASE: bot_users

DESTINY_API_KEY: "your Desyiny 2 api key"
TG_API_KEY: "your telegram api key"
DATABASE: "your mongodb link"
```

Run docker-compose: `docker-compose -f docker-compose.yml up`
