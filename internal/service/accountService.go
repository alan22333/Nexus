package service

import (
	"Nuxus/configs"
	"Nuxus/internal/dao"
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
	"Nuxus/pkg/erru"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"gorm.io/gorm"
)

type AccountService struct {
	userDAO *dao.UserDAO
	config *configs.Config
}

func NewAccountService(userDAO *dao.UserDAO,config *configs.Config) *AccountService {
	return &AccountService{
		userDAO: userDAO,
		config: config,
	}
}

func (a *AccountService) GetProfile(userId uint) (*models.User, error) {
	user, err := a.userDAO.GetUserById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, erru.ErrResourceNotFound
		}
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	return user, nil
}

func (a *AccountService) UpdateProfile(userId uint, reqDto dto.ProfileReqDTO) (*models.User, error) {
	user, err := a.userDAO.GetUserById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, erru.ErrResourceNotFound
		}
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	a.updateUser(user, &reqDto)
	user, err = a.userDAO.UpdateProfile(user)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	return user, nil
}

func (a *AccountService) updateUser(user *models.User, reqDto *dto.ProfileReqDTO) {
	user.Gender = reqDto.Gender
	user.Phone = reqDto.Phone
	user.QQ = reqDto.QQ
	user.Wechat = reqDto.Wechat
	user.Bio = reqDto.Bio

	user.IsPhonePublic = reqDto.Privacy.IsPhonePublic
	user.IsEmailPublic = reqDto.Privacy.IsEmailPublic
	user.IsQQPublic = reqDto.Privacy.IsQQPublic
	user.IsWechatPublic = reqDto.Privacy.IsWechatPublic
	user.IsGenderPublic = reqDto.Privacy.IsGenderPublic
}

// ---------------------头像------------------------------
// getQiniuZone 根据配置字符串返回七牛云的 Zone 对象
func getQiniuZone(zoneStr string) *storage.Zone {
	switch zoneStr {
	case "ZoneHuadong":
		return &storage.ZoneHuadong
	case "ZoneHuabei":
		return &storage.ZoneHuabei
	case "ZoneHuanan":
		return &storage.ZoneHuanan
	case "ZoneBeimei":
		return &storage.ZoneBeimei
	case "ZoneXinjiapo":
		return &storage.ZoneXinjiapo
	default:
		return &storage.ZoneHuadong // 默认华东
	}
}

// UpdateAvatar 使用七牛云 Kodo 处理头像上传
func (a *AccountService) UpdateAvatar(userID uint, file *multipart.FileHeader) (string, error) {
	qiniuConf := a.config.Qiniu

	// 1. 生成上传凭证
	putPolicy := storage.PutPolicy{
		Scope: qiniuConf.Bucket,
	}
	mac := qbox.NewMac(qiniuConf.AccessKey, qiniuConf.SecretKey)
	uploadToken := putPolicy.UploadToken(mac)

	// 2. 配置存储区域
	cfg := storage.Config{
		Zone:     getQiniuZone(qiniuConf.Zone),
		UseHTTPS: false, // 根据你的域名是否支持 HTTPS 修改
	}

	// 3. 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	// 4. 生成唯一的对象键 (文件名)
	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("avatars/%d/%s%s", userID, uuid.New().String(), ext)

	// 5. 打开上传的文件流
	src, err := file.Open()
	if err != nil {
		return "", erru.ErrInternalServer.Wrap(err)
	}
	defer src.Close()

	// 6. 执行上传
	err = formUploader.Put(context.Background(), &ret, uploadToken, key, src, file.Size, nil)
	if err != nil {
		return "", erru.ErrInternalServer.Wrap(err)
	}

	// 7. 构建公开访问的 URL
	avatarURL := fmt.Sprintf("http://%s/%s", qiniuConf.Domain, ret.Key)

	// 8. 将新的 URL 更新到数据库
	if err := a.userDAO.UpdateUserAvatar(userID, avatarURL); err != nil {
		return "", erru.ErrInternalServer.Wrap(err)
	}

	return avatarURL, nil
}
