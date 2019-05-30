package main

import (
	gin "github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()
	router.POST("/deployment/create", handlerAPI)
	router.Run(":80")
}
func handlerAPI(c *gin.Context) {
	body := c.Request.Body
	header := c.Request.Header
	method := c.Request.Method
	endpoint := c.Request.RequestURI
	timeout := time.Duration(10 * time.Second)
	client := http.Client{Timeout: timeout}
	defer body.Close()
	request, err := http.NewRequest(method, "http://localhost:8080"+endpoint, body)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}
	request.Header = header
	response, err := client.Do(request)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}
	defer response.Body.Close()
	bodyResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"result": string(bodyResp)})
}
