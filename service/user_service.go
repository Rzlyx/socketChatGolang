package service

import (
	"database/sql"
	"dou_yin/dao/mysql"
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"dou_yin/model/VO"
	"dou_yin/model/VO/param"
	"dou_yin/pkg/jwt"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

func Register(info *param.ParamRegister) (*PO.UserPO, error) {
	userInfo := &PO.UserPO{
		UserID:      snowflake.GenID() / 100000000000,
		UserName:    info.UserName,
		Password:    info.Password,
		Sex:         info.Sex,
		PhoneNumber: info.PhoneNumber,
		Email:       info.EMail,
		Signature:   &info.Signature,
		Birthday:    info.Birthday,
	}

	err := user_dao.Register(userInfo)
	fmt.Println("[Register], err is ", err.Error())
	return userInfo, err
}

func Login(p *param.ParamLogin) (user *PO.UserPO, token string, err error) {
	user, err = user_dao.Login(p.UserName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", errors.New("用户不存在")
	}

	if p.UserName == user.UserName && p.Password == user.Password {
		token, err = jwt.GenToken(user.UserID, user.UserName)

		return user, token, err
	}

	return nil, "", errors.New("信息错误")
}

func GetContactorList(Id string) (*PO.ContactorList, error) {

	p, err := user_dao.GetContactorList(Id)

	return p, err
}

func QueryContactorList(param param.QueryContactorListParam) (contactors DO.ContactList, err error) {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.UserID))
	if err != nil {
		return contactors, err
	}

	var extra PO.UserExtra
	if userInfo.Extra != nil {
		err = json.Unmarshal([]byte(*userInfo.Extra), &extra)
		if err != nil {
			return contactors, err
		}
	}

	if extra.ContactorList != nil {
		for _, contact := range *extra.ContactorList {
			contactDO := DO.ContactInfo{
				ID:      utils.ShiftToStringFromInt64(contact.ID),
				Name:    contact.Name,
				Message: contact.Message,
				Time:    contact.Time,
			}

			contactors.ContactorList = append(contactors.ContactorList, contactDO)
		}
	}

	return contactors, err
}

func SetContactorList(param param.SetContactorListParam) (err error) {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.UserID))
	if err != nil {
		return err
	}

	var extra PO.UserExtra
	if userInfo.Extra != nil {
		err = json.Unmarshal([]byte(*userInfo.Extra), &extra)
		if err != nil {
			return err
		}
	}

	if len(param.ContactorList) != 0 {
		var contactorsPO []PO.ContactInfoPO
		for _, item := range param.ContactorList {
			contactor := PO.ContactInfoPO{
				ID:      utils.ShiftToNum64(item.ID),
				Name:    item.Name,
				Message: item.Message,
				Time:    item.Time,
				Status:  item.Status,
			}

			contactorsPO = append(contactorsPO, contactor)
		}
		extra.ContactorList = &contactorsPO
	} else {
		extra.ContactorList = nil
	}
	extraJson, err := json.Marshal(extra)
	if err != nil {
		return err
	}
	extraStr := string(extraJson[:])
	userInfo.Extra = &extraStr

	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = user_dao.UpdateUserInfoByPO(tx, &userInfo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserInfo(param param.UpdateUserInfoParam) (err error) {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.UserID))
	if err != nil {
		return err
	}

	// todo: tx
	needUpdateRemark := false
	if userInfo.UserName != param.UserName {
		needUpdateRemark = true
	}

	userInfo.UserName = param.UserName
	userInfo.Password = param.Password
	userInfo.Sex = param.Sex
	if userInfo.Signature == nil {
		userInfo.Signature = new(string)
	}
	*userInfo.Signature = param.Signature
	userInfo.Status = param.Status
	userInfo.PhoneNumber = param.PhoneNumber
	userInfo.Email = param.Email
	userInfo.Birthday = param.Birthday
	userInfo.FriendCircleVisiable = param.FriendCircleVisiable

	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = user_dao.UpdateUserInfoByPO(tx, &userInfo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if needUpdateRemark {
		err = UpdateFriendRemark(utils.ShiftToNum64(param.UserID), param.UserName)
		if err != nil {
			return err
		}
	}
	return nil
}

func QueryUserInfo(param param.QueryUserInfoParam) (DO.UserInfo, error) {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.UserID))
	if err != nil {
		return DO.UserInfo{}, err
	}

	ret := DO.UserInfo{
		UserID:               userInfo.UserID,
		UserName:             userInfo.UserName,
		Password:             userInfo.Password,
		Sex:                  userInfo.Sex,
		PhoneNumber:          userInfo.PhoneNumber,
		Email:                userInfo.Email,
		Birthday:             userInfo.Birthday,
		Status:               userInfo.Status,
		FriendCircleVisiable: userInfo.FriendCircleVisiable,
	}
	if userInfo.Signature != nil {
		ret.Signature = *userInfo.Signature
	}

	return ret, nil
}

func SearchFriendOrGroup(param param.SearchFriendOrGroupParam) (context DO.SearchFriendOrGroupContexts, err error) {
	fmt.Println(param)
	users, err := user_dao.QueryLike(param.Context)
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Error(err.Error())
		return context, err
	} else {
		for _, user := range users {
			ctx := DO.SearchFriendOrGroupContext{
				Type: 0,
				ID:   user.UserID,
				Name: user.UserName,
			}
			if !IsSearchFriendOrGroupContextContain(context.Result, ctx) {
				context.Result = append(context.Result, ctx)
			}
		}
	}

	user, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.Context))
	if err != nil && err != sql.ErrNoRows {
		return context, err
	} else {
		if err != sql.ErrNoRows {
			ctx := DO.SearchFriendOrGroupContext{
				Type: 0,
				ID:   user.UserID,
				Name: user.UserName,
			}
			if !IsSearchFriendOrGroupContextContain(context.Result, ctx) {
				context.Result = append(context.Result, ctx)
			}
		}
	}

	groupInfos, err := group_dao.QueryLike(param.Context)
	if err != nil && err != sql.ErrNoRows {
		return context, err
	} else {
		for _, groupInfo := range groupInfos {
			ctx := DO.SearchFriendOrGroupContext{
				Type: 1,
				ID:   groupInfo.GroupID,
				Name: groupInfo.GroupName,
			}
			if !IsSearchFriendOrGroupContextContain(context.Result, ctx) {
				context.Result = append(context.Result, ctx)
			}
		}
	}

	groupInfo, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(param.Context))
	if err != nil && err != sql.ErrNoRows {
		return context, err
	} else {
		if err != sql.ErrNoRows {
			ctx := DO.SearchFriendOrGroupContext{
				Type: 1,
				ID:   groupInfo.GroupID,
				Name: groupInfo.GroupName,
			}
			if !IsSearchFriendOrGroupContextContain(context.Result, ctx) {
				context.Result = append(context.Result, ctx)
			}
		}
	}

	return context, nil
}

func IsSearchFriendOrGroupContextContain(list []DO.SearchFriendOrGroupContext, ctx DO.SearchFriendOrGroupContext) bool {
	for _, item := range list {
		if item.ID == ctx.ID {
			return true
		}
	}
	return false
}

func LogIn(userID int64, conn *websocket.Conn) (err error) {
	fmt.Printf("func(LogIn): [param: %v]\n", userID)
	userInfo, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	userInfo.Status = 1
	fmt.Println(userInfo)
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = user_dao.UpdateUserInfoByPO(tx, &userInfo)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	go SendHeartBeat(userID, conn)

	return nil
}

func SendHeartBeat(userID int64, conn *websocket.Conn) {
	var interval int64
	interval = 0
	for {
		UserChan[userID] <- VO.MessageVO{
			MsgType:    999,
			ReceiverID: utils.ShiftToStringFromInt64(userID),
		}
		fmt.Printf("func(SendHeartBeat): [param: %v]\n", userID)
		select {
		case <-UserHeartBeat[userID]:
			interval = 0
		case <-time.After(time.Second):
			logger.Log.Error("维持心跳失败 userID: " + utils.ShiftToStringFromInt64(userID) + " " + utils.ShiftToStringFromInt64(interval))
			interval++
			if interval == 5 {
				logger.Log.Error("维持心跳失败 userID: " + utils.ShiftToStringFromInt64(userID))
				err := LogOut(userID, conn)
				if err != nil {
					logger.Log.Error(err.Error())
				}
				return
			}
		}
		time.Sleep(time.Second * 2)
	}
}

func LogOut(userID int64, conn *websocket.Conn) (err error) {
	delete(UserHeartBeat, userID)
	delete(UserChan, userID)

	userInfo, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	userInfo.Status = 0
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = user_dao.UpdateUserInfoByPO(tx, &userInfo)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}
