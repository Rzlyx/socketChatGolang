package service

import (
	"dou_yin/dao/mysql"
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

func Register(p *param.ParamRegister) (*PO.User, error) {
	p1 := new(PO.User)
	p1.UserName = p.UserName
	p1.Password = p.Password
	p1.EMail = p.EMail
	p1.UserID = snowflake.GenID() / 100000000000
	err := mysql.Register(p1)
	fmt.Println("[Register], err is ", err.Error())
	return p1, err
}

func Login(p *param.ParamLogin) (user *PO.UserPO, token string, err error) {
	user, err = mysql.Login(p.UserName)
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

	p, err := mysql.GetContactorList(Id)

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

	return nil
}
