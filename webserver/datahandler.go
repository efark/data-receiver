package webserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthHandler returns Ok for all the requests.
func HealthHandler(c *gin.Context) {
	c.Status(http.StatusOK)
	return
}

// DataHandler has the logic to process the data requests.
func DataHandler(c *gin.Context) {
	service, ok := services[c.Param("service")]
	if !ok {
		err := fmt.Errorf("Service %q not found.", c.Param("service"))
		slog.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}

	extract := service.ext.Extract(c)
	err := service.ext.Validate(extract)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	err = service.auth.Authenticate(body, extract["signature"])
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	err = service.w.Write(string(body))
	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
	return
}
