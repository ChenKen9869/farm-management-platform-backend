package monitor

import (
	"go-backend/entity"
	"go-backend/server"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MonitorConnection struct {
	MessageChan chan string
	QuitChan chan int
}

var MonitorCentor map[uint]MonitorConnection

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary API of golang gin backend
// @Tags MonitorCentor
// @description connect with monitor centor
// @version 1.0
// @accept application/json
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /monitorCentor/connect [get]
func ConnectToMonitorCentor(ctx *gin.Context) {
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user info does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	userId := user.ID
	// 在监控中心注册管道前，先判断监控中心中是否已存在该用户的管道，如果存在则不要刷新管道，防止在掉线期间接收到的消息丢失
	// userId := uint(0)
	if _, exists := MonitorCentor[userId]; !exists {
		MonitorCentor[userId] = MonitorConnection {
			MessageChan: make(chan string),
			QuitChan: make(chan int),
		}
	}
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err.Error())
	}
	ws.SetCloseHandler(func(code int, text string) error {
		err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "closed"))
		if err != nil {
			return err
		}
		errClose := ws.Close()
		if errClose != nil {
			return errClose
		}
		return nil
	})
	defer ws.Close()
	ws.WriteMessage(websocket.TextMessage, []byte("connected"))
	ws.SetPingHandler(func(ping string) error { 
		err := ws.WriteMessage(websocket.PongMessage, []byte("pong"))
		if err != nil {
			return err
		}
		return nil
	})
	for {
		select {
		case message := <- MonitorCentor[userId].MessageChan:
			err := ws.WriteMessage(websocket.TextMessage, []byte(message))
			// 如果连接断开，则将最新取出的消息放回管道末尾
			if err != nil {
				MonitorCentor[userId].MessageChan <- message
				return
			}
		case <- MonitorCentor[userId].QuitChan:
			delete(MonitorCentor, userId)
			ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, "closed"))
			return
		}
	}
}

// @Summary API of golang gin backend
// @Tags MonitorCentor
// @description disconnect with monitor centor
// @version 1.0
// @accept application/json
// @param CompanyId query int true "company id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /monitorCentor/disconnect  [delete]
func DisconnectMonitorCentor(ctx *gin.Context) {
	userInfo, exists := ctx.Get("user")
	if !exists {
		server.Response(ctx, http.StatusInternalServerError, 500, nil, "user info does not exists in application context")
		return
	}
	user := userInfo.(entity.User)
	MonitorCentor[user.ID].Disconnect()
	server.ResponseSuccess(ctx, nil, server.Success)
}

func (con MonitorConnection)Disconnect() {
	con.QuitChan <- -1
}

func InitMonitor() {
	MonitorCentor = make(map[uint]MonitorConnection)
	ActiveFenceList = make(map[uint]ActiveFence)
}