.PHONY: startServer keepContainerOpen

startServer:
	cd server && go run cmd/server.go

keepContainerOpen:
	tail -f /dev/null
