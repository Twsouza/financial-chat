package main

import (
	"log"
	"os"
	"server/application/controller"
	"server/application/repositories"
	"server/application/services"
	"server/data"
	"server/router"
	"strconv"

	"github.com/joho/godotenv"
)

var db data.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("error loading AUTO_MIGRATE_DB env var: %s", os.Getenv("AUTO_MIGRATE_DB"))
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("error loading DEBUG env var: %s", os.Getenv("DEBUG"))
	}

	db.Env = os.Getenv("ENV")
	db.Debug = debug
	db.AutoMigrateDb = autoMigrateDb

	db.Dsn = os.Getenv("DSN")
	db.DsnTest = os.Getenv("DSN_TEST")
}

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	userRep := repositories.NewUserRepository(dbConn)
	userSvc := services.NewUserService(userRep)
	userCtrl := controller.NewUserController(userSvc)

	hub, err := services.NewHub()
	if err != nil {
		log.Fatalf("error creating chat hub: %v", err)
	}
	go hub.Run()

	chatCtrl := controller.NewChatController(hub)

	routes := router.SetupRouter(userCtrl, chatCtrl)
	routes.Run(":8080")
}
