package service

import (
	friendcircle_dao "dou_yin/dao/mysql/friendCircle_dao"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
)

func SendCirclebyParam(info *param.SendCircleParam) (*response.CreateCircle, error) {
	circle := DO.FriendCircleDO{
		NewsID:   snowflake.GenID(),
		SenderID: utils.ShiftToNum64(info.Sender),
		News:     &info.News,
		Type:     info.Type,
		Extra: &DO.FriendCircleExtra{
			Paths: new([]int64),
			List:  new([]DO.Comment),
		},
	}
	var blackList []int64
	if len(info.BlackList) > 0 {
		for _, id := range info.BlackList {
			blackList = append(blackList, utils.ShiftToNum64(id))
		}
	}
	circle.BlackList = &blackList
	circle.WhiteList = new([]int64)

	circlePO, err := DO.MGetFriendCirclePOFromDO(&circle)
	if err != nil {
		return nil, err
	}

	err = friendcircle_dao.CreateFriendCircle(circlePO)
	if err != nil {
		return nil, err
	}
	return &response.CreateCircle{
		NewsID: circle.NewsID,
	}, nil
}

func UploadCirclePhotoPath(info *param.UploadCirclePhotoParam, paths []int64) error {
	circlePO, err := friendcircle_dao.MGetFriendCircle(utils.ShiftToNum64(info.NewsID))
	if err != nil {
		return err

	}
	circleDO, err := DO.MGetFriendCircleDOFromPO(circlePO)
	if err != nil {
		return err
	}

	Paths := append(*circleDO.Extra.Paths, paths...)
	circleDO.Extra.Paths = &Paths

	CirclePO, err := DO.MGetFriendCirclePOFromDO(circleDO)
	if err != nil {
		return err
	}

	err = friendcircle_dao.UpdateFriendCircle(CirclePO)
	if err != nil {
		return err
	}
	return nil
}