package friendCircle

import (
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 发朋友圈
func SendCircle(c *gin.Context) {
	p := new(param.SendCircleParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[CreateGroupInfo] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.SendCirclebyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[CreateGroupInfo] CreateGroupInfoByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)
}


func QueryAllFriendCircle(c *gin.Context) {
	param := new(param.QueryAllFriendCircleParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	context, err := service.QueryAllFriendCircle(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, context)
}

func QueryFriendCircle(c *gin.Context) {
	param := new(param.QueryFriendCircleParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	context, err := service.QueryFriendCircle(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}
	
	response.ResponseSuccess(c, context)
}


// 上传图片
func UploadCirclePhoto(c *gin.Context) {
	p := new(param.UploadCirclePhotoParam)
	err := c.ShouldBind(p)
	if err != nil {
		fmt.Println("[UploadCirclePhoto], ShouldBind err ", err.Error())
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("[UploadCirclePhoto], MultipartForm err ", err.Error())
		return
	}
	paths := form.File["photo"]

	var PathIDs []int64

	fmt.Println("[UploadCirclePhoto], paths is ", PathIDs, " and len(photo) is ", len(paths))

	pwd := utils.GetCurrentPath()
	for _, file := range paths {
		path := snowflake.GenID()
		PathIDs = append(PathIDs, path)

		dst := fmt.Sprintf("%v/file/%v", pwd, utils.ShiftToStringFromInt64(path))
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			fmt.Println("[UploadCirclePhoto], SaveUploadedFile err ", err.Error())
			response.ResponseError(c, response.CodeInternError)
			return
		}
	}
	fmt.Println("[UploadCirclePhoto], paths is ", PathIDs, " and len(photo) is ", len(paths))

	err = service.UploadCirclePhotoPath(p, PathIDs)
	if err != nil {
		fmt.Println("[UploadCirclePhoto], UploadCirclePhotoPath err ", err.Error())
		response.ResponseError(c, response.CodeInternError)
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 点赞朋友圈
func IsLikeCircle(c *gin.Context) {
	p := new(param.IsLikeCircleParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[IsLikeCircle] ShouldBind err is ", err.Error())
		return
	}
	err = service.IsLikeCirclebyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[IsLikeCircle]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 评论朋友圈
func CommentCircle(c *gin.Context) {
	p := new(param.CommentCircleParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[CommentCircle] ShouldBind err is ", err.Error())
		return
	}
	err = service.CommentCirclebyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[CommentCircle]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

func DeleteFriendCircle(c *gin.Context) {
	p := new(param.DeleteFriendCircleParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[DeleteFriendCircle] ShouldBind err is ", err.Error())
		return
	}
	err = service.DeleteFriendCirclebyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[DeleteFriendCircle]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}