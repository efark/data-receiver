package webserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Routers return the http.Handler with the basic handlers for the data endpoint.
func Routers() http.Handler {
	e := gin.Default()

	e.GET("/health", HealthHandler)
	e.POST("/data/:service", DataHandler)

	return e
}
