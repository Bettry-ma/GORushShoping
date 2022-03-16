package common

import "net/http"

// FilterHandle 声明一个新的数据类型(函数类型)
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// Filter 拦截器结构体
type Filter struct {
	//用来存储需要拦截的URI
	filterMap map[string]FilterHandle
}

// NewFilter 初始化
func NewFilter() *Filter {
	return &Filter{make(map[string]FilterHandle)}
}

// RegisterFilterUri 注册拦截器
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

// GetFilterHandle 根据URI获取对应的handle
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// WebHandle 声明新的数据类型
type WebHandle func(rw http.ResponseWriter, req *http.Request)

// Handle 执行拦截器,返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			if path == r.RequestURI {
				//执行拦截业务逻辑
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		//执行正常注册的函数
		webHandle(rw, r)
	}
}
