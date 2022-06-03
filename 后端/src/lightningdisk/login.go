package lightningdisk

import (
	"fmt"
	"regexp"
)

/*
request
{
	"Method":"login",
	"Session":"",
	"Parame":{
		"Name": "Testtestt",
        "Passwd":"Testtestt"
	}
}
response
{
	"Code": 成功返回0, 账号不存在或密码错误返回1, 异常返回2,
	"Msg": "说明",
	"Ret": {
		"Session": "99EFF439-E07E-79E9-1020-5F5BFB72B129"
	}
}
*/
func loginHandle(request *RequestType, response *ResponseType) {
	// 检查账号密码
	name, ok := request.Parame["Name"].(string)
	if !ok {
		response.Code = 400
		response.Msg = "缺少用户名"
		return
	}
	passwd, ok := request.Parame["Passwd"].(string)
	if !ok {
		response.Code = 400
		response.Msg = "缺少密码"
		return
	}
	check, err := regexp.MatchString(`^[a-zA-Z0-9_]{8,16}$`, name)
	if err != nil || !check {
		response.Code = 400
		response.Msg = "账号名格式错误"
		return
	}
	check, err = regexp.MatchString(`^[a-zA-Z0-9_]{8,16}$`, passwd)
	if err != nil || !check {
		response.Code = 400
		response.Msg = "密码格式错误"
		return
	}
	// 搜索账号密码
	db := getDB()
	defer db.Close()
	stmt, err := db.Prepare(`SELECT * FROM users WHERE name=? AND passwd=?`)
	if err != nil {
		response.Code = 500
		response.Msg = "数据库异常"
		return
	}
	row := stmt.QueryRow(name, passwd)
	var user userType
	err = row.Scan(&user.Id, &user.Name, &user.Passwd, &user.RegisterTime, &user.RootFileId)
	if err != nil {
		fmt.Println(err)
		response.Code = 1
		response.Msg = "账号不存在或密码错误"
		return
	}
	loginSetLock.Lock()
	key, ok := checkSet[name]
	if ok {
		delete(loginSet, key)
		delete(checkSet, name)
	}
	var retUUID string
	for {
		retUUID = uuid()
		if _, ok := loginSet[retUUID]; !ok {
			loginSet[retUUID] = name
			checkSet[name] = retUUID
			break
		}
	}
	loginSetLock.Unlock()
	response.Code = 0
	response.Msg = "登陆成功"
	response.Ret["Session"] = retUUID
}

func init() {
	AddHandle("login", loginHandle, 0)
}
