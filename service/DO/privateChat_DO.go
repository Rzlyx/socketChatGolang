package DO

import "dou_yin/model/VO"

type MessageList struct {
	Messages []VO.MessageVO `json:"messages"`
}
