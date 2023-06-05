package service

import (
	"dou_yin/dao/redis"
	"dou_yin/model/VO"
	"dou_yin/pkg/jwt"
	"dou_yin/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var OtherMsgChan chan string
var OthserChan map[int64]chan string

var MsgChan chan VO.MessageVO            // 全局消息队列
var UserChan map[int64]chan VO.MessageVO // 每个用户分配一个chan
var IDChan chan int64                    // 接收login，分配一个chan

func ChanInit() {
	IDChan = make(chan int64, 100)
	MsgChan = make(chan VO.MessageVO, 10000)
	UserChan = make(map[int64]chan VO.MessageVO)
	UserChan[1] = make(chan VO.MessageVO, 10)
}

func MsgTransMit() {
	for msg := range MsgChan {
		if _, ok := UserChan[utils.ShiftToNum64(msg.ReceiverID)]; ok {
			fmt.Println("receive_id:", msg.ReceiverID)
			UserChan[utils.ShiftToNum64(msg.ReceiverID)] <- msg
		} else {
			id := msg.ReceiverID
			res, _ := json.Marshal(msg)

			err := redis.AddMsg(string(res), id)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func AddUser() {
	for msg := range IDChan {
		UserChan[msg] = make(chan VO.MessageVO, 10)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Connect(c *gin.Context) {
	//获取token,如果token无效，就返回
	token := c.Param("token")
	mc, err := jwt.ParseToken(token)
	IDChan <- mc.ID
	//连接升级
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	go func(userID int64) {
		fmt.Println(userID)
		for {
			message := <-UserChan[userID]
			msgJson, _ := json.Marshal(message)
			err = conn.WriteMessage(websocket.TextMessage, msgJson)
			if err != nil {
				fmt.Println(err)
			}
		}
	}(mc.ID)
	//ID233:=utils.ShiftToStringFromInt64(mc.ID)
	msg := new(VO.MessageVO)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println(err)
		}
		// filter
		if msg.MsgType == 0 {
			HandlePrivateChatMsg(*msg)
		}
	}
}

/********测试专用*******/
//go func() {
//	ticker := time.NewTicker(time.Second) // 创建每秒触发的定时器
//	defer ticker.Stop()
//
//	for range ticker.C {
//		message := "Hello, world!" // 要发送的消息内容
//
//		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//}()
/********测试专用*******/

func Connect2(c *gin.Context) {
	token := c.Param("token")
	mc, err := jwt.ParseToken(token)
	IDChan <- mc.ID
	//连接升级
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	go func(userID int64) {
		fmt.Println(userID)
		for {

			err = conn.WriteMessage(websocket.TextMessage, []byte(""))
			if err != nil {
				fmt.Println(err)
			}
		}
	}(mc.ID)
	//ID233:=utils.ShiftToStringFromInt64(mc.ID)
	var msg interface{}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println(err)
		}

	}
}
