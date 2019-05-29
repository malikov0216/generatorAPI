package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
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
	request, err := http.NewRequest(method, "localhost:8080/"+endpoint, body)
	if err != nil {
		log.Fatal(err)
	}
	request.Header = header
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	bodyResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bodyResp))
}
