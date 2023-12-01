Disclaimer:

The objective of this test was to implement a backend system; the front end was not the primary focus but was created as a means to validate and facilitate the use of the backend.

# Financial Chat

A web browser-based chat application written in Golang. Users may create accounts, login, and chat with other users. Users will be able to send commands to the chatbot and the chatbot will send the results of the command back to the chatroom. The application is built using the Gorilla WebSocket and Gin.

## How to run

All the configuration is done through environment variables. You can see, or copy, the .env.example file to get an idea of what variables are needed.

```bash
cp server/.env.example server/.env
cp bot/.env.example bot/.env
cp client/.env.example client/.env
```

(See Makefile `setEnvServerAndBot` target for more details)

Use docker-compose to start the project with a single command:

```bash
docker-compose up -d
```

## Services

### Server

The server will handle signups and logins from the users. Also, it will handle the chat rooms, and send the commands to the bot (which will process and send it back to the chat room).

#### Routes

- `/signup` - POST - Create a new user
  - Payload:
    - `username` - string
    - `password` - string
    - `email` - string
- `/login` - POST - Login with an existing user
  - Payload:
    - `username` - string
    - `password` - string

### Queue

The queue service will have two queues: `message` and `commands`. The message queue will send a message to the chat room, and the commands queue will send commands to the bot.

### Bot

The bot will process the commands queue and send the results back to the chat room. The bot services use a worker pool, the environment variable starts the `CONCURRENCY_WORKERS` number of workers. The bot will process the commands in parallel. The bot will send the results to the message queue.

#### Available commands

- `/stock=stock_code` - Get the current price of the stock, and data from https://stooq.com, and send a message to the queue with the stock quotes, aka `stock-bot`.

![Screenshot from 2023-10-07 22-47-40](https://github.com/Twsouza/financial-chat/assets/8239709/ed1326fe-cdc4-4c35-af07-ecc6e33aa651)
