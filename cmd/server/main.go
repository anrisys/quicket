package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	sayHi := func(c *gin.Context)  {
		c.JSON(http.StatusOK, "hi")
	}

	router.GET("/hi", sayHi)

	router.Run("localhost:8080")
}