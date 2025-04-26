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

// PeopleController 用于管理 People 实体的控制器
type PeopleController struct {
	peopleService service.PeopleService // 依赖的服务层，用于处理 People 数据的业务逻辑
}

// NewPeopleController 创建并返回一个新的 PeopleController 实例
func NewPeopleController(service service.PeopleService) *PeopleController {
	return &PeopleController{
		peopleService: service,
	}
}

func (controller *PeopleController) Import(ctx *gin.Context) {
	log.Println("import people")

	// 解析请求
	importRequest := request.ImportPeopleRequest{}
	if err := ctx.ShouldBindJSON(&importRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "无效的请求参数",
		})
		return
	}

	// 调用服务层处理导入
	result, err := controller.peopleService.Import(importRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 返回导入结果
	ctx.JSON(http.StatusOK, response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    result,
	})

	// 调用同步方法
	DataSync()
}

// Create 创建一个新的 People 实体
func (controller *PeopleController) Create(ctx *gin.Context) {
	log.Println("create people")

	// 解析并绑定请求体到 CreatePeopleRequest 结构体
	createPeopleRequest := request.CreatePeopleRequest{}
	err := ctx.ShouldBindJSON(&createPeopleRequest)
	utils.ErrorPanic(err) // 如果解析失败则抛出错误

	// 使用 litter 库格式化打印请求内容，方便调试
	log.Printf("%s", litter.Sdump(createPeopleRequest))

	// 调用服务层方法创建 People 实体
	controller.peopleService.Create(createPeopleRequest)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)

	// 调用同步方法
	DataSync()
}

// Update 更新指定的 People 实体
func (controller *PeopleController) Update(ctx *gin.Context) {
	log.Println("update people")

	// 解析并绑定请求体到 UpdatePeopleRequest 结构体
	updatePeopleRequest := request.UpdatePeopleRequest{}
	err := ctx.ShouldBindJSON(&updatePeopleRequest)
	utils.ErrorPanic(err) // 如果解析失败则抛出错误

	// 使用 litter 库格式化打印请求内容，方便调试
	log.Printf("%s", litter.Sdump(updatePeopleRequest))

	// 从 URL 参数中获取 peopleId 并设置到 UpdatePeopleRequest 中
	peopleId := cast.ToUint(ctx.Param("peopleId"))
	updatePeopleRequest.ID = peopleId

	// 调用服务层方法更新 People 实体
	controller.peopleService.Update(updatePeopleRequest)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)

	// 调用同步方法
	DataSync()
}

// Delete 删除指定的 People 实体
func (controller *PeopleController) Delete(ctx *gin.Context) {
	log.Println("delete people")

	// 从 URL 参数中获取 peopleId
	peopleId := cast.ToUint(ctx.Param("peopleId"))

	// 调用服务层方法删除指定的 People 实体
	controller.peopleService.Delete(peopleId)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)

	// 调用同步方法
	DataSync()
}

// FindById 根据 ID 查询指定的 People 实体
func (controller *PeopleController) FindById(ctx *gin.Context) {
	log.Println("findby peopleId")

	// 从 URL 参数中获取 peopleId
	peopleId := cast.ToUint(ctx.Param("peopleId"))

	// 调用服务层方法查询指定 ID 的 People 实体
	peopleResponse := controller.peopleService.FindById(peopleId)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    peopleResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// FindAll 查询所有的 People 实体，支持分页
func (controller *PeopleController) FindAll(ctx *gin.Context) {
	log.Println("findAll people")

	// 从请求上下文中解析分页参数
	pg := utils.NewPagination(ctx)

	// 调用服务层方法获取分页数据
	peopleResponse := controller.peopleService.FindAll(pg)

	// 构造响应并返回
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    peopleResponse,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
