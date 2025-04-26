package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"hoyang/ownsa/data/request"
	"hoyang/ownsa/data/response"
	"hoyang/ownsa/utils"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// PerformSync 执行数据同步操作
// 参数：syncType 同步类型，使用 uint 表示
func PerformSync(syncType uint) {
	log.Println("DataSync") // 记录同步操作日志

	return //部署

	// close .db
	//database.CloseDbConnection()
	// sync

	// 创建 DataSyncRequest 对象，设置同步类型
	dataSyncRequest := request.DataSyncRequest{
		SyncType: syncType,
	}

	// 构造表单数据，将同步类型设置为表单字段
	data := url.Values{}
	data.Set("type", strconv.FormatUint(cast.ToUint64(dataSyncRequest.SyncType), 10))

	// 将请求数据序列化为 JSON 格式，用于调试打印
	jsonData, err := json.Marshal(dataSyncRequest)
	utils.ErrorPanic(err)                          // 处理错误，如果失败则直接终止程序
	log.Printf("DataSync <- %s", string(jsonData)) // 打印同步请求数据

	// 从 .env 文件中读取配置
	confEnv, err := godotenv.Read() // .env 文件位于项目根目录
	if err != nil {
		log.Fatal(err) // 如果读取失败，直接记录并终止程序
	}

	// 从 .env 配置中读取后端基础 URL
	backendBaseURL := confEnv["BackendBaseURL"]

	// 构造 HTTP POST 请求，向后端同步数据
	req, err := http.NewRequest("POST", backendBaseURL+"api/datasync", strings.NewReader(data.Encode()))
	utils.ErrorPanic(err) // 如果请求构造失败，则抛出错误

	// 设置请求头，指定内容类型为表单数据
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	utils.ErrorPanic(err)   // 如果请求发送失败，则抛出错误
	defer resp.Body.Close() // 确保响应体被正确关闭，避免资源泄漏

	// 读取响应体数据
	body, err := io.ReadAll(resp.Body)
	utils.ErrorPanic(err) // 如果读取失败，则抛出错误

	// 将响应体数据反序列化为 DataSyncResponse 结构体
	var dataSyncResponse response.DataSyncResponse
	err = json.Unmarshal(body, &dataSyncResponse)
	utils.ErrorPanic(err) // 如果反序列化失败，则抛出错误

	// 打印同步响应数据
	log.Printf("DataSync -> %v+", dataSyncResponse)

	// exit
	//HttpSrv.Close()
}

// ConfigSync 执行配置同步操作
func ConfigSync() {
	// 调用 PerformSync，传入同步类型 3（表示配置同步）
	PerformSync(3)
}

// DataSync 执行数据同步操作
func DataSync() {
	// 调用 PerformSync，传入同步类型 1（表示数据同步）
	PerformSync(1)
}
