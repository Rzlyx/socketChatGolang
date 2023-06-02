package service

import (
	"dou_yin/dao/mysql"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/PO"
	"dou_yin/model/VO/param"
	"dou_yin/pkg/jwt"
	"dou_yin/pkg/snowflake"
	"dou_yin/service/DO"
	"encoding/json"
	"errors"
	"fmt"
)

func Register(p *param.ParamRegister) (err error, user *PO.User) {
	p1 := new(PO.User)
	p1.UserName = p.UserName
	p1.Password = p.Password
	p1.EMail = p.EMail
	p1.UserID = snowflake.GenID() / 100000000000
	err = mysql.Register(p1)
	fmt.Println(err)
	return err, p1
}

func Login(p *param.ParamLogin) (err error, user *PO.UserPO, token string) {
	user, err = mysql.Login(p.UserName)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("用户不存在"), nil, ""
	}

	if p.UserName == p.UserName && p.Password == user.Password {
		token, err = jwt.GenToken(user.UserID, user.UserName)
		return err, user, token
	}

	return errors.New("信息错误"), nil, ""
}

func GetContactorList(Id string) (*PO.ContactorList, error) {

	p, err := mysql.GetContactorList(Id)

	return p, err
}

func QueryContactorList(param param.QueryContactorList) (contactors DO.ContactList, err error) {
	userInfo, err := user_dao.QueryUserInfo(param.UserID)
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

	for _, contact := range extra.ContactorList {
		contactDO := DO.ContactInfo{
			ID:           contact.ID,
			Name:         contact.Name,
			Message:      contact.Message,
			FriendshipID: contact.FriendshipID,
		}

		contactors.ContactorList = append(contactors.ContactorList, contactDO)
	}

	return contactors, err
}

func SetContactorList(param param.SetContactorList) (err error) {
	// user_dao.Update
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserInfo(param param.UpdateUserInfo) (err error) {
	userInfo, err := user_dao.QueryUserInfo(param.UserID)
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

	// user_dao.Update
	if err != nil {
		return err
	}

	return nil
}
