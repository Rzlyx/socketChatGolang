package service

import (
	"database/sql"
	"dou_yin/dao/mysql/friendCircle_dao"
	"dou_yin/dao/mysql/friend_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/PO"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"encoding/json"
	"sort"
)

func QueryAllFriendCircle(param param.QueryAllFriendCircleParam) ([]response.FriendCircleContext, error) {
	friendships, err := friend_dao.QueryFriendshipList(utils.ShiftToNum64(param.UserID))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.UserID))
	if err != nil {
		return nil, err
	}
	var friendCircleBlackList PO.FriendCircleBlack
	if userInfo.FriendCircleBlack != nil {
		err = json.Unmarshal([]byte(*userInfo.FriendCircleBlack), &friendCircleBlackList)
		if err != nil {
			return nil, err
		}
	}

	var context []DO.FriendCircleDO
	myFriendCirclePOs, err := friendCircle_dao.QueryFriendCircle(utils.ShiftToNum64(param.UserID))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	for _, friendCirclePO := range myFriendCirclePOs {
		friendCircleDO, err := DO.MGetFriendCircleDOFromPO(&friendCirclePO)
		if err != nil {
			return nil, err
		}
		friendCircleDO.SenderName = userInfo.UserName

		context = append(context, *friendCircleDO)
	}

	for _, friendship := range friendships {
		var friendID int64
		var friendName string
		if friendship.FirstID == utils.ShiftToNum64(param.UserID) {
			friendID = friendship.SecondID
			friendName = friendship.FirstRemarkSecond
		} else {
			friendID = friendship.FirstID
			friendName = friendship.SecondRemarkFirst
		}

		// 判断是否把对方朋友圈拉黑
		if IsContains(friendCircleBlackList.BlackList, friendID) {
			continue
		}

		// 判断对方朋友圈可见性
		friendInfo, err := user_dao.QueryUserInfo(friendID)
		if err != nil {
			return nil, err
		}
		endTime := ""
		switch friendInfo.FriendCircleVisiable {
		case 1:
			endTime = utils.GetTime(0, 0, -1)
		case 2:
			endTime = utils.GetTime(0, -1, 0)
		case 999:
			continue
		}

		friendCirclePOs, err := friendCircle_dao.QueryFriendCircle(friendID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		for _, friendCirclePO := range friendCirclePOs {
			if endTime != "" {
				if endTime > friendCirclePO.CreateTime {
					continue
				}
			}
			friendCircleDO, err := DO.MGetFriendCircleDOFromPO(&friendCirclePO)
			if err != nil {
				return nil, err
			}
			friendCircleDO.SenderName = friendName

			// 判断对方朋友圈该条数据是否把我拉黑
			if friendCircleDO.BlackList != nil {
				if IsContains(*friendCircleDO.BlackList, utils.ShiftToNum64(param.UserID)) {
					continue
				}
			}

			context = append(context, *friendCircleDO)
		}
	}

	sort.SliceStable(context, func(i, j int) bool {
		return context[i].CreateTime > context[j].CreateTime
	})

	res := new([]response.FriendCircleContext)
	for index, item := range context {
		if param.ReadTime >= item.CreateTime {
			for i := index; i < param.Num && i < len(context); i++ {
				likes := new([]string)
				if context[i].Likes != nil {
					friendInfo, err := user_dao.QueryUserInfo(context[i].SenderID)
					if err != nil {
						return nil, err
					}

					*likes = append(*likes, friendInfo.UserName)
				}

				*res = append(*res, response.FriendCircleContext{
					NewsID:     context[i].NewsID,
					SenderID:   context[i].SenderID,
					SenderName: context[i].SenderName,
					News:       context[i].News,
					Type:       context[i].Type,
					CreateTime: context[i].CreateTime,
					Likes:      likes,
				})
			}
		}
	}

	return *res, nil
}

func QueryFriendCircle(param param.QueryFriendCircleParam) ([]response.FriendCircleContext, error) {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.UserID))
	if err != nil {
		return nil, err
	}
	var friendCircleBlackList PO.FriendCircleBlack
	if userInfo.FriendCircleBlack != nil {
		err = json.Unmarshal([]byte(*userInfo.FriendCircleBlack), &friendCircleBlackList)
		if err != nil {
			return nil, err
		}
	}

	var context []DO.FriendCircleDO
	if param.UserID == param.FriendID {
		myFriendCirclePOs, err := friendCircle_dao.QueryFriendCircle(utils.ShiftToNum64(param.UserID))
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		for _, friendCirclePO := range myFriendCirclePOs {
			friendCircleDO, err := DO.MGetFriendCircleDOFromPO(&friendCirclePO)
			if err != nil {
				return nil, err
			}
			friendCircleDO.SenderName = userInfo.UserName

			context = append(context, *friendCircleDO)
		}
	} else {
		friendship, err := friend_dao.QueryFriendshipBy2ID(utils.ShiftToNum64(param.UserID), utils.ShiftToNum64(param.FriendID))
		if err != nil {
			return nil, err
		}

		var friendID int64
		var friendName string
		if friendship.FirstID == utils.ShiftToNum64(param.UserID) {
			friendID = friendship.SecondID
			friendName = friendship.FirstRemarkSecond
		} else {
			friendID = friendship.FirstID
			friendName = friendship.SecondRemarkFirst
		}

		// 判断是否把对方朋友圈拉黑
		if IsContains(friendCircleBlackList.BlackList, friendID) {
			return nil, nil
		}

		// 判断对方朋友圈可见性
		friendInfo, err := user_dao.QueryUserInfo(friendID)
		if err != nil {
			return nil, err
		}
		endTime := ""
		switch friendInfo.FriendCircleVisiable {
		case 1:
			endTime = utils.GetTime(0, 0, -1)
		case 2:
			endTime = utils.GetTime(0, -1, 0)
		case 999:
			return nil, nil
		}

		friendCirclePOs, err := friendCircle_dao.QueryFriendCircle(friendID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		for _, friendCirclePO := range friendCirclePOs {
			if endTime != "" {
				if endTime > friendCirclePO.CreateTime {
					continue
				}
			}
			friendCircleDO, err := DO.MGetFriendCircleDOFromPO(&friendCirclePO)
			if err != nil {
				return nil, err
			}
			friendCircleDO.SenderName = friendName

			// 判断对方朋友圈该条数据是否把我拉黑
			if friendCircleDO.BlackList != nil {
				if IsContains(*friendCircleDO.BlackList, utils.ShiftToNum64(param.UserID)) {
					continue
				}
			}

			context = append(context, *friendCircleDO)
		}
	}

	sort.SliceStable(context, func(i, j int) bool {
		return context[i].CreateTime > context[j].CreateTime
	})

	res := new([]response.FriendCircleContext)
	for index, item := range context {
		if param.ReadTime >= item.CreateTime {
			for i := index; i < param.Num && i < len(context); i++ {
				likes := new([]string)
				if context[i].Likes != nil {
					friendInfo, err := user_dao.QueryUserInfo(context[i].SenderID)
					if err != nil {
						return nil, err
					}

					*likes = append(*likes, friendInfo.UserName)
				}

				*res = append(*res, response.FriendCircleContext{
					NewsID:     context[i].NewsID,
					SenderID:   context[i].SenderID,
					SenderName: context[i].SenderName,
					News:       context[i].News,
					Type:       context[i].Type,
					CreateTime: context[i].CreateTime,
					Likes:      likes,
				})
			}
		}
	}

	return *res, nil
}
