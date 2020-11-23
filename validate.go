package main

import (
	"errors"
	"fmt"
	"imooc-product/common"
	"imooc-product/encrypt"
	"net/http"
)

//执行正常的业务逻辑
func Check(w http.ResponseWriter , r *http.Request){

	fmt.Println("执行check")
}
//统一验证拦截器 每个接口都需要提前验证
func Auth(w http.ResponseWriter , r *http.Request)error{
	fmt.Println("执行验证")
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(r *http.Request) error{
	uidCookie , err := r.Cookie("uid")
	if err != nil {
		return errors.New("cookie 获取失败")
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("加密串 sign   cookie 获取失败")
	}
	//对信息进行解密
	signByte ,err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串被篡改")
	}
	fmt.Println("结果开始比对")
	fmt.Println("uid :",uidCookie)
	fmt.Println("解密后用户id：",string(signByte) )
	if checkInfo(uidCookie.Value,string(signByte)){
		return nil
	}
	return errors.New("身份检验失败")
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr==signStr{
		return true
	}
	return false
}

func main(){
	//1 过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check",Auth)
	//2 启动服务
	http.HandleFunc("/check",filter.Handle(Check))

	http.ListenAndServe(":8083",nil)
}
