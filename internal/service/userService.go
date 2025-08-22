package service

import (
	"Nuxus/configs"
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

type UserService struct {
	userDAO      *dao.UserDAO
	redisClient  *dao.RedisClient
	emailService *EmailService
	config       *configs.Config
}

func NewUserService(userDAO *dao.UserDAO, redisClient *dao.RedisClient, emailService *EmailService, config *configs.Config) *UserService {
	return &UserService{
		userDAO:      userDAO,
		redisClient:  redisClient,
		emailService: emailService,
		config:       config,
	}
}

func (us *UserService) Register(reqDto *dto.RegisterReqDTO) error {
	// 注册流程：
	// 1. 检查邮箱用户是否存在
	// 2. 生成验证码，放在redis里
	// 3. 发送验证码
	_, err := us.userDAO.GetUserByEmail(reqDto.Email)
	if err == nil {
		return erru.ErrEmailAlreadyUsed.Wrap(err)
	}

	isCool, err := us.redisClient.CheckSendCooldown(reqDto.Email)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	if isCool {
		return erru.New("操作太频繁，请稍后再试")
	}

	code := utils.GenerateRandomCode(6)

	err = us.redisClient.SetVerifyCode(reqDto.Email, code, 5*time.Minute)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if err := us.emailService.SendRegisterMail(reqDto.Email, code); err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	return nil
}

func (us *UserService) VerifyRegister(reqDto *dto.VerifyRegisterReqDTO) error {
	// 验证流程：
	// 1.从redis中取出验证码
	// 2.验证正确性
	// 3.创建用户，返回token
	code, err := us.redisClient.GetVerificationCode(reqDto.Email)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if code != reqDto.Code {
		return erru.New("验证码错误")
	}

	us.redisClient.DelVerificationCode(reqDto.Email)

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

	_, err = us.userDAO.CreateUser(user)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	return nil
}

func (us *UserService) Login(req *dto.LoginReqDTO) (*models.User, error) {
	user, err := us.userDAO.GetUserByIdentifier(req.Identifier)
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

func (us *UserService) RequestReset(reqDto *dto.RequestResetReqDTO) error {
	_, err := us.userDAO.GetUserByEmail(reqDto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return erru.ErrUserNotFound
		}
		return erru.ErrInternalServer.Wrap(err)
	}

	isCool, err := us.redisClient.CheckSendCooldown(reqDto.Email)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}
	if isCool {
		return erru.New("操作太频繁，请稍后再试")
	}

	code := utils.GenerateRandomCode(6)

	err = us.redisClient.SetVerifyCode(reqDto.Email, code, 5*time.Minute)
	if err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	if err := us.emailService.SendResetPasswordMail(reqDto.Email, code); err != nil {
		return erru.ErrInternalServer.Wrap(err)
	}

	return nil
}

func (us *UserService) VerifyReset(reqDto *dto.VerifyResetReqDTO) error {

	// 1.检查是否存在user
	user, err := us.userDAO.GetUserByEmail(reqDto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return erru.ErrUserNotFound
		}
		return erru.ErrInternalServer.Wrap(err)
	}
	// 2.检查验证码是否正确
	code, err := us.redisClient.GetVerificationCode(reqDto.Email)
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
	err = us.userDAO.UpdateUserPassword(user.ID, string(newPwd))

	return err
}
