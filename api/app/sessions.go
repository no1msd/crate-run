package app

import (
	"log"
	"net/http"

	"crate-run-api/middlewares"
	"crate-run-api/models"

	"github.com/gin-gonic/gin"
)

func (a Application) login(context *gin.Context) {
	credentials := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if context.ShouldBindJSON(&credentials) != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	user := models.User{Username: credentials.Username}
	passwordValid, err := user.CheckPassword(a.DB, credentials.Password)

	if err != nil {
		log.Printf("Handler error: %s", err.Error())
		context.Status(http.StatusInternalServerError)
		return
	}

	if !passwordValid {
		context.Status(http.StatusUnauthorized)
		return
	}

	token := a.TokenStore.CreateToken(credentials.Username)

	context.JSON(http.StatusOK, Session{Token: token, User: user})
}

func (a Application) logout(context *gin.Context) {
	a.TokenStore.InvalidateToken(context.Request.Header[middlewares.AuthTokenHeader][0])

	context.Status(http.StatusOK)
}
