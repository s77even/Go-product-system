package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"imooc-product/datamodels"
	"imooc-product/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName string , pwd string)(user *datamodels.User,isOk bool)
	AddUser(user *datamodels.User)(userID int64 ,err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewUserService(repositories repositories.IUserRepository)IUserService{
	return &UserService{UserRepository: repositories}
}

func validatePassword(userPassword string, hashed string)(isOk bool,err error){
	//CompareHashAndPassword Returns nil on success, or an error on failure.
	if err = bcrypt.CompareHashAndPassword([]byte(hashed),[]byte(userPassword));err!=nil{
		return false,errors.New("password error")
	}
	return true,nil
}
func (u *UserService)IsPwdSuccess(userName string , pwd string)(user *datamodels.User,isOk bool){
	user ,err := u.UserRepository.Select(userName)
	if err != nil{
		return
	}
	isOk , _ =validatePassword(pwd,user.HashPassword)
	if !isOk{
		return &datamodels.User{},false
	}
	return
}

func GeneratePassword(userPassword string)([]byte ,error){
	return bcrypt.GenerateFromPassword([]byte(userPassword),bcrypt.DefaultCost)
}

func (u *UserService)AddUser(user *datamodels.User)(userID int64 ,err error){
	pwdByte , errPwd := GeneratePassword(user.HashPassword)
	if errPwd!=nil{
		return userID, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

