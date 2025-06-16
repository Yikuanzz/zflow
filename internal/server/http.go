package server

import (
	"net/http"

	"zflow/internal/data"

	"github.com/gin-gonic/gin"
)

func NewServer() *http.Server {
	router := gin.Default()

	router.GET("/node_types", func(c *gin.Context) {
		c.JSON(http.StatusOK, data.GetDefaultNodeTypes())
	})

	router.GET("/connection_types", func(c *gin.Context) {
		c.JSON(http.StatusOK, data.GetDefaultConnectionTypes())
	})

	router.POST("/workflows", func(c *gin.Context) {

	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
