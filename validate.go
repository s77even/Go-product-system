package main

import (
	"fmt"
	"imooc-product/common"
	"net/http"
)

//执行正常的业务逻辑
func Check(w http.ResponseWriter , r *http.Request){

	fmt.Println("执行check")
}
//统一验证拦截器 每个接口都需要提前验证
func Auth(w http.ResponseWriter , r *http.Request)error{
	fmt.Println("执行验证")
	return nil 
	//return errors.New("验证失败！")
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
