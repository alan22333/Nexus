package controller

import (
	"Nuxus/internal/dto"
	"Nuxus/internal/middleware"
	"Nuxus/internal/models"
	"Nuxus/internal/res"
	"Nuxus/internal/service"
	"Nuxus/pkg/erru"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var reqDTO dto.RegisterReqDTO
	err := c.ShouldBindJSON(&reqDTO)
	if err != nil {
		c.Error(erru.ErrInvalidParams.Wrap(err))
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
		c.Error(erru.ErrInvalidParams.Wrap(err))
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
		c.Error(erru.ErrInvalidParams.Wrap(err))
		return
	}

	user, err := service.Login(&reqDTO)
	if err != nil {
		c.Error(err)
	}
	//token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.Error(erru.New("token 生成错误"))
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
		c.Error(erru.ErrInvalidParams.Wrap(err))
		return
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
		c.Error(erru.ErrInvalidParams.Wrap(err))
		return
	}

	err = service.VerifyReset(&reqDto)
	if err != nil {
		c.Error(err)
		return
	}
	res.OkWithMsg(c, "密码重置成功，请妥善保管")
}

func userModel2InfoDto(user *models.User) *dto.UserInfoDTO {
	return &dto.UserInfoDTO{
		ID:       user.ID,
		UserName: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}
}

// -------------------用户信息-----------------------------
func GetProfile(c *gin.Context) {
	userId := c.MustGet("userID").(uint)

	user, err := service.GetProfile(userId)
	if err != nil {
		c.Error(err)
		return
	}
	resDto := userModel2ProfileDto(user)

	res.OkWithData(c, resDto)
}

func UpdateProfile(c *gin.Context) {
	var reqDto dto.ProfileReqDTO
	err := c.ShouldBindJSON(&reqDto)
	if err != nil {
		res.FailWithAppErr(c, erru.ErrInvalidParams.Wrap(err))
		return
	}

	userId := c.MustGet("userID").(uint)

	user, err := service.UpdateProfile(userId, reqDto)
	if err != nil {
		c.Error(err)
		return
	}

	resDto := userModel2ProfileDto(user)

	res.OkWithData(c, resDto)
}

func userModel2ProfileDto(user *models.User) *dto.ProfileResDTO {
	return &dto.ProfileResDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Avatar:   user.Avatar,

		Gender: user.Gender,
		Phone:  user.Phone,
		QQ:     user.QQ,
		Wechat: user.Wechat,
		Bio:    user.Bio,

		Privacy: dto.PrivacyInfo{
			IsPhonePublic:  user.IsPhonePublic,
			IsEmailPublic:  user.IsEmailPublic,
			IsQQPublic:     user.IsQQPublic,
			IsWechatPublic: user.IsWechatPublic,
			IsGenderPublic: user.IsGenderPublic,
		},

		CreatedAt: user.CreatedAt,
	}
}

// --------------------头像------------------
//处理头像上传请求
func UpdateAvatar(c *gin.Context) {
	// 1. 从表单获取文件，"avatar" 是前端上传时文件字段的 name
	file, err := c.FormFile("avatar")
	if err != nil {
		_ = c.Error(erru.ErrInvalidParams.Wrap(err))
		return
	}

	// 2. 基础文件校验
	// 限制大小为 5MB
	if file.Size > 8*1024*1024 {
		_ = c.Error(erru.New("图片大小不能超过8MB"))
		return
	}
	// 限制文件类型
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		_ = c.Error(erru.New("只支持上传 jpg, jpeg, png 格式的图片"))
		return
	}

	// 3. 从 JWT 中间件设置的 context 获取用户 ID
	userID, _ := c.Get("userID")

	// 4. 调用 Service 层处理核心逻辑
	avatarURL, err := service.UpdateAvatar(userID.(uint), file)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// 5. 构造并返回成功响应
	res.OkWithData(c, dto.UpdateAvatarResDTO{AvatarURL: avatarURL})
}
