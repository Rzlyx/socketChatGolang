package response

import "dou_yin/service/DO"

type QueryPrivateChatMsgResp struct {
	MessageList DO.MessageList `json:"message_list"`
}
