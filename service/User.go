package service

import (
	"dou_yin/dao/mysql"
	"dou_yin/model/PO"
	"dou_yin/model/VO/param"
	"dou_yin/pkg/jwt"
	"dou_yin/pkg/snowflake"
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

func Login(p *param.ParamLogin) (err error, user *PO.User, token string) {
	err, user = mysql.Login(p.UserName)
	if err != nil {
		return errors.New("用户不存在"), nil, ""
	}
	if p.UserName == p.UserName && p.Password == user.Password {
		token, err = jwt.GenToken(user.UserID, user.UserName)
		return err, user, token
	}
	return errors.New("信息错误"), nil, ""
}

func GetContactorList(Id string) (*PO.ContactorList, error) {

	err, p := mysql.GetContactorList(Id)

	return p, err
}
