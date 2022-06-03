package lightningdisk

/*
request
{
	"Method":"checkLoginState",
	"Session":"E2D9ED06-D0B7-18D3-DE67-BD3F22BBE08D",
	"Parame":{
	}
}
}
response
{
	"Code":处于登陆状态返回0，未登录返回1,
	"Msg":"说明",
	"Parame":
}
*/
func checkLoginStateHandle(request *RequestType, response *ResponseType) {
	loginSetLock.RLock()
	var isLogin bool
	if _, ok := loginSet[request.Session]; ok {
		isLogin = true
	} else {
		isLogin = false
	}
	loginSetLock.RUnlock()
	if !isLogin {
		response.Code = 1
		response.Msg = "未登录"
		return
	}
	response.Code = 0
	response.Msg = "已登录"
}

func init() {
	AddHandle("checkLoginState", checkLoginStateHandle, 0)
}
