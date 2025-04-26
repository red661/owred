package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"hoyang/ownsa/data/request"
	"hoyang/ownsa/database"
	"hoyang/ownsa/middleware"
	"hoyang/ownsa/model"
	"hoyang/ownsa/service"

	"hoyang/ownsa/router"
	"hoyang/ownsa/utils"

	"hoyang/ownsa/repository"
)

func TestPingRoute(t *testing.T) {
	server := gin.New()
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	server.Use(gin.LoggerWithWriter(logger.Writer()))
	server.Use(gin.CustomRecovery(middleware.ErrorHandler))

	router := router.SetupRouter(server)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body.String())
}

func TestPwdHash(t *testing.T) {
	pwdHash, err := utils.HashPassword("3CoreP@w")
	utils.ErrorPanic(err)
	assert.Equal(t, pwdHash, "")
}

// 创建部门
func TestSaveDepartmentWithCredentials(t *testing.T) {

	confEnv := utils.GetEnvConf()
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	database.SetupDatabase(&confEnv, logger)

	departRepository := repository.NewDepartmentRepositoryImpl(database.DB.DbCredential)
	// peopleRepository := repository.NewPeopleRepositoryImpl(database.DB.DbCredential)

	// 创建部门
	depart := model.Department{Name: "H12"}
	newDepart, err := departRepository.Save(depart)
	if err != nil {
		fmt.Printf("Error saving department: %v\n", err)
		return
	}
	// 打印返回值中的信息
	fmt.Printf("Department Name: %s\n", newDepart.Name)
	fmt.Printf("Department ID: %d\n", newDepart.ID)

}

// 创建员工并关联部门
func TestNewPeopleDepartmentWithRelevance(t *testing.T) {
	confEnv := utils.GetEnvConf()
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	database.SetupDatabase(&confEnv, logger)

	// departRepository := repository.NewDepartmentRepositoryImpl(database.DB.DbCredential)
	peopleRepository := repository.NewPeopleRepositoryImpl(database.DB.DbCredential)
	// 创建员工并关联部门
	departmentId := uint(1)
	person := model.People{FirstName: "tmi", DepartmentID: departmentId}
	newPeopel, err := peopleRepository.Save(person)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
		return
	}
	fmt.Printf("Person's DepartmentID: %d\n", newPeopel.DepartmentID)
	fmt.Printf("Person's NAME: %s\n", newPeopel.FirstName)

}
func TestFindDepartmentWithCredentials(t *testing.T) {
	confEnv := utils.GetEnvConf()
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	database.SetupDatabase(&confEnv, logger)

	departRepository := repository.NewDepartmentRepositoryImpl(database.DB.DbCredential)
	peopleRepository := repository.NewPeopleRepositoryImpl(database.DB.DbCredential)

	departmentId := uint(1) // 假设要查询的部门 ID 为 1
	// 调用 FindById 方法
	department, _, err := departRepository.FindById(departmentId)
	if err != nil {
		t.Errorf("Failed to find department with ID %d: %v", departmentId, err)
		return
	}
	peopledepartment, _, err := peopleRepository.FindById(departmentId)
	if err != nil {
		t.Errorf("Failed to find department with ID %d: %v", departmentId, err)
		return
	}
	// 打印部门信息
	fmt.Printf("Department ID: %d\n", department.ID)
	fmt.Printf("Department Name: %s\n", department.Name)
	fmt.Printf("Department Name: %s\n", peopledepartment.FirstName)

}

// 删除部门
func TestUpdateUserControlData(t *testing.T) {
	confEnv := utils.GetEnvConf()
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	database.SetupDatabase(&confEnv, logger)

	userRepository := repository.NewControllerUserRepositoryImpl(database.DB.DbConfig)

	userId := uint(4)

	// 验证删除是否生效
	cu, err := userRepository.FindById(userId)
	if err != nil {
		t.Errorf("can't find userid=%v\n", userId)
		return
	}
	cu.Permission1 = 0
	userRepository.Update(*cu)
}

// 删除部门
func TestDeleteDepartmentWithRelevance(t *testing.T) {
	confEnv := utils.GetEnvConf()
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	database.SetupDatabase(&confEnv, logger)

	departRepository := repository.NewDepartmentRepositoryImpl(database.DB.DbCredential)

	departmentId := uint(1)
	departRepository.Delete(departmentId)

	// 验证删除是否生效
	_, _, err := departRepository.FindById(departmentId)
	if err == nil {
		t.Errorf("Department with ID %d was not deleted", departmentId)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Unexpected error when querying deleted department: %v", err)
	} else {
		fmt.Printf("Department with ID %d successfully deleted\n", departmentId)
	}
}

// 定义一个 MockDepartmentRepository 来模拟 DepartmentRepository 的行为
type MockDepartmentRepository struct {
	mock.Mock
}

func (m *MockDepartmentRepository) Save(dept model.Department) (*model.Department, error) {
	args := m.Called(dept)
	// 返回 *model.Department 和 error
	return args.Get(0).(*model.Department), args.Error(1)
}

func (m *MockDepartmentRepository) Delete(id uint) {
	// Delete 无返回值
	m.Called(id)
}

func (m *MockDepartmentRepository) FindById(id uint) (*model.Department, []*model.Credential, error) {
	// 对于此测试我们不关心此方法的实际调用，可根据需要实现
	args := m.Called(id)
	var dept *model.Department
	if args.Get(0) != nil {
		dept = args.Get(0).(*model.Department)
	}
	var creds []*model.Credential
	if args.Get(1) != nil {
		creds = args.Get(1).([]*model.Credential)
	}
	return dept, creds, args.Error(2)
}

func (m *MockDepartmentRepository) FindAll(pg *utils.Pagination) []*model.Department {
	// 对于此测试我们不关心此方法的实际调用，可根据需要实现
	args := m.Called(pg)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]*model.Department)
}

// control 调用 service
func TestDepartmentService_Create(t *testing.T) {
	// 创建 mock 仓库实例
	mockDepartmentRepository := new(MockDepartmentRepository)

	// 准备返回的数据
	expectedDepartment := &model.Department{
		ID:   1,
		Name: "HR",
	}

	// 当调用 Save 方法时，返回期望的部门数据和 nil error
	mockDepartmentRepository.On("Save", mock.AnythingOfType("model.Department")).
		Return(expectedDepartment, nil)

	// 创建服务实例
	validate := validator.New()
	deptService := service.NewDepartmentServiceImpl(mockDepartmentRepository, validate)

	// 创建请求数据

	req := request.CreateDepartmentRequest{Name: "HR"}

	// 调用 Create 方法
	resp := deptService.Create(req)

	// 验证结果
	assert.Equal(t, "HR", resp.Name)
	assert.Equal(t, uint(1), resp.ID)

	// 验证期望
	mockDepartmentRepository.AssertExpectations(t)
}

// 处理 http 请求
func TestPeopleDepartControl(t *testing.T) {
	// 设置 gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化测试环境（根据需要初始化数据库、mock数据等）
	confEnv := utils.GetEnvConf()
	logger := utils.NewLogger()
	database.SetupDatabase(&confEnv, logger)

	// 创建 gin.Engine 实例并注册路由
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(logger.Writer()))
	engine.Use(gin.CustomRecovery(middleware.ErrorHandler))
	router.SetupRouter(engine)
	router.CreateWebController()
	router.RegisterPeopleRoutes(&confEnv, engine, router.WebController.PeopleController)
	// 模拟请求 /api/people GET 方法来测试 FindAll
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/people", nil)
	// 如果需要鉴权，可以在 header 中加入 Token，如：
	// req.Header.Set("Authorization", "Bearer your_test_jwt_token")

	engine.ServeHTTP(w, req)

	// 断言状态码为 200
	assert.Equal(t, http.StatusOK, w.Code)

	// 根据你的返回数据结构进行断言
	// 假设 /api/people 返回结构如下:
	// {
	//    "code": 200,
	//    "success": true,
	//    "data": [
	//        { "id": 1, "name": "John Doe" },
	//        ...
	//    ]
	// }
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// 验证返回的 code 和 success 字段
	assert.Equal(t, float64(200), resp["code"])
	assert.Equal(t, true, resp["success"])

	// 验证 data 是一个 slice
	data, ok := resp["data"].([]interface{})
	assert.True(t, ok, "data 字段应该是一个数组")
	// 如果有期望至少返回1条数据，可以进一步断言 data 的长度和内容
	// assert.NotEmpty(t, data)
	if len(data) > 0 {
		firstPerson := data[0].(map[string]interface{})
		assert.Equal(t, float64(1), firstPerson["id"])
		assert.Equal(t, "John Doe", firstPerson["name"])
	}
	// 至此，已验证接口基本可用性，你可以根据实际业务需求继续添加更多断言。
}

func setupDatabaseForTesting(db *gorm.DB) error {
	// 自动迁移 EventMessageData 表
	return db.AutoMigrate(&model.EventMessageData{})
}
func getTestDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}

// TestInsertEventMessageData 测试插入 100 条测试数据
// func TestInsertEventMessageData(t *testing.T) {
// 	// 使用 SQLite 内存数据库
// 	db, err := getTestDB()
// 	if err != nil {
// 		t.Fatalf("Failed to connect to database: %v", err)
// 	}

// 	// 初始化表结构
// 	err = setupDatabaseForTesting(db)
// 	if err != nil {
// 		t.Fatalf("Failed to setup database: %v", err)
// 	}

// 	// 初始化仓库
// 	repo := repository.NewEventMessageDataRepositoryImpl(db)

// 	// 准备测试数据
// 	rand.Seed(time.Now().UnixNano())
// 	var testData []*model.EventMessageData
// 	for i := 1; i <= 100; i++ {
// 		accessTime := fmt.Sprintf("2024-12-23 %02d:%02d:%02d", rand.Intn(24), rand.Intn(60), rand.Intn(60))
// 		testData = append(testData, &model.EventMessageData{
// 			MsgId:           uint(i),
// 			AccessTime:      accessTime,
// 			IBName:          fmt.Sprintf("Module-%d", rand.Intn(10)+1),
// 			ReaderName:      fmt.Sprintf("Door-%d", rand.Intn(5)+1),
// 			PeopleCode:      fmt.Sprintf("P%03d", rand.Intn(1000)),
// 			PeopleFirstName: fmt.Sprintf("First%d", rand.Intn(100)),
// 			PeopleLastName:  fmt.Sprintf("Last%d", rand.Intn(100)),
// 			PeopleDepart:    fmt.Sprintf("Depart%d", rand.Intn(5)+1),
// 			CardNo:          fmt.Sprintf("C%010d", rand.Intn(1000000000)),
// 			Wiegand:         uint(rand.Intn(1000)),
// 			UniqueId:        uint(rand.Intn(1000)),
// 			PeopleId:        uint(rand.Intn(1000)),
// 			Content:         uint(rand.Intn(1000)),
// 			EventType:       uint(rand.Intn(10)),
// 			IBAddr:          rand.Intn(256),
// 			ReaderAddr:      rand.Intn(256),
// 			InputAddr:       rand.Intn(256),
// 			OutputAddr:      rand.Intn(256),
// 			FullName:        fmt.Sprintf("First%d Last%d", rand.Intn(100), rand.Intn(100)),
// 		})
// 	}

// 	// 插入测试数据
// 	if err := repo.SaveAll(testData, 20); err != nil {
// 		t.Fatalf("SaveAll failed: %v", err)
// 	}

// 	// 验证插入数据
// 	lastMsgId := repo.LastMsgId()
// 	assert.GreaterOrEqual(t, lastMsgId, uint(100), "Last MsgId should be at least 100")

// 	// 清理测试数据
// 	db.Exec("DELETE FROM event_message_data WHERE accesstime LIKE '2024-12-23%'")
// 	t.Logf("Inserted and validated 100 rows of test data in event_message_data")
// }
