package response

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParams
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedAuth
	CodeInvalidToken
	CodeNeedLogin
	CodePayError
	CodeOrderNotExist

	CodeInternError

	CodeNotFriend
	CodeNotApplied
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "无效的参数",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "无效密码",
	CodeServerBusy:      "系统繁忙",

	CodeNeedAuth:      "需要Auth",
	CodeInvalidToken:  "无效的Token",
	CodeNeedLogin:     "需要登录",
	CodePayError:      "支付失败",
	CodeOrderNotExist: "订单不存在或失效",

	CodeInternError: "系统内部错误",

	CodeNotFriend:  "不是好友",
	CodeNotApplied: "没有申请",
}

func (c ResCode) Msg() string {
	v, ok := codeMsgMap[c]
	if !ok {
		v = codeMsgMap[CodeServerBusy]
	}
	return v
}
