package DO

import (
	"dou_yin/model/PO"
	"encoding/json"
	"fmt"
)

type FriendCircleDO struct {
	NewsID     int64              `json:"news_id" db:"news_id"`
	SenderID   int64              `json:"sender_id" db:"sender_id"`
	News       *string            `json:"news" db:"news"`
	Type       int                `json:"type" db:"type"`
	BlackList  *[]int64           `json:"black_list" db:"black_list"`
	WhiteList  *[]int64           `json:"white_list" db:"white_list"`
	CreateTime string             `json:"create_time" db:"create_time"`
	Likes      *[]int64           `json:"likes" db:"likes"`
	IsDeleted  bool               `json:"is_deleted" db:"is_deleted"`
	Extra      *FriendCircleExtra `json:"extra" db:"extra"`
}

type FriendCircleExtra struct {
	List *[]Comment `json:"list"`
}

type Comment struct {
	SenderID   int64  `json:"sender_id"`
	Message    string `json:"message"`
	Caller     int64  `json:"caller"`
	CreateTime string `json:"create_time"`
}

func MGetFriendCircleDOFromPO(info *PO.FriendCirclePO) (*FriendCircleDO, error) {
	result := FriendCircleDO{
		NewsID:     info.NewsID,
		SenderID:   info.SenderID,
		News:       info.News,
		Type:       info.Type,
		CreateTime: info.CreateTime,
		IsDeleted:  info.IsDeleted,
	}
	var blackIds []int64
	if info.BlackList != nil {
		err := json.Unmarshal([]byte(*info.BlackList), &blackIds)
		if err != nil {
			fmt.Println("[MGetFriendCircleDOFromPO], Unmarshal blcakList err is ", err.Error())
			return nil, err
		}
	}
	result.BlackList = &blackIds

	var whiteList []int64
	if info.WhiteList != nil {
		err := json.Unmarshal([]byte(*info.WhiteList), &whiteList)
		if err != nil {
			fmt.Println("[MGetFriendCircleDOFromPO], Unmarshal whiteList err is ", err.Error())
			return nil, err
		}
	}
	result.WhiteList = &whiteList

	var likes []int64
	if info.Likes != nil {
		err := json.Unmarshal([]byte(*info.Likes), &likes)
		if err != nil {
			fmt.Println("[MGetFriendCircleDOFromPO], Unmarshal likes err is ", err.Error())
			return nil, err
		}
	}
	result.Likes = &likes

	var extra FriendCircleExtra
	if info.Extra != nil {
		err := json.Unmarshal([]byte(*info.Extra), &extra)
		if err != nil {
			fmt.Println("[MGetFriendCircleDOFromPO], Unmarshal extra err is ", err.Error())
			return nil, err
		}
	}
	result.Extra = &extra

	return &result, nil
}

func MGetFriendCirclePOFromDO(info *FriendCircleDO) (*PO.FriendCirclePO, error) {
	result := PO.FriendCirclePO{
		NewsID:     info.NewsID,
		SenderID:   info.SenderID,
		News:       info.News,
		Type:       info.Type,
		CreateTime: info.CreateTime,
		IsDeleted:  info.IsDeleted,
	}

	var blackIds string
	if len(*info.BlackList) > 0 {
		data, err := json.Marshal(*info.BlackList)
		if err != nil {
			fmt.Println("[MGetFriendCirclePOFromDO], Marshal BlackList err is ", err.Error())
			return nil, err
		}
		blackIds = string(data)
		result.BlackList = &blackIds
	} else {
		result.BlackList = nil
	}

	var whiteList string
	if len(*info.WhiteList) > 0 {
		data, err := json.Marshal(*info.WhiteList)
		if err != nil {
			fmt.Println("[MGetFriendCirclePOFromDO], Marshal WhiteList err is ", err.Error())
			return nil, err
		}
		whiteList = string(data)
		result.WhiteList = &whiteList
	} else {
		result.WhiteList = nil
	}

	var likes string
	if len(*info.Likes) > 0 {
		data, err := json.Marshal(*info.Likes)
		if err != nil {
			fmt.Println("[MGetFriendCirclePOFromDO], Marshal Likes err is ", err.Error())
			return nil, err
		}
		likes = string(data)
		result.Likes = &likes
	} else {
		result.Likes = nil
	}

	var extra string
	if info.Extra != nil {
		data, err := json.Marshal(*info.Extra)
		if err != nil {
			fmt.Println("[MGetFriendCirclePOFromDO], Marshal Extra err is ", err.Error())
			return nil, err
		}
		extra = string(data)
		result.Extra = &extra
	} else {
		result.Extra = nil
	}

	return &result, nil
}
