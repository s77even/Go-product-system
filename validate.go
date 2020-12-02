package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"imooc-product/common"
	"imooc-product/datamodels"
	"imooc-product/encrypt"
	"imooc-product/rabbitmq"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}//手动指定
var localHost = ""// 动态获取
var port = "8083"
//数量控制接口服务器内网IP 或者SLB内网IP
var GetOneIp = "127.0.0.1"
var GetOnePort = "8083"
var hashConsistent *common.Consistent
// rabbitmq
var rabbitMqValidate *rabbitmq.RabbitMQ
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
	hostUrl := "http://"+host+":"+port+"/checkRight"
	response, body, err := GetCurl(hostUrl, request)
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
	//uidInt, err := strconv.Atoi(uid)
	//if err != nil {
	//	return false
	//}
	//record := m.GetNewRecord(uidInt)
	//if record != nil {
	//	return true
	//}
	return true
}

//GetCurl 模拟请求
func GetCurl(hostUrl string,request *http.Request)(response *http.Response,body []byte,err error){
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter,r *http.Request){
	right := accessControl.GetDistributedRight(r)
	if !right{
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}
//执行正常的业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("执行check")
	queryForm , err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"][0])<=0 {
		_, _ = w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)

	// 获取用户cookie
	userCookie , err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false{
		w.Write([]byte("false"))
	}
	//获取数量控制 防止超卖
	hostUrl := "http://"+GetOneIp+":"+GetOnePort+"/getOne"
	responseValidate , validateBody, err := GetCurl(hostUrl,r)
	if err != nil {
		_, _ = w.Write([]byte("false"))
		return
	}
	//判断数量控制接口请求状态
	if responseValidate.StatusCode == 200{
		if string(validateBody)== "true"{
		//整合下单
			productID , err := strconv.ParseInt(productString,10,64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			userID,err := strconv.ParseInt(userCookie.Value,10,64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			//创建消息
			message := datamodels.NewMassage(userID,productID)
			byteMessage , err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
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
	localIp ,err := common.GetIntranceIP()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	//fmt.Println(localHost)

	// rabbitmq creat
	rabbitMqValidate = rabbitmq.NewSimpleRabbitMQ("imoocProduct")
	defer rabbitMqValidate.Destory()

	//1 过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)

	//2 启动服务
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))

	http.ListenAndServe(":8083", nil)
}
