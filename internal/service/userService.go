package service

import (
	"Nuxus/internal/dao"
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"Nuxus/pkg/erru"
	"Nuxus/pkg/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(reqDto *dto.RegisterReqDTO) error {
	// 注册流程：
	// 1. 检查邮箱用户是否存在
	// 2. 生成验证码，放在redis里
	// 3. 发送验证码
	_, err := dao.GetUserByEmail(reqDto.Email)
	if err == nil {
		return erru.ErrEmailAlreadyUsed.Wrap(err)
	}

	isCool, err := dao.CheckSendCooldown(reqDto.Email)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	if isCool {
		return erru.New("操作太频繁，请稍后再试")
	}

	code := utils.GenerateRandomCode(6)

	err = dao.SetVerifyCode(reqDto.Email, code, 5*time.Minute)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if err := SendRegisterMail(reqDto.Email, code); err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	return nil
}

func VerifyRegister(reqDto *dto.VerifyRegisterReqDTO) error {
	// 验证流程：
	// 1.从redis中取出验证码
	// 2.验证正确性
	// 3.创建用户，返回token
	code, err := dao.GetVerificationCode(reqDto.Email)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if code != reqDto.Code {
		return erru.New("验证码错误")
	}

	dao.DelVerificationCode(reqDto.Email)

	// 按理来说不可能
	// _, err = dao.GetUserByEmail(reqDto.Email)
	// if err == nil {
	// 	// 如果 err 为 nil，说明找到了用户，邮箱已被注册
	// 	return nil, erru.ErrEmailAlreadyUsed
	// }

	// 加密
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(reqDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return erru.New("加密失败")
	}
	user := &models.User{
		Username: reqDto.UserName,
		Email:    reqDto.Email,
		Password: string(encryptedPwd),
	}

	_, err = dao.CreateUser(user)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	return nil
}

func Login(req *dto.LoginReqDTO) (*models.User, error) {
	user, err := dao.GetUserByIdentifier(req.Identifier)
	if err != nil {
		// 如果错误是 gorm.ErrRecordNotFound，说明用户不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, erru.ErrUserNotFound
		}
		// 其他数据库错误
		return nil, erru.ErrInternalServer.Wrap(err)
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, erru.ErrPasswordIncorrect
	}

	return user, nil
}

func RequestReset(reqDto *dto.RequestResetReqDTO) error {
	_, err := dao.GetUserByEmail(reqDto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return erru.ErrUserNotFound
		}
		return erru.ErrInternalServer.Wrap(err)
	}

	isCool, err := dao.CheckSendCooldown(reqDto.Email)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	if isCool {
		return erru.New("操作太频繁，请稍后再试")
	}

	code := utils.GenerateRandomCode(6)

	err = dao.SetVerifyCode(reqDto.Email, code, 5*time.Minute)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if err := SendResetPasswordMail(reqDto.Email, code); err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	return nil
}

func VerifyReset(reqDto *dto.VerifyResetReqDTO) error {

	// 1.检查是否存在user
	user, err := dao.GetUserByEmail(reqDto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return erru.ErrUserNotFound
		}
		return erru.ErrInternalServer.Wrap(err)
	}
	// 2.检查验证码是否正确
	code, err := dao.GetVerificationCode(reqDto.Email)
	if err != nil || code != reqDto.Code {
		return erru.New("验证码错误")
	}

	// 3.检查密码是否相同
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqDto.Password))
	if err == nil {
		return erru.New("新密码和原密码相同")
	}

	// 4.更新密码
	newPwd, err := bcrypt.GenerateFromPassword([]byte(reqDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return erru.New("加密失败")
	}
	err = dao.UpdateUserPassword(user.ID, string(newPwd))

	return err
}
