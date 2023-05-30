package param

type ParamRegister struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	EMail    string `json:"e_mail" form:"e_mail" binding:"required"`
}

type ParamLogin struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}
