package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (Application) relayEnvVars(context *gin.Context) {
	context.String(http.StatusOK, `
		function getEnvironmentVariables() {
			return {
				basePath: '%s'
			}
		}
	`, os.Getenv("BASE_URL"))
}

func (a Application) registerStaticRoutes() {
	a.Engine.StaticFile("/crate.audio.worklet.js", "/content/crate.audio.worklet.js")
	a.Engine.StaticFile("/crate.js", "/content/crate.js")
	a.Engine.StaticFile("/crate.pck", "/content/crate.pck")
	a.Engine.StaticFile("/crate.png", "/content/crate.png")
	a.Engine.StaticFile("/crate.wasm", "/content/crate.wasm")
	a.Engine.StaticFile("/favicon.png", "/content/favicon.png")
	a.Engine.StaticFile("/pwa-icon.png", "/content/pwa-icon.png")
	a.Engine.StaticFile("/.webmanifest", "/content/.webmanifest")
	a.Engine.StaticFile("/service-worker.js", "/content/service-worker.js")
	a.Engine.StaticFile("/index.html", "/content/index.html")
	a.Engine.StaticFile("/", "/content/index.html")

	a.Engine.GET("/env.js", a.relayEnvVars)
}
