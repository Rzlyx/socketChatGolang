package param

import "dou_yin/service/DO"

type ParamRegister struct {
	UserName    string `json:"username" form:"username" binding:"required"`
	Password    string `json:"password" form:"password" binding:"required"`
	Sex         int    `json:"sex" form:"sex" binding:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required"`
	EMail       string `json:"e_mail" form:"e_mail" binding:"required"`
	Signature 	string `json:"signature" form:"signature" binding:"required"`
}

type ParamLogin struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type QueryContactorListParam struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
}

type SetContactorListParam struct {
	UserID        int64            `json:"user_id"`
	ContactorList []DO.ContactInfo `json:"contactor_list"`
}

type UpdateUserInfoParam struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	Sex         int    `json:"sex"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"e_mail"`
	Signature   string `json:"signature"`
	Birthday    string `json:"birthday"`
	Status      int    `json:"status"`
	//PrivateChatWhite     *string `json:"private_chat_white"`
	//PrivateChatBlack     *string `json:"private_chat_black"`
	//FriendCircleWhite    *string `json:"friend_circle_white"`
	//FriendCircleBlack    *string `json:"friend_circle_black"`
	FriendCircleVisiable int `json:"friend_circle_visiable"`
	//GroupChatWhite       *string `json:"group_chat_white""`
	//GroupChatBlack       *string `json:"group_chat_black"`
	//GroupChatGray        *string `json:"group_chat_grey"`
	//CreateTime           *string `json:"create_time"`
	//IsDeleted            bool    `json:"is_deleted"`
	//Extra                *string `json:"extra"`
}
