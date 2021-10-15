package tests

import (
	"crate-run-api/app"
	"crate-run-api/models"
	"crate-run-api/tokenstore"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func createTestApp() app.Application {
	instanceId := uuid.New().String()

	db, _ := models.OpenDatabase("file:" + instanceId + "?mode=memory&cache=shared")

	ts := tokenstore.New()
	r := gin.Default()
	gin.SetMode(gin.TestMode)

	app := app.Application{Engine: r, DB: db, TokenStore: ts}
	app.RegisterRoutes()

	return app
}
