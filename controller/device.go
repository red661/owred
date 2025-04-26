package controller

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sanity-io/litter"
	"github.com/spf13/cast"

	"hoyang/ownsa/data/request"
	"hoyang/ownsa/data/response"
	"hoyang/ownsa/model"
	"hoyang/ownsa/service"
	"hoyang/ownsa/utils"
)

// DeviceController 是一个用于管理设备相关功能的控制器，提供多种接口操作
// 包括设备属性查询、更新，接口板的增删改查等功能。
type DeviceController struct {
	deviceService         service.DeviceService         // 设备相关操作的服务层
	controllerUserService service.ControllerUserService // 用户相关操作的服务层
}

// NewDeviceController 创建一个新的 DeviceController 实例。
func NewDeviceController(
	deviceService service.DeviceService,
	controllerUserService service.ControllerUserService,
) *DeviceController {
	return &DeviceController{
		deviceService:         deviceService,
		controllerUserService: controllerUserService,
	}
}

// FindControllerProp 查询设备的控制器属性信息。
func (controller *DeviceController) FindControllerProp(ctx *gin.Context) {
	log.Println("查询控制器属性信息")

	controllerPropResponse := controller.deviceService.FindControllerProp()

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    controllerPropResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateControllerProp 更新设备的控制器属性。
func (controller *DeviceController) UpdateControllerProp(ctx *gin.Context) {
	log.Println("更新控制器属性")

	var updateControllerPropRequest request.UpdateControllerPropRequest
	err := ctx.ShouldBindJSON(&updateControllerPropRequest)
	utils.ErrorPanic(err)

	log.Printf("更新请求: %s", litter.Sdump(updateControllerPropRequest))

	controllerPropResponse := controller.deviceService.UpdateControllerProp(updateControllerPropRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    controllerPropResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)

	// 配置同步和数据写入磁盘
	ConfigSync()
	syscall.Sync()

	//time.Sleep(3 * time.Second)
	//syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

// FindMT2InterfaceBoardById 根据接口板 ID 查询 MT2 类型接口板信息。
func (controller *DeviceController) FindMT2InterfaceBoardById(ctx *gin.Context) {
	log.Println("通过 ID 查询 MT2 类型接口板")

	interfaceBoardId := cast.ToUint(ctx.Param("interfaceBoardId"))
	deviceResponse := controller.deviceService.FindMT2InterfaceBoardById(interfaceBoardId)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    deviceResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// FindAllMT2InterfaceBoard 查询所有 MT2 类型接口板信息。
func (controller *DeviceController) FindAllMT2InterfaceBoard(ctx *gin.Context) {
	log.Println("查询所有 MT2 类型接口板信息")

	deviceResponse := controller.deviceService.FindAllMT2InterfaceBoard()
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    deviceResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// FindMIOInterfaceBoardById 根据接口板 ID 查询 MIO 类型接口板信息。
func (controller *DeviceController) FindMIOInterfaceBoardById(ctx *gin.Context) {
	log.Println("通过 ID 查询 MIO 类型接口板")

	interfaceBoardId := cast.ToUint(ctx.Param("interfaceBoardId"))
	deviceResponse := controller.deviceService.FindMIOInterfaceBoardById(interfaceBoardId)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    deviceResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// FindAllMIOInterfaceBoard 查询所有 MIO 类型接口板信息。
func (controller *DeviceController) FindAllMIOInterfaceBoard(ctx *gin.Context) {
	log.Println("查询所有 MIO 类型接口板信息")

	deviceResponse := controller.deviceService.FindAllMIOInterfaceBoard()
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    deviceResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// AddMT2InterfaceBoard 添加新的 MT2 类型接口板。
func (controller *DeviceController) AddMT2InterfaceBoard(ctx *gin.Context) {
	log.Println("添加 MT2 类型接口板")

	var createInterfaceBoardRequest request.CreateMT2InterfaceBoardRequest
	err := ctx.ShouldBindJSON(&createInterfaceBoardRequest)
	utils.ErrorPanic(err)

	log.Printf("添加请求: %s", litter.Sdump(createInterfaceBoardRequest))

	interfaceBoardResponse := controller.deviceService.AddMT2InterfaceBoard(createInterfaceBoardRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    interfaceBoardResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
	DataSync()
}

// AddMIOInterfaceBoard 添加新的 MIO 类型接口板。
func (controller *DeviceController) AddMIOInterfaceBoard(ctx *gin.Context) {
	log.Println("添加 MIO 类型接口板")

	var createInterfaceBoardRequest request.CreateMIOInterfaceBoardRequest
	err := ctx.ShouldBindJSON(&createInterfaceBoardRequest)
	utils.ErrorPanic(err)

	log.Printf("添加请求: %s", litter.Sdump(createInterfaceBoardRequest))

	interfaceBoardResponse := controller.deviceService.AddMIOInterfaceBoard(createInterfaceBoardRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    interfaceBoardResponse,
	}

	DataSync()
	ctx.JSON(http.StatusOK, webResponse)

}

// UpdateMT2InterfaceBoard 更新 MT2 类型接口板信息。
func (controller *DeviceController) UpdateMT2InterfaceBoard(ctx *gin.Context) {
	log.Println("更新 MT2 类型接口板信息")

	var updateInterfaceBoardRequest request.UpdateMT2InterfaceBoardRequest
	err := ctx.ShouldBindJSON(&updateInterfaceBoardRequest)
	utils.ErrorPanic(err)

	log.Printf("更新请求: %s", litter.Sdump(updateInterfaceBoardRequest))

	interfaceBoardId := cast.ToUint(ctx.Param("interfaceBoardId"))
	// if interfaceBoardId == 0 {
	// 	panic("内置接口板无法修改")
	// }
	updateInterfaceBoardRequest.IBId = interfaceBoardId
	interfaceBoardResponse := controller.deviceService.UpdateMT2InterfaceBoard(updateInterfaceBoardRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    interfaceBoardResponse,
	}

	DataSync()
	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateMIOInterfaceBoard 更新 MIO 类型接口板信息。
func (controller *DeviceController) UpdateMIOInterfaceBoard(ctx *gin.Context) {
	log.Println("更新 MIO 类型接口板信息")

	updateInterfaceBoardRequest := request.UpdateMIOInterfaceBoardRequest{}
	err := ctx.ShouldBindJSON(&updateInterfaceBoardRequest)
	utils.ErrorPanic(err)

	log.Printf("更新请求: %s", litter.Sdump(updateInterfaceBoardRequest))

	interfaceBoardId := cast.ToUint(ctx.Param("interfaceBoardId"))
	if interfaceBoardId == 0 {
		panic("默认接口板无法修改")
	}
	updateInterfaceBoardRequest.IBId = interfaceBoardId
	interfaceBoardResponse := controller.deviceService.UpdateMIOInterfaceBoard(updateInterfaceBoardRequest)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    interfaceBoardResponse,
	}

	DataSync()
	ctx.JSON(http.StatusOK, webResponse)
}

// DeleteMT2InterfaceBoard 删除指定 ID 的 MT2 类型接口板。
func (controller *DeviceController) DeleteMT2InterfaceBoard(ctx *gin.Context) {
	log.Println("删除 MT2 类型接口板")

	interfaceBoardId := cast.ToUint(ctx.Param("interfaceBoardId"))
	controller.deviceService.DeleteMT2InterfaceBoard(interfaceBoardId)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
	DataSync()
}

// DeleteMIOInterfaceBoard 删除指定 ID 的 MIO 类型接口板。
func (controller *DeviceController) DeleteMIOInterfaceBoard(ctx *gin.Context) {
	log.Println("删除 MIO 类型接口板")

	interfaceBoardId := cast.ToUint(ctx.Param("interfaceBoardId"))
	controller.deviceService.DeleteMT2InterfaceBoard(interfaceBoardId)

	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
	DataSync()
}

func (controller *DeviceController) StatusSync(ctx *gin.Context) {
	log.Println("StatusSync") // 记录日志：同步设备状态

	confEnv, err := godotenv.Read() // 读取 .env 文件，获取环境变量
	if err != nil {
		log.Fatal(err) // 如果读取失败，程序直接退出
	}

	backendBaseURL := confEnv["BackendBaseURL"] // 从环境变量中获取后端服务的基础 URL
	// 创建 HTTP POST 请求，请求后端的 /api/statussync 接口
	req, err := http.NewRequest("POST", backendBaseURL+"api/statussync", nil)
	utils.ErrorPanic(err) // 检查是否发生错误

	req.Header.Set("Content-Type", "application/json") // 设置请求头，表示请求体内容为 JSON 格式

	client := &http.Client{}

	resp, err := client.Do(req) // 发送请求并接收响应
	utils.ErrorPanic(err)
	defer resp.Body.Close() // 确保响应体关闭以释放资源

	body, err := io.ReadAll(resp.Body) // 读取响应体
	utils.ErrorPanic(err)

	log.Printf("StatusSync -> %s", string(body)) // 打印日志，记录响应结果

	//var statusSyncResponse response.StatusSyncResponse
	//err = json.Unmarshal(body, &statusSyncResponse)
	//utils.ErrorPanic(err)

	// 构造返回给前端的响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    string(body),
		//Data: "{\"retcode\":200,\"content\":[{\"ibaddr\":0,\"ibtype\":1,\"ibstate\":1,\"fire\":1,\"box\":0,\"pow1\":0,\"pow2\":0,\"door1\":0,\"door1-timeout\":0,\"door1-forced\":0,\"door1-long\":2,\"door2\":0,\"door2-timeout\":0,\"door2-forced\":0,\"door2-long\":1},{\"ibaddr\":1,\"ibtype\":2,\"ibstate\":1,\"fire\":1,\"box\":0,\"pow1\":0,\"pow2\":0,\"door1\":0,\"door1-timeout\":0,\"door1-forced\":0,\"door1-long\":2,\"door2\":0,\"door2-timeout\":0,\"door2-forced\":0,\"door2-long\":1},{\"ibaddr\":2,\"ibtype\":3,\"ibstate\":1,\"fire\":0,\"box\":0,\"pow1\":0,\"pow2\":0,\"input1\":0,\"input2\":0,\"input3\":0,\"input4\":0,\"input5\":0,\"input6\":0,\"input7\":0,\"input8\":0}]}",
	}

	ctx.JSON(http.StatusOK, webResponse) // 将响应以 JSON 格式返回给前端
}

func (controller *DeviceController) DoorOpen(ctx *gin.Context) {
	// 记录日志：门禁打开请求
	log.Println("DoorOpen")

	// 解析前端传入的 JSON 请求体，绑定到 DoorOpenRequest 结构体
	doorOpenRequest := request.DoorOpenRequest{}
	err := ctx.ShouldBindJSON(&doorOpenRequest)
	utils.ErrorPanic(err)

	log.Printf("DoorOpen <- %s", litter.Sdump(doorOpenRequest))

	// 构造 URL 表单数据
	data := url.Values{}
	data.Set("ibaddr", strconv.FormatUint(cast.ToUint64(doorOpenRequest.IBAddr), 10))
	data.Set("outputaddr", strconv.FormatUint(cast.ToUint64(doorOpenRequest.OutputAddr), 10))
	data.Set("mode", strconv.FormatUint(cast.ToUint64(doorOpenRequest.Mode), 10))
	utils.ErrorPanic(err)

	// 读取 .env 文件，获取环境变量
	confEnv, err := godotenv.Read() // .env in project root path.
	if err != nil {
		log.Fatal(err)
	}

	// 从环境变量中获取后端服务的基础 URL
	backendBaseURL := confEnv["BackendBaseURL"]
	// 创建 HTTP POST 请求，请求后端的 /api/dooropen 接口
	req, err := http.NewRequest("POST", backendBaseURL+"api/dooropen", strings.NewReader(data.Encode()))
	utils.ErrorPanic(err)

	// 设置请求头，表示请求体为表单数据
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	// 发送请求并接收响应
	resp, err := client.Do(req)
	utils.ErrorPanic(err)
	defer resp.Body.Close() // 确保响应体关闭以释放资源

	body, err := io.ReadAll(resp.Body) // 读取响应体
	utils.ErrorPanic(err)
	// 打印日志，记录响应结果
	log.Printf("DoorOpen -> %s", string(body))

	var doorOpenResponse response.DoorOpenResponse
	err = json.Unmarshal(body, &doorOpenResponse)
	utils.ErrorPanic(err)

	// 构造返回给前端的响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    doorOpenResponse,
	}

	ctx.JSON(http.StatusOK, webResponse) // 将响应以 JSON 格式返回给前端
}

// FireCancel 处理设备的取消火灾报警请求。
// 该函数接收一个 gin.Context 参数，用于处理HTTP请求和响应。
// 它从请求体中解析 FireCancelRequest 对象，然后根据请求的数据构建一个 POST 请求，
// 发送到后端服务进行处理。处理结果会被解析并作为 JSON 响应返回给客户端。
func (controller *DeviceController) FireCancel(ctx *gin.Context) {
	log.Println("FireCancel")

	// 解析请求体中的 FireCancelRequest 对象。
	fireCancelRequest := request.FireCancelRequest{}
	err := ctx.ShouldBindJSON(&fireCancelRequest)
	utils.ErrorPanic(err)

	// 打印请求对象的日志。
	log.Printf("FireCancel <- %s", litter.Sdump(fireCancelRequest))

	// 构建要发送到后端服务的请求数据。
	data := url.Values{}
	data.Set("ibaddr", strconv.FormatUint(cast.ToUint64(fireCancelRequest.IBAddr), 10))

	// 加载环境变量配置。
	confEnv, err := godotenv.Read() // .env in project root path.
	if err != nil {
		log.Fatal(err)
	}

	// 构建后端服务的 URL 并创建 HTTP 请求。
	backendBaseURL := confEnv["BackendBaseURL"]
	req, err := http.NewRequest("POST", backendBaseURL+"api/firecancel", strings.NewReader(data.Encode()))
	utils.ErrorPanic(err)

	// 设置请求的 Content-Type 头。
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 创建 HTTP 客户端。
	client := &http.Client{}

	// 发送请求并获取响应。
	resp, err := client.Do(req)
	utils.ErrorPanic(err)
	defer resp.Body.Close()

	// 读取响应体。
	body, err := io.ReadAll(resp.Body)
	utils.ErrorPanic(err)

	// 打印响应体的日志。
	log.Printf("FireCancel -> %s", string(body))

	// 解析响应体为 FireCancelResponse 对象。
	var fireCancelResponse response.FireCancelResponse
	err = json.Unmarshal(body, &fireCancelResponse)
	utils.ErrorPanic(err)

	// 构建并返回 JSON 响应。
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    fireCancelResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// GetFactorySet 处理获取工厂设置的请求。
// 该方法从 deviceService 中检索控制器属性，并将其用于构建工厂设置响应。
// 参数:
//
//	ctx *gin.Context - Gin框架的上下文，用于处理HTTP请求和响应。
func (controller *DeviceController) GetFactorySet(ctx *gin.Context) {
	log.Println("GetFactorySet")

	// 从deviceService中查找控制器属性。
	controllerPropResponse := controller.deviceService.FindControllerProp()

	// 构建工厂设置响应对象，填充从控制器属性中获取的数据。
	factorySetResponse := response.FactorySetResponse{
		BPType:       controllerPropResponse.BPType,
		ProductType:  controllerPropResponse.ProductType,
		IsDoubleLine: controllerPropResponse.IsDoubleLine,
		IsOwnsa:      controllerPropResponse.IsOwnsa,
		IsBakCtl:     controllerPropResponse.IsBakCtl,
	}

	// 构建最终的响应对象，包括状态码、成功标志和数据。
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    factorySetResponse,
	}

	// 以JSON格式发送响应。
	ctx.JSON(http.StatusOK, webResponse)
}

// FactorySet 处理工厂设置请求。
// 该函数接收一个 gin.Context 对象作为参数，从中读取请求数据并进行处理。
// 它首先验证用户是否具有工厂设置的权限，然后根据请求数据更新控制器属性。
// 最后，它返回更新后的控制器属性信息。
func (controller *DeviceController) FactorySet(ctx *gin.Context) {
	log.Println("FactorySet")

	// 解析工厂设置请求体
	factorySetRequest := request.FactorySetRequest{}
	err := ctx.ShouldBindJSON(&factorySetRequest)
	utils.ErrorPanic(err)

	// 验证用户是否已登录并获取用户ID
	uid, exists := ctx.Get("id")
	if !exists {
		panic("Not authorized")
	}

	// 只有 UserType = Factory 的账号才能进行工厂设置
	user := controller.controllerUserService.FindById(uid.(uint))
	if user.UserType != model.UserTypeFactory {
		panic("Permission denied")
	}

	// 将布尔值转换为对应的uint值，用于表示是否启用某些特性
	var is_double_line uint = 0
	if *factorySetRequest.IsDoubleLine {
		is_double_line = 1
	}
	var is_ownsa uint = 0
	if *factorySetRequest.IsOwnsa {
		is_ownsa = 1
	}
	var is_bak_ctl uint = 0
	if *factorySetRequest.IsBakCtl {
		is_bak_ctl = 1
	}

	// 准备更新控制器属性的请求数据
	updateControllerPropRequest := request.UpdateControllerPropRequest{
		BPType:      factorySetRequest.BPType,
		ProductType: factorySetRequest.ProductType,
	}

	updateControllerPropRequest.IsDoubleLine = &is_double_line
	updateControllerPropRequest.IsOwnsa = &is_ownsa
	updateControllerPropRequest.IsBakCtl = &is_bak_ctl

	// 调用服务层方法更新控制器属性并获取响应数据
	controllerPropResponse := controller.deviceService.UpdateControllerProp(updateControllerPropRequest)

	// 构建成功响应并返回给客户端
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    controllerPropResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)

	// 确保数据写入磁盘并同步配置
	syscall.Sync()
	ConfigSync()
}

// SyncDatetime 同步设备时间
// 该方法通过HTTP请求接收时间信息，并将系统时间和硬件时钟同步到请求的时间
func (controller *DeviceController) SyncDatetime(ctx *gin.Context) {
	// 记录方法开始执行的日志
	log.Println("SyncDatetime")

	// 初始化时间同步请求对象
	datetimeSyncRequest := request.DatetimeSyncRequest{}
	// 绑定HTTP请求中的JSON数据到datetimeSyncRequest对象
	err := ctx.ShouldBindJSON(&datetimeSyncRequest)
	// 如果发生错误，立即终止程序
	utils.ErrorPanic(err)

	// 记录请求参数的日志
	log.Printf("SyncDatetime <- %s", litter.Sdump(datetimeSyncRequest))

	// 定义时间格式布局
	layout := time.DateTime
	// 解析请求中的时间字符串为time.Time对象
	dt, err := time.Parse(layout, *datetimeSyncRequest.Datetime)
	// 如果发生错误，立即终止程序
	utils.ErrorPanic(err)

	// 使用解析的时间更新系统时间
	// date -s
	args := []string{"-s", dt.Format(layout)}
	res := exec.Command("date", args...).Run()

	// 同步硬件时钟
	// hwclock --systohc
	args = []string{"-w", dt.Format(layout)}
	exec.Command("hwclock", args...).Run()

	// 构建HTTP响应对象
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    res,
	}

	// 发送HTTP响应
	ctx.JSON(http.StatusOK, webResponse)
}

// FactoryReset 执行设备的出厂重置操作
// 此函数通过创建一个特定的文件来触发系统级别的出厂重置流程，
// 在系统重启后删除数据库，以实现真正的重置。
// 参数:
//
//	ctx *gin.Context: Gin框架的上下文对象，用于处理HTTP请求和响应。
func (controller *DeviceController) FactoryReset(ctx *gin.Context) {
	log.Println("FactoryReset")

	// 创建 SystemResetFactoryFile 文件，重启后将删除数据库，让程序启动自动创建
	conf := utils.GetEnvConf()
	os.Create(conf["SystemResetFactoryFile"])

	// 构建响应对象，表示操作成功
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
	}

	// 发送JSON响应给客户端
	ctx.JSON(http.StatusOK, webResponse)
}

// DeviceInit 初始化设备控制器。
// 该方法响应设备初始化请求，并协调设备服务进行初始化。.
// 参数:
//
//	ctx *gin.Context: Gin框架的上下文，用于处理HTTP请求和响应。
func (controller *DeviceController) DeviceInit(ctx *gin.Context) {
	// 记录设备初始化开始的日志。
	log.Println("DeviceInit")

	// 调用设备服务的初始化方法，并获取初始化响应。
	controllerPropResponse := controller.deviceService.DeviceInit()

	// 构建HTTP响应，包含设备初始化的结果。
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    controllerPropResponse,
	}

	// 以JSON格式返回设备初始化的响应。
	ctx.JSON(http.StatusOK, webResponse)

	// 同步文件系统和配置，确保设备初始化的状态被持久化。
	syscall.Sync()
	ConfigSync()
}

// UpgradeApp 升级应用程序。
// 该方法通过HTTP POST请求接收一个文件，然后将该文件作为应用程序的升级包处理。
// 参数:
//
//	ctx *gin.Context: Gin框架的上下文，用于处理HTTP请求和响应。
func (controller *DeviceController) UpgradeApp(ctx *gin.Context) {
	log.Println("UpgradeApp")

	// 获取上传的文件。
	file, err := ctx.FormFile("file")
	utils.ErrorPanic(err)

	log.Println(file.Filename)

	// 获取环境配置。
	conf := utils.GetEnvConf()
	// 如果存在 conf["AppFile"] 则先备份
	if _, err := os.Stat(conf["AppFile"]); err == nil {
		utils.CopyFile(conf["AppFile"], conf["AppBackupFile"])
		os.Remove(conf["AppFile"])
	}

	// 将上传的文件保存到指定路径。
	ctx.SaveUploadedFile(file, conf["AppFile"])
	// 为升级包文件添加执行权限。
	os.Chmod(conf["AppFile"], 0755)

	// 构建响应对象。
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	// 返回HTTP响应。
	ctx.JSON(http.StatusOK, webResponse)
}

// BackupApp 执行应用程序的备份操作
// 该函数首先检查应用程序文件是否存在，如果存在则进行备份
// 如果备份文件已存在，则先删除旧的备份文件
// 如果应用程序文件不存在，则返回相应的错误信息
func (controller *DeviceController) BackupApp(ctx *gin.Context) {
	log.Println("BackupApp")

	// 获取环境配置
	conf := utils.GetEnvConf()

	// 检查应用程序文件是否存在
	if _, err := os.Stat(conf["AppFile"]); err == nil {
		// 如果备份文件存在则先删除
		if _, err := os.Stat(conf["AppBackupFile"]); err == nil {
			os.Remove(conf["AppBackupFile"])
		}

		// 执行文件备份操作
		err := utils.CopyFile(conf["AppFile"], conf["AppBackupFile"])

		// 构建响应并返回
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: err == nil,
			Data:    nil,
		}

		ctx.JSON(http.StatusOK, webResponse)
	} else if errors.Is(err, os.ErrNotExist) {
		// 如果应用程序文件不存在，返回错误信息
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: false,
			Message: "App not exists.",
		}

		ctx.JSON(http.StatusOK, webResponse)
	} else {
		// 如果发生其他错误，返回通用错误信息
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: false,
			Message: "App not available.",
		}

		ctx.JSON(http.StatusOK, webResponse)
	}
}

// RestoreApp 用于恢复应用程序到先前的版本。
// 该方法首先检查备份文件是否存在，如果存在，则删除当前应用程序文件并用备份文件替换。
// 如果备份文件不存在或不可用，将返回相应的错误信息。
func (controller *DeviceController) RestoreApp(ctx *gin.Context) {
	log.Println("恢复应用程序")

	// 获取环境配置，以便知道备份文件和应用程序文件的路径。
	conf := utils.GetEnvConf()

	// 检查应用程序备份文件是否存在。
	if _, err := os.Stat(conf["AppBackupFile"]); err == nil {
		// 备份文件存在，删除当前应用程序文件。
		os.Remove(conf["AppFile"])
		// 将备份文件重命名为当前应用程序文件。
		os.Rename(conf["AppBackupFile"], conf["AppFile"])
		// 确保应用程序文件具有正确的权限。
		os.Chmod(conf["AppFile"], 0755)

		// 创建一个成功的响应并返回。
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: true,
			Data:    nil,
		}

		ctx.JSON(http.StatusOK, webResponse)
	} else if errors.Is(err, os.ErrNotExist) {
		// 如果错误是因为文件不存在，说明没有先前的版本。
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: false,
			Message: "Previous version not exists.",
		}

		ctx.JSON(http.StatusOK, webResponse)
	} else {
		// 其他错误，可能是权限问题或文件系统问题。
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: false,
			Message: "Previous version not available.",
		}

		ctx.JSON(http.StatusOK, webResponse)
	}
}

// ExportApp 处理应用导出功能。
// 该函数读取应用配置以获取应用文件路径，
// 设置HTTP响应头以便客户端下载文件，
// 如果文件不存在或无法读取，则返回JSON格式的错误消息。
func (controller *DeviceController) ExportApp(ctx *gin.Context) {
	log.Println("ExportApp")

	// 获取环境配置以检索应用文件路径
	conf := utils.GetEnvConf()

	// 检查应用文件是否存在
	if _, err := os.Stat(conf["AppFile"]); err == nil {
		// 获取文件名并设置响应头
		filename := path.Base(conf["AppFile"])
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Transfer-Encoding", "binary")
		ctx.Header("Content-Disposition", "attachment; filename="+filename)
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Length", "0")
		ctx.File(conf["AppFile"])
	} else if errors.Is(err, os.ErrNotExist) {
		// 文件不存在时返回错误响应
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: false,
			Message: "App not exists.",
		}

		ctx.JSON(http.StatusOK, webResponse)
	} else {
		// 其他错误情况返回错误响应
		webResponse := response.Response{
			Code:    http.StatusOK,
			Success: false,
			Message: "App not available.",
		}

		ctx.JSON(http.StatusOK, webResponse)
	}
}

// Reboot 重启系统。
//
// 该方法首先同步文件系统缓冲区缓存与底层存储设备，以确保数据一致性。
// 然后使用 LINUX_REBOOT_CMD_RESTART 命令重启系统。
func (controller *DeviceController) Reboot(ctx *gin.Context) {
	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

// SysInfo 提供系统信息的API接口
// 该方法从设备服务中获取产品类型，读取系统版本信息，获取内核版本和发布日期
// 同时，它还收集设备的序列号，MAC地址和IP地址等信息，并将其返回给调用者
func (controller *DeviceController) SysInfo(ctx *gin.Context) {
	// 初始化系统信息响应对象
	systemInfoResponse := response.SystemInfoResponse{}

	// 读取系统版本信息并解析内核版本
	version, err := os.ReadFile("/proc/version")
	utils.ErrorPanic(err)

	kernelRegex := regexp.MustCompile(`kernel_v[^\s]+`)
	kernelVersion := kernelRegex.FindString(string(version))
	systemInfoResponse.SysVer = kernelVersion

	// 解析系统版本信息中的发布日期
	dateRegex := regexp.MustCompile(`\w{3} \w{3} \d{2} \d{2}:\d{2}:\d{2} UTC \d{4}`)
	dateTime := dateRegex.FindString(string(version))
	systemInfoResponse.ReleaseDate = dateTime

	// 读取设备的序列号信息
	cfg0, err := os.ReadFile("/sys/fsl_otp/HW_OCOTP_CFG0")
	utils.ErrorPanic(err)
	cfg1, err := os.ReadFile("/sys/fsl_otp/HW_OCOTP_CFG1")
	utils.ErrorPanic(err)
	cfg2, err := os.ReadFile("/sys/fsl_otp/HW_OCOTP_CFG2")
	utils.ErrorPanic(err)
	cfg3, err := os.ReadFile("/sys/fsl_otp/HW_OCOTP_CFG3")
	utils.ErrorPanic(err)
	serialNumber := strings.TrimPrefix(string(cfg0), "0x") + strings.TrimPrefix(string(cfg1), "0x") + strings.TrimPrefix(string(cfg2), "0x") + strings.TrimPrefix(string(cfg3), "0x")
	systemInfoResponse.SysSerialNumber = strings.ReplaceAll(serialNumber, "\n", "")

	// 获取网络接口信息，特别是主要的网络接口的MAC地址和IP地址
	interfaces, err := net.Interfaces()
	utils.ErrorPanic(err)
	var primaryInterface net.Interface
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 && !strings.HasPrefix(iface.Name, "can") {
			primaryInterface = iface
			break
		}
	}

	// 如果未找到主要的网络接口，则抛出异常
	if primaryInterface.Name == "" {
		panic("No primary network interface found")
	}
	systemInfoResponse.MACAddr = primaryInterface.HardwareAddr.String()
	addrs, err := primaryInterface.Addrs()
	utils.ErrorPanic(err)
	ips := []string{}
	for _, addr := range addrs {
		ips = append(ips, addr.String())
	}
	systemInfoResponse.IPAddr = strings.Join(ips, ", ")

	// 构建响应对象并返回系统信息
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    systemInfoResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *DeviceController) RestoreData(ctx *gin.Context) {
	log.Println("RestoreData - Receiving uploaded database files")

	// 定义可恢复的数据库文件路径
	filesToRestore := map[string]string{
		"config.db":        "./appdata/db/config.db",
		"iolink.db":        "./appdata/db/iolink.db",
		"credential.db":    "./appdata/db/credential.db",
		"other_group.db":   "./appdata/db/other_group.db",
		"event_message.db": "./appdata/db/event_message.db",

		// // 部署
		// "config.db":        "./ownsadb/config.db",
		// "iolink.db":        "../iolink.db",
		// "credential.db":    "./ownsadb/credential.db",
		// "other_group.db":   "./ownsadb/other_group.db",
		// "event_message.db": "./ownsadb/event_message.db",
	}

	// 解析前端上传的文件
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request format. Expecting multipart/form-data.",
		})
		return
	}

	uploadedFiles := form.File["files"] // 获取前端上传的文件列表
	if len(uploadedFiles) == 0 {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "No files uploaded.",
		})
		return
	}

	// 执行恢复操作
	var failedFiles []string
	for _, fileHeader := range uploadedFiles {
		filename := fileHeader.Filename
		restorePath, exists := filesToRestore[filename]

		// 只允许上传已定义的数据库文件
		if !exists {
			log.Printf("Invalid file upload attempt: %s", filename)
			failedFiles = append(failedFiles, filename)
			continue
		}

		// 打开上传的文件
		srcFile, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open uploaded file %s: %v", filename, err)
			failedFiles = append(failedFiles, filename)
			continue
		}
		defer srcFile.Close()

		// 创建目标文件
		dstFile, err := os.Create(restorePath)
		if err != nil {
			log.Printf("Failed to create destination file %s: %v", restorePath, err)
			failedFiles = append(failedFiles, filename)
			continue
		}

		// 复制文件内容
		_, err = io.Copy(dstFile, srcFile)
		dstFile.Close() // 关闭文件
		if err != nil {
			log.Printf("Failed to copy data for %s: %v", filename, err)
			failedFiles = append(failedFiles, filename)
			continue
		}

		// 设置文件读写权限
		if err := os.Chmod(restorePath, 0666); err != nil {
			log.Printf("Failed to set permissions for %s: %v", restorePath, err)
			failedFiles = append(failedFiles, filename)
			continue
		}

		log.Printf("Successfully restored file: %s", filename)
	}

	// 构建响应
	success := len(failedFiles) == 0
	ctx.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Success: success,
		Message: func() string {
			if success {
				return "Database restore completed successfully."
			}
			return "Database restore completed with some failures."
		}(),
		Data: map[string]interface{}{
			"failedFiles": failedFiles,
		},
	})
}

// ExportData 允许用户下载数据库文件
func (controller *DeviceController) ExportData(ctx *gin.Context) {
	log.Println("ExportData")

	// 获取前端传递的数据库文件名
	var request struct {
		File string `json:"file"` // 需要导出的数据库文件
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request payload.",
		})
		return
	}

	// 允许导出的数据库文件列表
	allowedFiles := map[string]string{
		"config.db":        "./appdata/db/config.db",
		"iolink.db":        "./appdata/db/iolink.db",
		"credential.db":    "./appdata/db/credential.db",
		"other_group.db":   "./appdata/db/other_group.db",
		"event_message.db": "./appdata/db/event_message.db",

		// // 部署
		// "config.db":        "./ownsadb/config.db",
		// "iolink.db":        "../iolink.db",
		// "credential.db":    "./ownsadb/credential.db",
		// "other_group.db":   "./ownsadb/other_group.db",
		// "event_message.db": "./ownsadb/event_message.db",
	}

	// 确保用户请求的文件在允许的文件列表中
	dbFilePath, exists := allowedFiles[request.File]
	if !exists {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid database file request.",
		})
		return
	}

	// 检查数据库文件是否存在
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, response.Response{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Database file not found.",
		})
		return
	}

	// 发送文件给前端
	filename := path.Base(dbFilePath)
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(dbFilePath)
}

func (controller *DeviceController) OpenBackupFolder(ctx *gin.Context) {
	log.Println("OpenBackupFolder")

	// **备份文件夹路径**
	backupFolder := "./appdata/backup"
	// backupFolder := "./backupdb"  //部署

	// **检查目录是否存在**
	if _, err := os.Stat(backupFolder); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"success": false,
			"message": "Backup folder not found.",
		})
		return
	}

	// **创建 ZIP 文件**
	zipFileName := "backup.zip"
	zipFilePath := filepath.Join("./appdata/", zipFileName) // 保存在 `appdata/` 目录
	// zipFilePath := filepath.Join("./backupdb/", zipFileName) //部署
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"success": false,
			"message": "Failed to create ZIP file.",
		})
		return
	}

	// **创建 ZIP Writer**
	zipWriter := zip.NewWriter(zipFile)

	// **遍历备份目录并压缩**
	err = filepath.Walk(backupFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// **计算 ZIP 内的相对路径**
		relPath, err := filepath.Rel(filepath.Dir(backupFolder), path)
		if err != nil {
			return err
		}

		// **如果是目录，创建文件夹**
		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		// **创建 ZIP 内的文件**
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// **打开原始文件**
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// **写入 ZIP**
		_, err = io.Copy(zipEntry, srcFile)
		return err
	})

	// **关闭 ZIP Writer，确保数据写入完成**
	zipWriter.Close()
	zipFile.Close()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"success": false,
			"message": "Failed to add files to ZIP.",
		})
		return
	}

	// **提供 ZIP 文件下载**
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+zipFileName)
	ctx.Header("Content-Type", "application/zip")
	ctx.File(zipFilePath)

	// **可选：下载完成后不删除 ZIP 文件**
	defer os.Remove(zipFilePath) // ⚠️ 如果需要持久存储 ZIP，请注释这行
}
