package service

import (
	"database/sql"
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/PO"
	"dou_yin/model/VO/param"
	"dou_yin/pkg/jwt"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"encoding/json"
	"errors"
	"fmt"
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

	err = user_dao.UpdateUserInfoByPO(&userInfo)
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

	err = user_dao.UpdateUserInfoByPO(&userInfo)
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
	users, err := user_dao.QueryLike(param.Context)
	if err != nil && err != sql.ErrNoRows {
		return context, err
	} else {
		for _, user := range users {
			context.Result = append(context.Result, DO.SearchFriendOrGroupContext{
				Type: 0,
				ID:   user.UserID,
				Name: user.UserName,
			})
		}
	}

	user, err := user_dao.QueryUserInfo(utils.ShiftToNum64(param.Context))
	if err != nil && err != sql.ErrNoRows {
		return context, err
	} else {
		context.Result = append(context.Result, DO.SearchFriendOrGroupContext{
			Type: 0,
			ID:   user.UserID,
			Name: user.UserName,
		})
	}

	groupInfos, err := group_dao.QueryLike(param.Context)
	if err != nil && err != sql.ErrNoRows {
		return context, err
	} else {
		for _, groupInfo := range groupInfos {
			context.Result = append(context.Result, DO.SearchFriendOrGroupContext{
				Type: 1,
				ID:   groupInfo.GroupID,
				Name: groupInfo.GroupName,
			})
		}
	}

	groupInfo, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(param.Context))
	if err != nil && err != sql.ErrNoRows {
		return context, err
	} else {
		context.Result = append(context.Result, DO.SearchFriendOrGroupContext{
			Type: 1,
			ID:   groupInfo.GroupID,
			Name: groupInfo.GroupName,
		})
	}

	return context, nil
}
