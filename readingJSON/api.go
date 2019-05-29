package main
import (
	"fmt"
	"github.com/gin-gonic/gin"
)
func main () {
	fmt.Println("asd")
	router := gin.Default()
	router.POST(":8080", handler)
}
func handler(c *gin.Context){
