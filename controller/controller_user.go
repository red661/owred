package controller

import (
	"errors"
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

// ControllerUserController 结构体定义，用于用户相关的控制器
type ControllerUserController struct {
	controllerUserService service.ControllerUserService // 服务层，包含用户相关操作的业务逻辑
}

// NewControllerUserController 构造函数，初始化用户控制器实例
func NewControllerUserController(service service.ControllerUserService) *ControllerUserController {
	return &ControllerUserController{
		controllerUserService: service,
	}
}

// Create 创建一个新的用户
func (controller *ControllerUserController) Create(ctx *gin.Context) {
	log.Println("create controllerUser") // 日志记录操作类型

	// 解析 JSON 请求体到请求结构体
	createControllerUserRequest := request.CreateControllerUserRequest{}
	err := ctx.ShouldBindJSON(&createControllerUserRequest)
	utils.ErrorPanic(err) // 如果解析失败，抛出错误

	// 打印请求内容，便于调试
	log.Printf("%s", litter.Sdump(createControllerUserRequest))

	// 调用服务层方法创建用户
	createdControllerUser := controller.controllerUserService.Create(createControllerUserRequest)

	// 构造成功响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    createdControllerUser,
	}

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// Update 更新用户信息
func (controller *ControllerUserController) Update(ctx *gin.Context) {
	log.Println("update controllerUser")

	// 解析 JSON 请求体到请求结构体
	updateControllerUserRequest := request.UpdateControllerUserRequest{}
	err := ctx.ShouldBindJSON(&updateControllerUserRequest)
	utils.ErrorPanic(err)
	updateControllerUserRequest.ID = cast.ToUint(ctx.Param("controllerUserId"))
	// 打印请求内容
	log.Printf("%s", litter.Sdump(updateControllerUserRequest))

	// 从 URL 参数中获取用户 ID，并设置到请求结构体中
	controllerUserId := cast.ToUint(ctx.Param("controllerUserId"))
	updateControllerUserRequest.ID = controllerUserId

	// 调用服务层方法更新用户
	controller.controllerUserService.Update(updateControllerUserRequest)

	// 构造响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// ChangePassword 修改用户密码
func (controller *ControllerUserController) ChangePassword(ctx *gin.Context) {
	log.Println("change controllerUser password")

	// 解析 JSON 请求体到请求结构体
	changePasswordRequest := request.ChangePasswordRequest{}
	err := ctx.ShouldBindJSON(&changePasswordRequest)
	utils.ErrorPanic(err)

	// 打印请求内容
	log.Printf("%s", litter.Sdump(changePasswordRequest))

	// 从上下文中获取用户 ID
	controllerUserId, exists := ctx.Get("id")
	if !exists {
		utils.ErrorPanic(errors.New("not authorized")) // 如果未获取到 ID，抛出未授权错误
	}

	// 构造更新密码的请求数据
	controllerUserData := request.UpdateControllerUserRequest{}
	controllerUserData.ID = cast.ToUint(controllerUserId)
	controllerUserData.NewPassword = &changePasswordRequest.NewPassword
	controllerUserData.Password = &changePasswordRequest.Password

	// 调用服务层方法更新密码
	controller.controllerUserService.Update(controllerUserData)

	// 构造响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// Delete 删除用户
func (controller *ControllerUserController) Delete(ctx *gin.Context) {
	log.Println("delete controllerUser")

	// 从 URL 参数中获取用户 ID
	controllerUserId := cast.ToUint(ctx.Param("controllerUserId"))

	// 调用服务层方法删除用户
	controller.controllerUserService.Delete(controllerUserId)

	// 构造响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// FindById 根据 ID 查找用户
func (controller *ControllerUserController) FindById(ctx *gin.Context) {
	log.Println("findby controllerUserid")

	// 从 URL 参数中获取用户 ID
	controllerUserId := cast.ToUint(ctx.Param("controllerUserId"))

	// 调用服务层方法查找用户
	controllerUserResponse := controller.controllerUserService.FindById(controllerUserId)

	// 构造响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    controllerUserResponse,
	}

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// FindAll 查找所有用户
func (controller *ControllerUserController) FindAll(ctx *gin.Context) {
	log.Println("findAll controllerUser")

	// 调用服务层方法获取所有用户
	controllerUserResponse := controller.controllerUserService.FindAll()

	// 构造响应
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    controllerUserResponse,
	}

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// Login 用户登录
func (controller *ControllerUserController) Login(ctx *gin.Context) {
	log.Println("login controllerUser")

	// 构造响应
	webResponse := response.Response{}

	// 解析登录请求体
	loginControllerUserRequest := request.LoginControllerUserRequest{}
	err := ctx.ShouldBindJSON(&loginControllerUserRequest)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	// 调用服务层登录方法
	tokenResponse, err := controller.controllerUserService.Login(loginControllerUserRequest)
	if err != nil {
		webResponse.Code = http.StatusNotFound
		webResponse.Success = false
		webResponse.Message = err.Error()
	} else {
		webResponse.Code = http.StatusOK
		webResponse.Success = true
		webResponse.Data = tokenResponse
	}

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// TokenVerify 验证用户令牌
func (controller *ControllerUserController) TokenVerify(ctx *gin.Context) {
	log.Println("tokenVerify")

	// 构造响应
	webResponse := response.Response{}
	webResponse.Code = http.StatusOK
	webResponse.Success = true

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}

// Logout 用户登出
func (controller *ControllerUserController) Logout(ctx *gin.Context) {
	log.Println("logout controllerUser")

	// 从上下文中获取用户 ID
	tokenUserId, exists := ctx.Get("id")
	if !exists {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("token without user"))
	}

	// 调用服务层方法登出用户
	controller.controllerUserService.Logout(tokenUserId.(uint))

	// 构造响应
	webResponse := response.Response{}
	webResponse.Code = http.StatusOK
	webResponse.Success = true

	// 返回响应
	ctx.JSON(http.StatusOK, webResponse)
}
