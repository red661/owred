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

// DepartmentController 用于管理 Department 实体的控制器
type DepartmentController struct {
	departmentService service.DepartmentService // 依赖的服务层，用于处理 Department 数据的业务逻辑
}

// NewDepartmentController 创建并返回一个新的 DepartmentController 实例
func NewDepartmentController(service service.DepartmentService) *DepartmentController {
	return &DepartmentController{
		departmentService: service,
	}
}

// Create 创建一个新的 Department 实体
func (controller *DepartmentController) Create(ctx *gin.Context) {
	log.Println("create department")

	// 解析并绑定请求体到 CreateDepartmentRequest 结构体
	createDepartmentRequest := request.CreateDepartmentRequest{}
	err := ctx.ShouldBindJSON(&createDepartmentRequest)
	utils.ErrorPanic(err) // 如果解析失败则抛出错误

	// 使用 litter 库格式化打印请求内容，方便调试
	log.Printf("%s", litter.Sdump(createDepartmentRequest))

	// 调用服务层方法创建 Department 实体
	controller.departmentService.Create(createDepartmentRequest)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)

	// 调用同步方法（如果需要）
	DataSync()
}

// Update 更新指定的 Department 实体
func (controller *DepartmentController) Update(ctx *gin.Context) {
	log.Println("update department")

}

// Delete 删除指定的 Department 实体
func (controller *DepartmentController) Delete(ctx *gin.Context) {
	log.Println("delete department")

	// 从 URL 参数中获取 departmentId
	departmentId := cast.ToUint(ctx.Param("departmentId"))

	// 调用服务层方法删除指定的 Department 实体
	controller.departmentService.Delete(departmentId)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)

	// 调用同步方法（如果需要）
	DataSync()
}

// FindById 根据 ID 查询指定的 Department 实体
func (controller *DepartmentController) FindById(ctx *gin.Context) {
	log.Println("find department by ID")

	// 从 URL 参数中获取 departmentId
	departmentId := cast.ToUint(ctx.Param("departmentId"))

	// 调用服务层方法查询指定 ID 的 Department 实体
	departmentResponse := controller.departmentService.FindById(departmentId)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    departmentResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// FindAll 查询所有的 Department 实体，支持分页
func (controller *DepartmentController) FindAll(ctx *gin.Context) {
	log.Println("findAll department")

	// 从请求上下文中解析分页参数
	pg := utils.NewPagination(ctx)

	// 调用服务层方法获取分页数据
	departmentResponse := controller.departmentService.FindAll(pg)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    departmentResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
