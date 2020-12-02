package common

import (
	"net/http"
	"strings"
)

//声明一个函数类型
type  FilterHandle func(rw http.ResponseWriter, req *http.Request)error

type Filter struct {
	//用来存储需要拦截的uri
	filterMap map[string]FilterHandle
}

func NewFilter()*Filter{
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

func (f *Filter)RegisterFilterUri(uri string,handler FilterHandle){
	f.filterMap[uri]=handler
}

func (f *Filter)GetFilterHandle(uri string)FilterHandle{
	return f.filterMap[uri]
}

type webHandle func(rw http.ResponseWriter, req *http.Request)

// Handle 执行拦截器
func(f *Filter)Handle(webHandle webHandle)func(rw http.ResponseWriter, r *http.Request){
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap{
			if strings.Contains(r.RequestURI,path){
				err := handle(rw, r)
				if err != nil {
					_, _ = rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		webHandle(rw,r)
	}
}