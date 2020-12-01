package main

import (
	"errors"
	"fmt"
	"imooc-product/common"
	"imooc-product/encrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}
var localHost = "127.0.0.1"

var port = "8081"
var hashConsistent *common.Consistent

//用来存放控制信息
type AccessControl struct {
	//用来存放用户想要存放的信息
	sourceArray map[int]interface{}
	*sync.RWMutex
}

var accessControl = &AccessControl{
	sourceArray: make(map[int]interface{}),
}

func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourceArray[uid]
	return data
}
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[uid] = "hello"
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	if hostRequest == localHost {
		// 本机校验
		return m.GetDataFromMap(uid.Value)
	} else {
		// 代理校验
		return GetDataFromOtherMap(hostRequest, req )

	}
}

//GetDataFromOtherMap 从其他节点获取结果
func GetDataFromOtherMap(host string, request *http.Request) bool {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/access", nil)
	if err != nil {
		return false
	}
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

//GetDataFromMap 获取本机map 处理业务逻辑 返回结果类型为bool类型
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	record := m.GetNewRecord(uidInt)
	if record != nil {
		return true
	}
	return
}

//执行正常的业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {

	fmt.Println("执行check")
}

//统一验证拦截器 每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证")
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("cookie 获取失败")
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("加密串 sign   cookie 获取失败")
	}
	//对信息进行解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串被篡改")
	}
	fmt.Println("结果开始比对")
	fmt.Println("uid :", uidCookie)
	fmt.Println("解密后用户id：", string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份检验失败")
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	//负载均衡器 设置
	//采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	//添加节点
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
	//1 过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	//2 启动服务
	http.HandleFunc("/check", filter.Handle(Check))

	http.ListenAndServe(":8083", nil)
}
