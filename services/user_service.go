package services

import (
	"GORushShoping/datamodels"
	"GORushShoping/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

// ValidatePassWord 用于将用户密码与进行哈希过的密码进行比对,验证密码是否正确
func ValidatePassWord(userPassword, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("hashedPassword compare not match")
	}
	return true, nil
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	var err error
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, err = ValidatePassWord(pwd, user.HashPassword)
	if isOk {
		return
	}
	return &datamodels.User{}, false
}

// GeneratePassword 将用户的明文密码使用哈希加密
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errGen := GeneratePassword(user.HashPassword)
	if errGen != nil {
		return userId, errGen
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func NewUserService(up repositories.IUserRepository) IUserService {
	return &UserService{up}
}