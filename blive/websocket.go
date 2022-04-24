package blive

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	Host   = "blive.ericlamm.xyz"
	WsPath = "/ws?id=ddstats_client"
)

const public = "https://blive.ericlamm.xyz/ws?id=ddstats_client"

var (
	logger = logrus.WithField("service", "blive")
)

func StartWebSocket(ctx context.Context) {

	websocketHost := fmt.Sprintf("wss://%s%s", Host, WsPath)

	logrus.Debugf("prepare to connect %v", websocketHost)
	con, _, err := websocket.DefaultDialer.Dial(websocketHost, nil)

	if err != nil {
		logger.Errorf("連線到 Websocket %s 時出現錯誤: %v", websocketHost, err)
		logger.Warnf("十秒後重試")
		<-time.After(time.Second * 10)
		StartWebSocket(ctx)
		return
	}

	logger.Infof("連線到 Websocket %s 成功", websocketHost)

	con.SetCloseHandler(func(code int, text string) error {
		return con.WriteMessage(websocket.CloseMessage, nil)
	})

	go onReceiveMessage(ctx, con)
}

func onReceiveMessage(ctx context.Context, conn *websocket.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf("關閉 Websocket 時出現錯誤: %v", err)
		} else {
			logger.Debugf("連接關閉成功。")
		}
	}()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("正在關閉 Websocket...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "stop"))
			if err != nil {
				logger.Errorf("發送 websocket 關閉訊息時出現錯誤: %v", err)
			}
			return
		default:
			_, _, err := conn.ReadMessage()
			// Error
			if err != nil {
				logger.Errorf("Websocket 嘗試讀取消息時出現錯誤: %v", err)
				go retryDelay(ctx)
				return
			}
		}
	}
}

func retryDelay(ctx context.Context) {
	logger.Warnf("五秒後重連...")
	<-time.After(time.Second * 5)
	StartWebSocket(ctx)
}
