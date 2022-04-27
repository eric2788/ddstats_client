package blive

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func RegisterRoute(c *gin.RouterGroup) {
	c.GET("", getSubscribe)
	c.PUT("", putSubscribe)

}

func getSubscribe(c *gin.Context) {
	rooms, err := GetSubscribes()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, rooms)
	}
}

func putSubscribe(c *gin.Context) {
	op := c.Query("type")
	rooms, ok := c.GetPostFormArray("subscribes")

	if !ok {
		c.JSON(400, gin.H{
			"message": "subscribes is required",
		})
		return
	}

	if !slices.Contains([]string{"add", "remove"}, op) {
		c.JSON(400, gin.H{
			"message": "type must be either add or remove",
		})
		return
	}

	err := PutSubscribe(rooms, op == "add")

	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "success",
		})
	}
}
