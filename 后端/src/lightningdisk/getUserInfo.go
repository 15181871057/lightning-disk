package lightningdisk

/*
request
{
	"Method":"getUserInfo",
	"Session":"99EFF439-E07E-79E9-1020-5F5BFB72B129",
	"Parame":{
	}
}
response
{
	"Code": 0,
	"Msg": "获取成功",
	"Ret": {
		"id": 2,
		"name": "Testtestt",
		"registerTime": 1650518554,
		"rootFileId": 7
	}
}
*/
func getUserInfo(request *RequestType, response *ResponseType) {
	userName, err := GetNameByUUID(request.Session)
	if err != nil {
		response.Code = 1
		response.Msg = "登陆状态有误"
		return
	}
	// 搜索账号密码
	db := getDB()
	defer db.Close()
	stmt, err := db.Prepare(`SELECT * FROM users WHERE name=?`)
	if err != nil {
		response.Code = 500
		response.Msg = "数据库异常"
		return
	}
	row := stmt.QueryRow(userName)
	var user userType
	err = row.Scan(&user.Id, &user.Name, &user.Passwd, &user.RegisterTime, &user.RootFileId)
	if err != nil {
		response.Code = 500
		response.Msg = "读取数据时发生错误"
		return
	}
	response.Code = 0
	response.Msg = "获取成功"
	response.Ret["id"] = user.Id
	response.Ret["name"] = user.Name
	response.Ret["registerTime"] = user.RegisterTime
	response.Ret["rootFileId"] = user.RootFileId
}

func init() {
	AddHandle("getUserInfo", getUserInfo, NEEDLOGIN)
}
