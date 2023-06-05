package service

import (
	// "dou_yin/dao/redis"
	"dou_yin/model/VO"
	"dou_yin/pkg/jwt"
	"dou_yin/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
		if msg.MsgType == 0 { // 私聊
			if _, ok := UserChan[utils.ShiftToNum64(msg.ReceiverID)]; ok {
				fmt.Println("receive_id:", msg.ReceiverID)
				UserChan[utils.ShiftToNum64(msg.ReceiverID)] <- msg
			}
			//  else {
			// 	id := msg.ReceiverID
			// 	res, _ := json.Marshal(msg)

			// 	err := redis.AddMsg(string(res), id)
			// 	if err != nil {
			// 		fmt.Println("[MsgTransMit], 私聊 err is ", err.Error())
			// 	}
			// }
		} else { // 群聊
			userIds, err := GetAllUserIDsbyGroupID(msg.ReceiverID)
			if err != nil {
				fmt.Println("[MsgTransMit], 群聊 err is ", err.Error())
			}
			for _, id := range *userIds {
				if id != utils.ShiftToNum64(msg.SenderID) && id != 999999 {
					if _, ok := UserChan[id]; ok {
						Type, err := GroupMSGType(utils.ShiftToStringFromInt64(id), msg.ReceiverID)
						if err != nil {
							fmt.Println("[MsgTransMit], GroupMSGType err is ", err.Error())
						}
						msg.MsgType = int(Type)
						if Type == GROUP_WHITE_LIST || Type == GROUP_GRAY_LIST {
							fmt.Println("receive_group:", msg.ReceiverID, " receive_id:", id)
							UserChan[id] <- msg
						}
					}
				}
			}
		}
	}
}

func StartSendGroupNewMsg(UserId, GroupID string, Type int) error {
	msgs, err := QueryGroupNewMsgList(UserId, GroupID)
	if err != nil {
		return err
	}
	id := utils.ShiftToNum64(UserId)
	if _, ok := UserChan[id]; ok {
		for _, msg := range msgs {
			msg.MsgType = Type
			fmt.Println("receive_group:", GroupID, " receive_id:", id)
			UserChan[id] <- msg
		}
	}

	return nil
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
		} else { // 群聊
			IsSlience, err := IsGroupSlienceList(msg.SenderID, msg.ReceiverID)
			if err != nil {
				msg.ErrString = "系统内部错误，请稍后再试"
			}
			if IsSlience {
				msg.ErrString = "已被禁言"
			} else {
				err = CreatGroupMsg(*msg)
				if err != nil {
					msg.ErrString = "系统内部错误，请稍后再试"
				}
				MsgChan <- *msg
				if strings.HasPrefix(msg.Message, "@GPT") {
					result, err := GetGPTMessage(msg)
					if err != nil {
						msg.ErrString = "系统内部错误，请稍后再试"
					}
					MsgChan <- *result
				}
			}
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
