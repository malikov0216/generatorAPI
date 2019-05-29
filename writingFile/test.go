package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/deployment/create", handlerAPI)
	router.Run(":1234")
}
func handlerAPI(c *gin.Context) {
	var receivingData []map[string]string
	receivingData = append(receivingData, map[string]string{
		"deployment-name": c.PostForm("deployment-name"),
		"in":              "formData",
		"type":            "string",
	})
	var sliceData []map[string]interface{}
	sliceData = append(sliceData, map[string]interface{}{
		"asd": c.FormFile("asda"),
		"bkb": "asda"
	})
}
