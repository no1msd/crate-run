// Package app contains the API controllers
package app

import (
	"context"
	"crate-run-api/middlewares"
	"crate-run-api/models"
	"crate-run-api/tokenstore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Application holds the dependencies used by the API controllers
type Application struct {
	Engine     *gin.Engine
	DB         *gorm.DB
	TokenStore *tokenstore.MemoryStore
}

// RegisterRoutes registers all routes served by the API
func (a Application) RegisterRoutes() {
	r := a.Engine
	p := r.Group("").Use(middlewares.AuthRequired(a.DB, a.TokenStore))

	r.GET("/highscores", a.getGlobalHighScores)
	p.GET("/highscores/user", a.getUserHighScores)
	r.GET("/highscores/level/:number", a.getLevelHighScores)
	p.PUT("/highscores/level/:number", a.registerHighScore)
	r.POST("/session", a.login)
	p.DELETE("/session", a.logout)
	r.POST("/user", a.register)
	p.GET("/user", a.userDetails)
	p.PATCH("/user/nickname", a.changeNickname)

	a.registerStaticRoutes()
}

// Run starts the HTTP server with the given listening address and port. It handles SIGINT,
// SIGTERM and SIGHUP gracefully, giving the server a 5 second timeout to shut down.
func (a Application) Run(listenerAddress string, listenerPort uint) {
	srv := &http.Server{
		Addr:    listenerAddress + ":" + strconv.Itoa(int(listenerPort)),
		Handler: a.Engine,
	}

	go func() {
		log.Printf("Server listening on: %s:%d", listenerAddress, listenerPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Cannot listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}

// Session is the data returned by the login and register handlers
type Session struct {
	Token string      `json:"token" binding:"required"`
	User  models.User `json:"user"  binding:"required"`
}

func getUser(context *gin.Context) *models.User {
	userContext, exists := context.Get("user")

	if !exists {
		return nil
	}

	user := userContext.(models.User)
	return &user
}
