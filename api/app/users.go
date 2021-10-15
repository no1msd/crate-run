package app

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"crate-run-api/models"

	"github.com/gin-gonic/gin"
)

func (a Application) register(context *gin.Context) {
	credentials := struct {
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
	}{}

	if context.ShouldBindJSON(&credentials) != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	user := models.User{Username: credentials.Username}
	if err := user.Register(a.DB, credentials.Password); err != nil {
		if errors.Is(err, models.ErrUsernameTaken) {
			context.Status(http.StatusConflict)
		} else {
			log.Printf("Handler error: %s", err.Error())
			context.Status(http.StatusInternalServerError)
		}

		return
	}

	token := a.TokenStore.CreateToken(user.Username)

	context.JSON(http.StatusOK, Session{Token: token, User: user})
}

func (a Application) userDetails(context *gin.Context) {
	context.JSON(http.StatusOK, getUser(context))
}

func (a Application) changeNickname(context *gin.Context) {
	nickname, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	user := getUser(context)
	if err := user.ChangeNickname(a.DB, string(nickname)); err != nil {
		log.Printf("Handler error: %s", err.Error())
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusOK)
}
