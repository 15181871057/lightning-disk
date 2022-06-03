package lightningdisk

import (
	"regexp"
	"time"
)

/*
request
{
	"Method":"register",
	"Session":"",
	"Parame":{
		"Name": "Testtest",
        "Passwd":"Testtest"
	}
}
response
{
	"Code": 0,
	"Msg": 成功返回0, 账号或密码格式不正确1，账号已存在返回2，缺少参数返回3，异常返回4,
	"Ret": {}
}
*/
func registerHandle(request *RequestType, response *ResponseType) {
	// 输入解析
	name, ok := request.Parame["Name"].(string)
	if !ok {
		response.Code = 400
		return
	}
	passwd, ok := request.Parame["Passwd"].(string)
	if !ok {
		response.Code = 400
		return
	}
	check, err := regexp.MatchString(`^[a-zA-Z0-9_]{8,16}$`, name)
	if err != nil || !check {
		response.Code = 1
		response.Msg = "账号名或密码格式错误"
		return
	}
	check, err = regexp.MatchString(`^[a-zA-Z0-9_]{8,16}$`, passwd)
	if err != nil || !check {
		response.Code = 1
		response.Msg = "账号名或密码格式错误"
		return
	}
	// 插入数据库
	db := getDB()
	defer db.Close()
	tx, err := db.Begin()
	registerTime := int(time.Now().Unix())
	if err != nil {
		if err != nil {
			response.Code = 500
			response.Msg = "事务创建失败"
			return
		}
	}
	result, err := tx.Exec(
		`INSERT INTO files VALUES(null, "/", 0, ?, "Dir")`, registerTime)
	if err != nil {
		response.Code = 2
		response.Msg = "创建用户根目录失败"
		tx.Rollback()
		return
	}
	userRootFileId, err := result.LastInsertId()
	if err != nil {
		response.Code = 2
		response.Msg = "获取根目录id失败"
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`INSERT INTO users VALUES(null, ?, ?, ?, ?)`,
		name, passwd, registerTime, userRootFileId)
	if err != nil {
		response.Code = 2
		response.Msg = "账号已存在"
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		response.Code = 500
		response.Msg = "事务提交失败"
		return
	}
	response.Code = 0
	response.Msg = "注册成功"
}

func init() {
	AddHandle("register", registerHandle, 0)
}
