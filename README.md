# Financial Chat

A web browser-based chat application written in Golang. Users may create accounts, login, and chat with other users. Users will be able to send commands to the chat bot and the chat bot will send the results of the command back to the chatroom. The application is built using the Gorilla WebSocket and Gin.

## How to run

All the configuration is done through environment variables. You can see, or copy, the .env.example file to get an idea of what variables are needed.

```bash
cp server/.env.example server/.env
```

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
