// Package middlewares contains the middlewares used in the application
package middlewares

import (
	"log"
	"net/http"

	"crate-run-api/models"
	"crate-run-api/tokenstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthTokenHeader is the HTTP header key where the AuthRequired middleware looks for the
// authorization token when handling requests.
const AuthTokenHeader string = "X-Auth-Token"

// AuthRequired is a middleware that checks the request header for an access token and validates
// it with a token store. If the token is invalid it aborts the request with HTTP Unauthorized.
// If it's valid it sets the current user's data in the context in the "user" key.
func AuthRequired(db *gorm.DB, tokenStore *tokenstore.MemoryStore) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header[AuthTokenHeader]

		if len(token) == 0 || !tokenStore.IsTokenValid(token[0]) {
			context.AbortWithStatus(http.StatusUnauthorized)
		} else {
			username, _ := tokenStore.GetUsernameForToken(token[0])

			var user models.User
			err := db.Where(models.User{Username: username}).First(&user).Error
			if err != nil {
				log.Printf("AuthRequired error: %s", err.Error())
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			context.Set("user", user)
			context.Next()
		}
	}
}
