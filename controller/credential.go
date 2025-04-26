package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
	"github.com/spf13/cast"

	"hoyang/ownsa/data/request"
	"hoyang/ownsa/data/response"
	"hoyang/ownsa/service"
	"hoyang/ownsa/utils"
)

// CredentialController 用于管理凭证（Credential）的控制器
type CredentialController struct {
	credentialService service.CredentialService // 依赖 CredentialService 提供具体业务逻辑
}

// NewCredentialController 创建一个新的 CredentialController 实例
// 参数：service 提供业务逻辑的 CredentialService
func NewCredentialController(service service.CredentialService) *CredentialController {
	return &CredentialController{
		credentialService: service,
	}
}

// Create 处理创建凭证的 HTTP 请求
// 路由：POST /credentials
func (controller *CredentialController) Create(ctx *gin.Context) {
	log.Println("create credential") // 记录日志

	// 解析请求体中的 JSON 数据到 CreateCredentialRequest 结构体
	createCredentialRequest := request.CreateCredentialRequest{}
	err := ctx.ShouldBindJSON(&createCredentialRequest)
	utils.ErrorPanic(err) // 如果解析失败，抛出错误

	log.Printf("%s", litter.Sdump(createCredentialRequest)) // 打印请求数据，用于调试

	// 调用业务层的 Create 方法，返回创建的凭证数据
	credentialResponse := controller.credentialService.Create(createCredentialRequest)

	// 构造标准化的 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    credentialResponse,
	}

	// 返回 JSON 响应
	ctx.JSON(http.StatusOK, webResponse)

	// 触发数据同步
	DataSync()
}

// Update 处理更新凭证的 HTTP 请求
// 路由：PUT /credentials/:credentialId
func (controller *CredentialController) Update(ctx *gin.Context) {
	log.Println("update credential")

	// 解析请求体中的 JSON 数据到 UpdateCredentialRequest 结构体
	updateCredentialRequest := request.UpdateCredentialRequest{}
	err := ctx.ShouldBindJSON(&updateCredentialRequest)
	utils.ErrorPanic(err)

	log.Printf("%s", litter.Sdump(updateCredentialRequest)) // 打印请求数据，用于调试

	// 从 URL 参数中获取 credentialId，并将其转换为 uint 类型
	credentialId := cast.ToUint(ctx.Param("credentialId"))
	updateCredentialRequest.UniqueId = uint(credentialId)

	// 调用业务层的 Update 方法更新凭证
	controller.credentialService.Update(updateCredentialRequest)

	// 构造标准化的 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	// 返回 JSON 响应
	ctx.JSON(http.StatusOK, webResponse)

	// 触发数据同步
	DataSync()
}

// Delete 处理删除凭证的 HTTP 请求
// 路由：DELETE /credentials/:credentialId
func (controller *CredentialController) Delete(ctx *gin.Context) {
	log.Println("delete credential")

	// 从 URL 参数中获取 credentialId，并将其转换为 uint 类型
	credentialId := cast.ToUint(ctx.Param("credentialId"))

	// 调用业务层的 Delete 方法删除凭证
	controller.credentialService.Delete(credentialId)

	// 构造标准化的 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	// 返回 JSON 响应
	ctx.JSON(http.StatusOK, webResponse)

	// 触发数据同步
	DataSync()
}

// FindById 根据凭证 ID 查找凭证信息
// 路由：GET /credentials/:credentialId
func (controller *CredentialController) FindById(ctx *gin.Context) {
	log.Println("findby credentialId")

	// 从 URL 参数中获取 credentialId，并将其转换为 uint 类型
	credentialId := cast.ToUint(ctx.Param("credentialId"))

	// 调用业务层的 FindById 方法获取凭证信息
	credentialResponse := controller.credentialService.FindById(credentialId)

	// 构造标准化的 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    credentialResponse,
	}

	// 返回 JSON 响应
	ctx.JSON(http.StatusOK, webResponse)
}

// FindAll 获取所有凭证信息
// 路由：GET /credentials
func (controller *CredentialController) FindAll(ctx *gin.Context) {
	log.Println("findAll credential")

	// 调用业务层的 FindAll 方法获取所有凭证信息
	credentialResponse := controller.credentialService.FindAll()

	// 构造标准化的 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    credentialResponse,
	}

	// 返回 JSON 响应
	ctx.JSON(http.StatusOK, webResponse)
}

// ListDoor 获取凭证关联的门信息
// 路由：GET /credentials/:credentialId/doors
func (controller *CredentialController) ListDoor(ctx *gin.Context) {
	log.Println("list door")

	// 从 URL 参数中获取 credentialId，并将其转换为 uint 类型
	credentialId := cast.ToUint(ctx.Param("credentialId"))

	// 调用业务层的 ListDoor 方法获取门信息
	doorResponse := controller.credentialService.ListDoor(credentialId)

	// 构造标准化的 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    doorResponse,
	}

	// 返回 JSON 响应
	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *CredentialController) AllDoor(ctx *gin.Context) {
	log.Println("Fetching all doors")

	// 调用业务层获取所有门信息
	allDoors := controller.credentialService.AllDoor()

	// 组织 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    allDoors,
	}

	// 返回 JSON 结果
	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateDoor 更新凭证关联的门信息
// 路由：PUT /credentials/:credentialId/doors
func (controller *CredentialController) UpdateDoor(ctx *gin.Context) {
	log.Println("update door")

	// 解析请求体中的 JSON 数据到 UpdateCredentialDoorRequest 结构体
	updateCredentialDoorRequest := request.UpdateCredentialDoorRequest{}
	err := ctx.ShouldBindJSON(&updateCredentialDoorRequest)
	utils.ErrorPanic(err)

	// 从 URL 参数中获取 credentialId，并将其赋值给请求结构体
	credentialId := cast.ToUint(ctx.Param("credentialId"))
	updateCredentialDoorRequest.CredentialId = &credentialId

	log.Printf("%s", litter.Sdump(updateCredentialDoorRequest)) // 打印请求数据，用于调试

	// 调用业务层的 UpdateDoor 方法更新门信息
	doorResponse := controller.credentialService.UpdateDoor(updateCredentialDoorRequest)

	// 构造标准化的 HTTP 响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    doorResponse,
	}

	// 返回 JSON 响应
	ctx.JSON(http.StatusOK, webResponse)

	// 触发数据同步
	DataSync()
}
