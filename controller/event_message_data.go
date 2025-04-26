package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"

	"hoyang/ownsa/data/response"
	"hoyang/ownsa/service"
	"hoyang/ownsa/utils"
)

// EventMessageDataController 控制器，用于处理事件消息数据的相关请求
type EventMessageDataController struct {
	eventMessageDataService service.EventMessageDataService // 依赖的服务层
}

// WebSocket Upgrader，用于将 HTTP 请求升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // 读缓冲区大小
	WriteBufferSize: 1024, // 写缓冲区大小
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源的请求升级 WebSocket
		return true
	},
}

// NewEventMessageDataController 创建一个新的 EventMessageDataController 实例
func NewEventMessageDataController(service service.EventMessageDataService) *EventMessageDataController {
	return &EventMessageDataController{
		eventMessageDataService: service,
	}
}

// FindById 根据事件消息数据的 ID 查询具体信息
func (controller *EventMessageDataController) FindById(ctx *gin.Context) {
	log.Println("findby eventMessageDataId")

	// 从 URL 参数中获取 eventMessageDataId 并转换为 uint 类型
	eventMessageDataId := cast.ToUint(ctx.Param("eventMessageDataId"))

	// 调用服务层方法获取数据
	eventMessageResponse := controller.eventMessageDataService.FindById(eventMessageDataId)

	// 构造返回的响应体
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    eventMessageResponse,
	}

	// 返回 JSON 格式的响应
	ctx.JSON(http.StatusOK, webResponse)
}

// FindAll 查询所有事件消息数据，并支持分页
func (controller *EventMessageDataController) FindAll(ctx *gin.Context) {
	log.Println("findAll eventMessageData")

	// 从请求中解析分页参数
	pg := utils.NewPagination(ctx)

	// 调用服务层方法获取分页数据
	eventMessageResponse := controller.eventMessageDataService.FindAll(pg)

	// 构造返回的响应体
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    eventMessageResponse,
	}

	// 返回 JSON 格式的响应
	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *EventMessageDataController) FindByTimeRange(ctx *gin.Context) {
	log.Println("FindByTimeRange eventMessageData")

	// 从请求中获取开始时间和结束时间
	startTime := ctx.Query("startTime")
	endTime := ctx.Query("endTime")

	// 调用服务层方法获取时间段内的数据
	eventMessageResponse := controller.eventMessageDataService.FindByTimeRange(startTime, endTime)

	// 构造返回的响应体
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    eventMessageResponse,
	}

	// 返回 JSON 格式的响应
	ctx.JSON(http.StatusOK, webResponse)
}

// WebSocketServer WebSocket 服务，用于实时推送事件消息数据
func (controller *EventMessageDataController) WebSocketServer(ctx *gin.Context) {
	// 将 HTTP 请求升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("%s, error while Upgrading websocket connection\n", err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 不断读取和处理 WebSocket 消息
	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("%s, error while reading message\n", err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}

		// 将接收到的消息（事件 ID）转换为 uint
		msgId := cast.ToUint(string(p))

		// 调用服务层方法查询事件 ID 之后的数据
		eventMessageResponse := controller.eventMessageDataService.FindAfterId(msgId)

		// 将响应数据转换为 JSON 格式
		data, err := json.Marshal(eventMessageResponse)
		if err != nil {
			log.Printf("%s, error while marshalling response data\n", err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}

		// 将响应数据写回客户端
		err = conn.WriteMessage(messageType, data)
		if err != nil {
			log.Printf("%s, error while writing message\n", err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}
	}
}

// Sync 同步事件消息数据
func (controller *EventMessageDataController) Sync(ctx *gin.Context) {
	log.Println("sync eventMessageData")

	// 调用服务层方法执行同步操作
	controller.eventMessageDataService.Sync()

	// 构造返回的响应体（同步完成后无需返回具体数据）
	webResponse := response.Response{
		Code:    http.StatusOK,
		Success: true,
		Data:    nil,
	}

	// 返回 JSON 格式的响应
	ctx.JSON(http.StatusOK, webResponse)
}
