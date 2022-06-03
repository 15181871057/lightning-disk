package lightningdisk

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// uuid
func uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

//isLoginUUID

var loginSetLock sync.RWMutex

// uuidToName
var loginSet map[string]string = map[string]string{}

// nameToUUID
var checkSet map[string]string = map[string]string{}

func GetNameByUUID(UUID string) (string, error) {
	loginSetLock.RLock()
	userName, ok := loginSet[UUID]
	loginSetLock.RUnlock()
	if ok {
		return userName, nil
	}
	return "", errors.New("UUID不存在")
}

// DB

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Dx990522@/test")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// InitMySql
func initSql() {
	db := getDB()
	defer db.Close()
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users(
		id INT AUTO_INCREMENT,
		name VARCHAR(32) NOT NULL UNIQUE,
		passwd VARCHAR(32) NOT NULL,
		registerTime int NOT NULL,
		rootFileId int NOT NULL,
		PRIMARY KEY(id, name)
	);
	`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS fileType(
		id INT AUTO_INCREMENT PRIMARY KEY,
		typeName VARCHAR(128) NOT NULL UNIQUE
	)
	`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS files(
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(128) NOT NULL,
		size INT NOT NULL,
		upLoadTime INT NOT NULL,
		fileType VARCHAR(128) NOT NULL,
		owner int NOT NULL,
		FOREIGN KEY(fileType) REFERENCES fileType(typeName),
		FOREIGN KEY(owner) REFERENCES users(id)
	);
	`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS fileTree(
		father INT NOT NULL,
		child INT PRIMARY KEY NOT NULL,
		FOREIGN KEY(father) REFERENCES files(id),
		FOREIGN KEY(child) REFERENCES files(id)
	);
	`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS upload(
		id int PRIMARY KEY AUTO_INCREMENT,
		fileId int NOT NULL,
		size int NOT NULL,
		time int NOT NULL,
		FOREIGN KEY(fileId) REFERENCES files(id)
	);
	`)
	if err != nil {
		panic(err)
	}
}

// Handle
const (
	NEEDLOGIN int64 = 1 << 0
)

func checkIsLogin(request *RequestType) bool {
	// DEBUG
	return true
	session := request.Session
	if len(session) == 0 {
		return false
	}
	loginSetLock.RLock()
	var isLogin bool
	if _, ok := loginSet[request.Session]; ok {
		isLogin = true
	} else {
		isLogin = false
	}
	loginSetLock.RUnlock()
	return isLogin
}

type RequestHandle struct {
	handle  func(*RequestType, *ResponseType)
	flagSet int64
}

var handleMap map[string]RequestHandle = map[string]RequestHandle{}

func routeHandle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	url := req.URL.Path
	if req.Method != "POST" || req.URL.Path != "/api" {
		http.NotFound(w, req)
		return
	}
	fmt.Println(url)
	buf := make([]byte, 4096)
	jsonData := make([]byte, 0)
	for {
		size, err := req.Body.Read(buf)
		jsonData = append(jsonData, buf[:size]...)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
	}
	var request RequestType
	var response ResponseType = ResponseType{
		Code: 0, Msg: "",
		Ret: make(map[string]interface{})}
	err := json.Unmarshal(jsonData, &request)
	fmt.Println(string(jsonData))
	fmt.Println(err)
	if err != nil {
		response.Code = 404
		response.Msg = "格式错误"
		goto returnMark
	} else {
		handle, ok := handleMap[request.Method]
		if !ok {
			response.Code = 404
			response.Msg = "Mathod Not Found"
			goto returnMark
		}
		if handle.flagSet&NEEDLOGIN != 0 && !checkIsLogin(&request) {
			response.Code = 401
			response.Msg = "未登录"
			goto returnMark
		}
		handle.handle(&request, &response)
	}
returnMark:
	index := 0
	writeBuf, _ := json.Marshal(&response)
	for index < len(writeBuf) {
		size, _ := w.Write(writeBuf[index:])
		index += size
	}
}

func AddHandle(re string, handle func(*RequestType, *ResponseType), flagSet int64) {
	handleMap[re] = RequestHandle{handle: handle, flagSet: flagSet}
}

/*
request
{
	"session":uuid，
	"filePath":"文件路径"
}
response
{
	"code":成功返回0,未登录返回1,文件已存在返回2,异常返回3,
	"path":"ws路径",
}
*/
func getUpLoadFileWSHandle(request *RequestType, response *ResponseType) {
}

/*
request
{
	"session":uuid，
	"filePath":"文件路径"
}
response
{
	"code":成功返回0,未登录返回1,文件不存在返回2,异常返回3,
	"path":"ws路径",
}
*/
func getDownLoadFileWSHandle(request *RequestType, response *ResponseType) {
}

// WebSocketHandle

func Start(address string) {
	initSql()
	AddHandle("getUpLoadFileWS", getUpLoadFileWSHandle, NEEDLOGIN)
	AddHandle("getDownLoadFileWS", getDownLoadFileWSHandle, NEEDLOGIN)
	http.HandleFunc("/", routeHandle)
	http.ListenAndServe(address, nil)
}
