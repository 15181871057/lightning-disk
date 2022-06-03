package lightningdisk

/*
request
{
	"Method":"deleteFile",
	"Session":"99EFF439-E07E-79E9-1020-5F5BFB72B129",
	"Parame":{
        "Id": 4 文件id
	}
}
response
{
	{
	"Code": 0, 不管是否存在都返回0
	"Msg": "删除成功", 说明
	"Ret": {}
}
*/
func deleteFileHandle(request *RequestType, response *ResponseType) {
	_, ok := request.Parame["Id"]
	if !ok {
		response.Code = 400
		response.Msg = "缺少参数"
		return
	}
	needDeleteFileId := int(request.Parame["Id"].(float64))
	db := getDB()
	defer db.Close()
	stmt, err := db.Prepare(`
		WITH RECURSIVE cte_select(id) AS (
		SELECT father from fileTree WHERE fileTree.father=?
		UNION ALL
		select f.child from fileTree AS f JOIN cte_select AS c ON f.father=c.id)
		DELETE fileTree, Files FROM fileTree JOIN  Files WHERE father in (SELECT * FROM cte_select);
		`)
	if err != nil {
		response.Code = 500
		response.Msg = "数据库准备时发生异常"
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(needDeleteFileId)
	if err != nil {
		response.Code = 500
		response.Msg = "数据库删除时发生异常"
		return
	}
	response.Code = 0
	response.Msg = "删除成功"
}

func init() {
	AddHandle("deleteFile", deleteFileHandle, NEEDLOGIN)
}
