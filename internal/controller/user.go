package controller

import (
	"Nuxus/internal/dto"
	"Nuxus/internal/middleware"
	"Nuxus/internal/res"
	"Nuxus/internal/service"
	"Nuxus/pkg/erru"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var reqDTO dto.RegisterReqDTO
	err := c.ShouldBindJSON(&reqDTO)
	if err != nil {
		res.FailWithAppErr(c, erru.ErrInvalidParams.Wrap(err))
		return
	}
	// log.Println(reqDTO)
	err = service.Register(&reqDTO)
	if err != nil {
		c.Error(err)
		return
	}

	res.OkWithMsg(c, "验证码已发送，请确认")
}

func VerifyRegister(c *gin.Context) {
	var reqDto dto.VerifyRegisterReqDTO
	err := c.ShouldBindJSON(&reqDto)
	if err != nil {
		res.FailWithAppErr(c, erru.ErrInvalidParams.Wrap(err))
		return
	}

	err = service.VerifyRegister(&reqDto)
	if err != nil {
		c.Error(err)
		return
	}
	res.OkWithMsg(c, "注册成功")
}

func Login(c *gin.Context) {
	var reqDTO dto.LoginReqDTO
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		res.FailWithAppErr(c, erru.ErrInvalidParams.Wrap(err))
		return
	}

	user, err := service.Login(&reqDTO)
	if err != nil {
		c.Error(err)
	}
	//token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		res.FailWithAppErr(c, erru.New("token 生成错误"))
		return
	}
	// encapsulate
	userInfo := dto.UserInfoDTO{
		ID:       user.ID,
		UserName: user.Username,
		Email:    user.Email,
	}
	resDTO := dto.LoginResponseDTO{
		User:  userInfo,
		Token: token,
	}
	res.OkWithData(c, resDTO)
}

func RequestReset(c *gin.Context) {
	var reqDTO dto.RequestResetReqDTO
	err := c.ShouldBindJSON(&reqDTO)
	if err != nil {
		res.FailWithAppErr(c, erru.ErrInvalidParams.Wrap(err))
	}
	// log.Println(reqDTO)
	err = service.RequestReset(&reqDTO)
	if err != nil {
		c.Error(err)
		return
	}
	res.OkWithMsg(c, "验证码已发送，请确认")
}

func VerifyReset(c *gin.Context) {
	var reqDto dto.VerifyResetReqDTO
	err := c.ShouldBindJSON(&reqDto)
	if err != nil {
		res.FailWithAppErr(c, erru.ErrInvalidParams.Wrap(err))
		return
	}

	err = service.VerifyReset(&reqDto)
	if err != nil {
		c.Error(err)
		return
	}
	res.OkWithMsg(c, "密码重置成功，请妥善保管")
}
