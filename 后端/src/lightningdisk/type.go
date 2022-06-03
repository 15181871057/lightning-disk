package lightningdisk

type RequestType struct {
	Method  string
	Session string
	Parame  map[string]interface{}
}

/*
通用Code
400 缺少参数/参数类型错误/参数格式错误
401 未登录
404 路径错误或请求的方法不存在
500 服务器处理中发生异常
*/
type ResponseType struct {
	Code int
	Msg  string
	Ret  map[string]interface{}
}

type userType struct {
	Id           int
	Name         string
	Passwd       string
	RegisterTime int
	RootFileId   int
}

type File struct {
	Id         int
	Name       string
	Type       string
	Size       int
	UpLoadTime int
}
