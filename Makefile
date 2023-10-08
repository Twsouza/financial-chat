.PHONY: startServer startBot keepContainerOpen

startServer:
	cd server && go run cmd/server.go

startBot:
	cd bot && go run cmd/server.go

keepContainerOpen:
	tail -f /dev/null

setEnvServerAndBot:
	cp server/.env.example server/.env
	cp bot/.env.example bot/.env
