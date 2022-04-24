package main

import (
	"bytes"
	"context"
	"ddstats_client/blive"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"

	_ "embed"
)

//go:embed static/index.html
var indexHtml []byte

func main() {

	blive.StartWebSocket(context.Background())

	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		c.DataFromReader(200, int64(len(indexHtml)), "text/html", bytes.NewReader(indexHtml), map[string]string{})
	})

	route.POST("/offline", func(c *gin.Context) {
		rooms, ok := c.GetPostFormArray("subscribes")

		if !ok {
			c.JSON(400, gin.H{
				"error": "缺少 subscribes 列表参数",
			})
			return
		}

		if err := SaveOffline(rooms); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		} else {
			c.JSON(200, gin.H{
				"success": true,
			})
		}
	})
	route.GET("/roomName/:room_id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("room_id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		name, err := blive.GetUserName(id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"name": name,
		})
	})

	if rooms, ok := GetFromOffline(); ok {
		if err := blive.SubscribeRequest(rooms); err != nil {
			logrus.Errorf("尝试从离线重新订阅时出现错误: %v", err)
		} else {
			logrus.Infof("成功从离线重新订阅 %v 个房间。", len(rooms))
		}
	}

	if err := route.Run(":9090"); err != nil {
		logrus.Fatal(err)
	}

}

func SaveOffline(room []string) error {
	_ = os.MkdirAll("data", os.ModePerm)
	_ = os.Remove("data/offline.json")
	f, err := os.Create("data/offline.json")
	if err != nil {
		return err
	}
	b, err := json.Marshal(room)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

func GetFromOffline() ([]string, bool) {
	content, err := os.ReadFile("data/offline.json")
	if err != nil {
		return nil, false
	}
	var rooms []string
	err = json.Unmarshal(content, &rooms)
	if err != nil {
		return nil, false
	}
	return rooms, true
}
