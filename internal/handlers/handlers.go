package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthcheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
