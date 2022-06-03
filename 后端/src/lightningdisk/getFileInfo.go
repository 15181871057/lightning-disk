package lightningdisk

/*
request
{
	"Method":"getFileInfo",
	"Session":"99EFF439-E07E-79E9-1020-5F5BFB72B129",
	"Parame":{
        "Id": 6
	}
}
response
{
	"Code": 0, 成功返回0,未登录返回1,文件不存在返回2,异常返回3,
	"Msg": "",
	"Ret": {
		"Files": [ 若type为file，则只有该文件，如果为dir，则为多个，第一个保证为所查询的文件
			{
				"Id": 6,
				"Name": "/",
				"Type": "Dir",
				"Size": 0,
				"UpLoadTime": 1650517281
			}
		],
		"Type": "Dir" 目录返回dir,文件返回file,
	}
}
*/
func getFileInfoHandle(request *RequestType, response *ResponseType) {
	_, ok := request.Parame["Id"].(float64)
	if !ok {
		response.Code = 400
		response.Msg = "缺少文件id或者文件id格式错误"
	}
	requestFileId := int(request.Parame["Id"].(float64))
	db := getDB()
	defer db.Close()
	stmt, err := db.Prepare(`SELECT * FROM files WHERE id=?`)
	if err != nil {
		response.Code = 500
		response.Msg = "数据库异常"
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(requestFileId)
	var file File
	err = row.Scan(&file.Id, &file.Name, &file.Size, &file.UpLoadTime, &file.Type)
	if err != nil {
		response.Code = 2
		response.Msg = "文件不存在或数据库异常"
		return
	}
	files := make([]File, 0)
	files = append(files, file)
	if file.Type == "Dir" {
		response.Ret["Type"] = "Dir"
		stmt, err = db.Prepare(`select * from files WHERE id IN (SELECT child FROM fileTree WHERE father=?)`)
		if err != nil {
			response.Code = 3
			response.Msg = "准备搜索文件树时发生异常"
			return
		}
		defer stmt.Close()
		rows, err := stmt.Query(requestFileId)
		if err != nil {
			response.Code = 3
			response.Msg = "搜索文件树时返回异常"
			return
		}
		defer rows.Close()
		for rows.Next() {
			var retFile File
			err = rows.Scan(&retFile.Id, &retFile.Name, &retFile.Size, &retFile.UpLoadTime, &retFile.Type)
			if err != nil {
				response.Code = 3
				response.Msg = "创建搜索结果时发生异常"
				return
			}
			files = append(files, retFile)
		}
	} else {
		response.Ret["Type"] = "File"
	}
	response.Ret["Files"] = files
}

func init() {
	AddHandle("getFileInfo", getFileInfoHandle, NEEDLOGIN)
}
