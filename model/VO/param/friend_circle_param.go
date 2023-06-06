package param

import "dou_yin/model/VO"

type SendCircleParam struct {
	Sender    string   `json:"sender" form:"sender" binding:"required"`
	News      string   `json:"news" form:"news" binding:"required"`
	Type      int      `json:"type" form:"type" binding:"required"`
	BlackList []string `json:"black_list" form:"black_list" binding:"required"`
	CircleType string `` // 0-私密 1-公开
	// WhiteList []string `json:"" form:"" binding:"required"`
}

type UploadCirclePhotoParam struct {
	NewsID  string       `json:"news_id" form:"news_id" binding:"required"`
	Message VO.MessageVO `json:"message" binding:"required"`
}
