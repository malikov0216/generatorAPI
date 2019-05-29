package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", Handler)
	router.Run(":8080")
}

func Handler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "some text"})
}