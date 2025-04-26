package router

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"

	"hoyang/ownsa/controller"
	"hoyang/ownsa/data/response"
	"hoyang/ownsa/database"
	"hoyang/ownsa/middleware"
	"hoyang/ownsa/repository"
	"hoyang/ownsa/service"
	"hoyang/ownsa/utils"
)

// WebControllerGroup 用于定义各个控制器的集合
type WebControllerGroup struct {
	ControllerUserController   *controller.ControllerUserController   // 用户管理控制器
	CredentialController       *controller.CredentialController       // 认证控制器
	DeviceController           *controller.DeviceController           // 设备控制器
	EventMessageDataController *controller.EventMessageDataController // 事件消息控制器
	PeopleController           *controller.PeopleController           // 人员控制器
	DepartmentController       *controller.DepartmentController       // 部门控制器
}

var WebController *WebControllerGroup // WebControllerGroup 实例
var HttpSrv *http.Server              // HTTP 服务实例

// NewProfileHttpServer 启动一个新的性能分析服务器，监听指定地址
func NewProfileHttpServer(addr string) {
	go func() {
		log.Fatalln(http.ListenAndServe(addr, nil)) // 启动性能分析服务
	}()
}

// CreateHttpServ 创建并启动 HTTP 服务
func CreateHttpServ(confEnv *map[string]string, out *log.Logger) error {
	// 检查运行模式是否为生产模式
	isReleaseMode := (*confEnv)["GIN_MODE"] == gin.ReleaseMode
	if isReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置为生产模式
		gin.DisableConsoleColor()    // 禁用控制台颜色
	}

	// 创建 Gin 引擎实例
	server := gin.New()
	server.Use(gin.LoggerWithWriter(out.Writer()))          // 使用日志记录
	server.Use(gin.CustomRecovery(middleware.ErrorHandler)) // 使用自定义的错误恢复中间件
	server.MaxMultipartMemory = 16 << 20                    // 16MB 限制上传文件的大小

	// 设置跨域请求头
	server.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 对预检请求返回204
			return
		}

		c.Next() // 继续处理请求
	})

	// 配置数据库
	database.SetupDatabase(confEnv, out)

	// 设置路由
	routes := SetupRouter(server)

	// 创建并注册控制器
	CreateWebController()

	// 注册各个控制器的路由
	RegisterControllerUserRoutes(confEnv, routes, WebController.ControllerUserController)
	RegisterPeopleRoutes(confEnv, routes, WebController.PeopleController)
	RegisterDepartmentRoutes(confEnv, routes, WebController.DepartmentController)
	RegisterEventMessageDataRoutes(confEnv, routes, WebController.EventMessageDataController)
	RegisterCredentialRoutes(confEnv, routes, WebController.CredentialController)
	RegisterDeviceRoutes(confEnv, routes, WebController.DeviceController)

	// 配置服务地址，根据平台确定
	servAddr := ":8080"
	if runtime.GOARCH == "arm" && runtime.GOOS == "linux" {
		servAddr = ":80"
	}

	// 创建 HTTP 服务实例
	HttpSrv = &http.Server{
		Addr:           servAddr,         // 服务监听地址
		Handler:        routes,           // 路由
		ReadTimeout:    10 * time.Second, // 读请求超时
		WriteTimeout:   10 * time.Second, // 写请求超时
		MaxHeaderBytes: 1 << 20,          // 请求头最大字节数
	}

	out.Println("CreateHttpServ") // 日志输出

	// 在非生产模式下启动性能分析服务器
	if !isReleaseMode {
		NewProfileHttpServer(":9999")
	}

	// 启动 HTTP 服务
	return HttpSrv.ListenAndServe()
}

// SetupRouter 设置路由
func SetupRouter(service *gin.Engine) *gin.Engine {
	// 静态文件路由
	service.Static("/assets", "./web/assets")
	service.Static("/static", "./web/static")
	service.Static("/audio", "./web/audio")
	service.Static("/images", "./web/images")
	service.Static("/uploads", "./web/uploads")

	// 单个文件路由
	service.StaticFile("/platform-config.json", "./web/platform-config.json")
	service.StaticFile("/version.json", "./web/version.json")
	service.StaticFile("/logo.svg", "./web/logo.svg")
	service.StaticFile("/favicon.ico", "./web/favicon.ico")

	// 加载 HTML 模板
	service.LoadHTMLFiles(
		"../ownsa_web/index.html",
		// "./web/index.html", //部署
	)

	// OWNSA 管理界面首页
	service.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 获取异步路由配置
	service.GET("/api/getAsyncRoutes", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Response{
			Code:    http.StatusOK,
			Success: true,
			Data:    []string{},
		})
	})

	// Ping 测试接口
	service.GET("/ping/", func(c *gin.Context) {
		if runtime.GOOS == "linux" {
			utils.PrintMemUsage()
			var s = utils.GetMemUsageString()
			c.String(http.StatusOK, s)
		} else {
			c.String(http.StatusOK, "pong")
		}
	})

	// 性能分析接口
	service.GET("/prof/", func(c *gin.Context) {
		// 启动内存分析
		f, err := os.OpenFile("/tmp/ownsa.mem.prof", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}

		defer f.Close()

		pprof.Lookup("heap").WriteTo(f, 0)
	})

	// 垃圾回收接口
	service.GET("/gc/", func(c *gin.Context) {
		runtime.GC()
		c.String(http.StatusOK, "OK")
	})

	return service
}

// CreateWebController 创建 Web 控制器实例
func CreateWebController() {
	// 验证器
	validate := validator.New()
	// 创建各个数据仓库
	controllerUserRepository := repository.NewControllerUserRepositoryImpl(database.DB.DbConfig)
	peopleRepository := repository.NewPeopleRepositoryImpl(database.DB.DbCredential)
	departmentRepository := repository.NewDepartmentRepositoryImpl(database.DB.DbCredential)
	eventMessageDataRepository := repository.NewEventMessageDataRepositoryImpl(database.DB.DbEventMessage)
	credentialRepository := repository.NewCredentialRepositoryImpl(database.DB.DbCredential)
	credentialAccessRepository := repository.NewCredentialAccessRepositoryImpl(database.DB.DbCredential)
	doorGroupRepository := repository.NewDoorGroupRepositoryImpl(database.DB.DbOtherGroup)
	interfaceBoardRepository := repository.NewInterfaceBoardRepositoryImpl(database.DB.DbConfig)
	controllerPropRepository := repository.NewControllerPropRepositoryImpl(database.DB.DbConfig)
	cardReaderPropRepository := repository.NewCardReaderPropRepositoryImpl(database.DB.DbConfig)
	inputPropRepository := repository.NewInputPropRepositoryImpl(database.DB.DbConfig)
	outputPropRepository := repository.NewOutputPropRepositoryImpl(database.DB.DbConfig)
	schedGroupRepository := repository.NewSchedGroupRepositoryImpl(database.DB.DbOtherGroup)
	accessGroupRepository := repository.NewAccessGroupRepositoryImpl(database.DB.DbOtherGroup)
	// 创建各个服务实例
	controllerUserService := service.NewControllerUserServiceImpl(
		controllerUserRepository,
		controllerPropRepository,
		validate)
	peopleService := service.NewPeopleServiceImpl(
		peopleRepository,
		credentialRepository,
		credentialAccessRepository,
		doorGroupRepository,
		accessGroupRepository,
		departmentRepository,
		validate)
	departmentService := service.NewDepartmentServiceImpl(
		departmentRepository,
		validate)
	eventMessageDataService := service.NewEventMessageDataServiceImpl(eventMessageDataRepository, validate)
	credentialService := service.NewCredentialServiceImpl(
		credentialRepository,
		credentialAccessRepository,
		doorGroupRepository,
		peopleRepository,
		departmentRepository,
		cardReaderPropRepository,
		outputPropRepository,
		interfaceBoardRepository,
		schedGroupRepository,
		accessGroupRepository,
		validate,
	)
	deviceService := service.NewDeviceServiceImpl(
		interfaceBoardRepository,
		controllerPropRepository,
		controllerUserRepository,
		cardReaderPropRepository,
		inputPropRepository,
		outputPropRepository,
		validate,
	)

	WebController = &WebControllerGroup{}

	WebController.ControllerUserController = controller.NewControllerUserController(controllerUserService)
	WebController.EventMessageDataController = controller.NewEventMessageDataController(eventMessageDataService)
	WebController.PeopleController = controller.NewPeopleController(peopleService)
	WebController.DepartmentController = controller.NewDepartmentController(departmentService)
	WebController.CredentialController = controller.NewCredentialController(credentialService)
	WebController.DeviceController = controller.NewDeviceController(deviceService, controllerUserService)
}

// 注册用户相关的路由
func RegisterControllerUserRoutes(confEnv *map[string]string, service *gin.Engine, controllerUserController *controller.ControllerUserController) {
	// 创建 API 路由组
	router := service.Group("/api")

	// 创建公开和私有的用户路由组
	controllerUserPublicRouter := router.Group("/controllerUser")
	controllerUserPrivateRouter := router.Group("/controllerUser")

	// 公开路由：登录和创建用户
	controllerUserPublicRouter.POST("/login", controllerUserController.Login)
	controllerUserPublicRouter.POST("", controllerUserController.Create)

	// 私有路由：需要身份验证
	controllerUserPrivateRouter.Use(middleware.TokenAuthMiddleware(confEnv))
	{
		// 获取所有用户信息
		controllerUserPrivateRouter.GET("", controllerUserController.FindAll)
		// 根据 ID 获取用户信息
		controllerUserPrivateRouter.GET("/:controllerUserId", controllerUserController.FindById)
		// 更新用户信息
		controllerUserPrivateRouter.PATCH("/:controllerUserId", controllerUserController.Update)
		// 修改密码
		controllerUserPrivateRouter.PATCH("", controllerUserController.ChangePassword)
		// 删除用户
		controllerUserPrivateRouter.DELETE("/:controllerUserId", controllerUserController.Delete)
		// 验证 token
		controllerUserPrivateRouter.GET("/tokenVerify", controllerUserController.TokenVerify)
		// 用户登出
		controllerUserPrivateRouter.GET("/logout", controllerUserController.Logout)
	}
}

// 注册部门相关的路由
func RegisterDepartmentRoutes(confEnv *map[string]string, service *gin.Engine, departmentController *controller.DepartmentController) {
	router := service.Group("/api")
	departmentPrivateRouter := router.Group("/department")

	// 私有路由：需要身份验证
	departmentPrivateRouter.Use(middleware.TokenAuthMiddleware(confEnv))
	{
		// 获取所有部门信息
		departmentPrivateRouter.GET("", departmentController.FindAll)
		// 根据 ID 获取部门信息
		departmentPrivateRouter.GET("/:departmentId", departmentController.FindById)
		// 创建部门
		departmentPrivateRouter.POST("", departmentController.Create)
		// // 更新部门信息
		// departmentPrivateRouter.PATCH("/:departmentId", departmentController.Update)
		// 删除部门
		departmentPrivateRouter.DELETE("/:departmentId", departmentController.Delete)
	}
}

// 注册人员相关的路由
func RegisterPeopleRoutes(confEnv *map[string]string, service *gin.Engine, peopleController *controller.PeopleController) {
	router := service.Group("/api")
	peoplePrivateRouter := router.Group("/people")

	// 私有路由：需要身份验证
	peoplePrivateRouter.Use(middleware.TokenAuthMiddleware(confEnv))
	{
		// 获取所有人员信息
		peoplePrivateRouter.GET("", peopleController.FindAll)
		// 根据 ID 获取人员信息
		peoplePrivateRouter.GET("/:peopleId", peopleController.FindById)
		// 创建人员
		peoplePrivateRouter.POST("", peopleController.Create)
		// 更新人员信息
		peoplePrivateRouter.PATCH("/:peopleId", peopleController.Update)
		// 删除人员
		peoplePrivateRouter.DELETE("/:peopleId", peopleController.Delete)
		// csv批量导入员工信息
		peoplePrivateRouter.POST("/import", peopleController.Import)
	}
}

// 注册事件消息数据相关的路由
func RegisterEventMessageDataRoutes(confEnv *map[string]string, service *gin.Engine, eventMessageDataController *controller.EventMessageDataController) {
	router := service.Group("/api")
	eventMessageDataPrivateRouter := router.Group("/event")

	// 私有路由：需要身份验证
	eventMessageDataPrivateRouter.Use(middleware.TokenAuthMiddleware(confEnv))
	{
		// 获取所有事件消息数据
		eventMessageDataPrivateRouter.GET("", eventMessageDataController.FindAll)
		// 根据消息 ID 获取事件消息数据
		eventMessageDataPrivateRouter.GET("/:messageId", eventMessageDataController.FindById)
		// 同步事件数据
		eventMessageDataPrivateRouter.GET("/sync", eventMessageDataController.Sync)
		// WebSocket 服务
		eventMessageDataPrivateRouter.GET("/ws", eventMessageDataController.WebSocketServer)
		// 新增获取某时间段事件数据的接口
		eventMessageDataPrivateRouter.GET("/peopletime", eventMessageDataController.FindByTimeRange)
	}
}

// 注册凭证相关的路由
func RegisterCredentialRoutes(confEnv *map[string]string, service *gin.Engine, credentialController *controller.CredentialController) {
	router := service.Group("/api")
	credentialPrivateRouter := router.Group("/credential")

	// 私有路由：需要身份验证
	credentialPrivateRouter.Use(middleware.TokenAuthMiddleware(confEnv))
	{
		// 获取所有凭证
		credentialPrivateRouter.GET("", credentialController.FindAll)
		// 根据凭证 ID 获取凭证
		credentialPrivateRouter.GET("/:credentialId", credentialController.FindById)
		// 创建凭证
		credentialPrivateRouter.POST("", credentialController.Create)
		// 更新凭证
		credentialPrivateRouter.PATCH("/:credentialId", credentialController.Update)
		// 删除凭证
		credentialPrivateRouter.DELETE("/:credentialId", credentialController.Delete)
		// 获取门信息
		credentialPrivateRouter.GET("/alldoor", credentialController.AllDoor)
		// 根据凭证 ID 获取门信息
		credentialPrivateRouter.GET("/door/:credentialId", credentialController.ListDoor)
		// 更新门信息
		credentialPrivateRouter.PATCH("/door/:credentialId", credentialController.UpdateDoor)
	}
}

// 注册设备相关的路由
func RegisterDeviceRoutes(confEnv *map[string]string, service *gin.Engine, deviceController *controller.DeviceController) {
	router := service.Group("/api")
	devicePrivateRouter := router.Group("/device")
	devicePublicRouter := router.Group("/device")

	// 公开路由：获取出厂设置
	devicePublicRouter.GET("/factorySet", deviceController.GetFactorySet)

	// 私有路由：需要身份验证
	devicePrivateRouter.Use(middleware.TokenAuthMiddleware(confEnv))
	{
		// 获取控制器属性
		devicePrivateRouter.GET("/controller", deviceController.FindControllerProp)
		// 更新控制器属性
		devicePrivateRouter.POST("/controller", deviceController.UpdateControllerProp)

		// 获取所有 MT2 接口板
		devicePrivateRouter.GET("/mt2InterfaceBoard", deviceController.FindAllMT2InterfaceBoard)
		// 获取所有 MIO 接口板
		devicePrivateRouter.GET("/mioInterfaceBoard", deviceController.FindAllMIOInterfaceBoard)
		// 添加 MT2 接口板
		devicePrivateRouter.POST("/mt2InterfaceBoard", deviceController.AddMT2InterfaceBoard)
		// 更新 MT2 接口板
		devicePrivateRouter.PATCH("/mt2InterfaceBoard/:interfaceBoardId", deviceController.UpdateMT2InterfaceBoard)
		// 添加 MIO 接口板
		devicePrivateRouter.POST("/mioInterfaceBoard", deviceController.AddMIOInterfaceBoard)
		// 更新 MIO 接口板
		devicePrivateRouter.PATCH("/mioInterfaceBoard/:interfaceBoardId", deviceController.UpdateMIOInterfaceBoard)
		// 删除 MT2 接口板
		devicePrivateRouter.DELETE("/mt2InterfaceBoard/:interfaceBoardId", deviceController.DeleteMT2InterfaceBoard)
		// 删除 MIO 接口板
		devicePrivateRouter.DELETE("/mioInterfaceBoard/:interfaceBoardId", deviceController.DeleteMIOInterfaceBoard)
		// 根据接口板 ID 获取 MT2 接口板
		devicePrivateRouter.GET("/mt2InterfaceBoard/:interfaceBoardId", deviceController.FindMT2InterfaceBoardById)
		// 根据接口板 ID 获取 MIO 接口板
		devicePrivateRouter.GET("/mioInterfaceBoard/:interfaceBoardId", deviceController.FindMIOInterfaceBoardById)

		// 同步设备状态
		devicePrivateRouter.GET("/statusSync", deviceController.StatusSync)
		// 开门操作
		devicePrivateRouter.POST("/doorOpen", deviceController.DoorOpen)
		// 火警取消操作
		devicePrivateRouter.POST("/fireCancel", deviceController.FireCancel)

		// 获取系统信息
		devicePrivateRouter.GET("/sysinfo", deviceController.SysInfo)

		// 同步时间
		devicePrivateRouter.POST("/timeSync", deviceController.SyncDatetime)

		// 恢复出厂设置(创建删除与重建数据文件，当系统重启的时候，会删除数据库，进行初始化数据，全部清0)
		devicePrivateRouter.GET("/factoryReset", deviceController.FactoryReset)

		// 升级应用
		devicePrivateRouter.POST("/upgradeApp", deviceController.UpgradeApp)
		// 备份应用
		devicePrivateRouter.GET("/backupApp", deviceController.BackupApp)
		// 恢复应用
		devicePrivateRouter.GET("/restoreApp", deviceController.RestoreApp)
		// 导出应用
		devicePrivateRouter.GET("/exportApp", deviceController.ExportApp)

		// 设置出厂设置(factory账户权限)
		devicePrivateRouter.POST("/factorySet", deviceController.FactorySet)
		// 初始化设备
		devicePrivateRouter.GET("/deviceInit", deviceController.DeviceInit)
		// 重启设备
		devicePrivateRouter.GET("/reboot", deviceController.Reboot)

		// 导出现有数据库
		devicePrivateRouter.POST("/exportData", deviceController.ExportData)
		// 恢复为上传数据
		devicePrivateRouter.POST("/restoreData", deviceController.RestoreData)
		// 导出自动备份的文件夹
		devicePrivateRouter.GET("/openBackupFolder", deviceController.OpenBackupFolder)

	}
}
