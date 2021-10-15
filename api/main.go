package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"crate-run-api/app"
	"crate-run-api/models"
	"crate-run-api/tokenstore"
	"crate-run-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	arguments, err := utils.ParseArguments()
	if err != nil {
		log.Fatal(err.Error())
	}

	configFile, err := os.OpenFile(arguments.ConfigPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()

	config, err := utils.LoadOrCreateConfig(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := models.OpenDatabase(config.Database.Path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	if !arguments.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	app := app.Application{
		Engine:     gin.Default(),
		TokenStore: tokenstore.New(),
		DB:         db,
	}
	app.RegisterRoutes()
	app.Run(config.Listener.Address, config.Listener.Port)
}
