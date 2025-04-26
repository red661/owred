package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hoyang/ownsa/model"
	"hoyang/ownsa/utils"
)

// DbInstance 结构体包含了四个数据库连接实例
type DbInstance struct {
	DbConfig       *gorm.DB // 主数据库连接实例
	DbCredential   *gorm.DB // 用户凭证数据库连接实例
	DbOtherGroup   *gorm.DB // 其他数据库连接实例
	DbEventMessage *gorm.DB // 事件消息数据库连接实例
}

// DB 全局数据库实例
var DB *DbInstance

// Migrate 用于自动迁移所有数据库表结构
func Migrate() {
	// 在主数据库（DbConfig）中自动迁移表
	DB.DbConfig.AutoMigrate(&model.ControllerUser{})
	DB.DbConfig.AutoMigrate(&model.ControllerProp{})
	DB.DbConfig.AutoMigrate(&model.InterfaceBoard{})
	DB.DbConfig.AutoMigrate(&model.CardReaderProp{})
	DB.DbConfig.AutoMigrate(&model.InputProp{})
	DB.DbConfig.AutoMigrate(&model.OutputProp{})
	DB.DbConfig.AutoMigrate(&model.SystemVariableParam{})

	// 在用户凭证数据库（DbCredential）中自动迁移表
	DB.DbCredential.AutoMigrate(&model.People{})
	DB.DbCredential.AutoMigrate(&model.Credential{})
	DB.DbCredential.AutoMigrate(&model.CredentialAccess{})
	DB.DbCredential.AutoMigrate(&model.VeinData{})

	// 在其他数据库（DbOtherGroup）中自动迁移表
	DB.DbOtherGroup.AutoMigrate(&model.AccessGroup{})
	DB.DbOtherGroup.AutoMigrate(&model.DoorGroup{})
	DB.DbOtherGroup.AutoMigrate(&model.SchedGroup{})
	DB.DbOtherGroup.AutoMigrate(&model.APBController{})
	DB.DbOtherGroup.AutoMigrate(&model.ValidWiegandRule{})

	// 在事件消息数据库（DbEventMessage）中自动迁移表
	DB.DbEventMessage.AutoMigrate(&model.EventMessageData{})
}

// CloseDbConnection 关闭所有数据库连接
func CloseDbConnection() {
	// 关闭主数据库连接
	sqlDB, err := DB.DbConfig.DB()
	utils.ErrorPanic(err)
	defer sqlDB.Close()

	// 关闭用户凭证数据库连接
	sqlDB, err = DB.DbCredential.DB()
	utils.ErrorPanic(err)
	defer sqlDB.Close()

	// 关闭其他数据库连接
	sqlDB, err = DB.DbOtherGroup.DB()
	utils.ErrorPanic(err)
	defer sqlDB.Close()

	// 关闭事件消息数据库连接
	sqlDB, err = DB.DbEventMessage.DB()
	utils.ErrorPanic(err)
	defer sqlDB.Close()
}

// SetupDatabase 初始化并设置数据库连接
func SetupDatabase(confEnv *map[string]string, out *log.Logger) {
	// 设置日志级别，根据环境变量 "GIN_MODE" 决定
	logLevel := logger.Info
	if (*confEnv)["GIN_MODE"] == "release" {
		logLevel = logger.Error
	}
	logLevel = logger.Info //调试日志开启

	// 配置日志
	newLogger := logger.New(
		out, // gin-gonic logger
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值
			LogLevel:                  logLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			ParameterizedQueries:      false,       // 不记录 SQL 参数
			Colorful:                  true,        // 启用彩色日志
		},
	)

	// 设置 GORM 配置，包含日志设置
	conf := gorm.Config{
		Logger: newLogger,
	}

	// 初始化 DB 实例
	DB = &DbInstance{}

	// 数据库连接字符串格式
	connStrFmt := "%s?charset=utf8mb4&parseTime=True&loc=Local" // 连接字符串格式，支持UTF8编码和本地时区

	// 连接主数据库
	DbConfigSqlite3Url := (*confEnv)["DbConfigPath"]
	dsn := fmt.Sprintf(connStrFmt, DbConfigSqlite3Url)

	dbConfig, err := gorm.Open(sqlite.Open(dsn), &conf)
	if err != nil {
		panic("failed to connect database")
	}
	// 执行SQLite自动清理和压缩
	dbConfig.Exec("PRAGMA auto_vacuum=1;")
	dbConfig.Exec("vacuum;")

	DB.DbConfig = dbConfig

	// 连接用户凭证数据库
	DbCredentialSqlite3Url := (*confEnv)["DbCredentialPath"]
	dsn = fmt.Sprintf(connStrFmt, DbCredentialSqlite3Url)

	dbCredential, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 执行SQLite自动清理和压缩
	dbCredential.Exec("PRAGMA auto_vacuum=1;")
	dbCredential.Exec("vacuum;")

	DB.DbCredential = dbCredential

	// 连接其他数据库
	DbOtherGroupSqlite3Url := (*confEnv)["DbOtherGroupPath"]
	dsn = fmt.Sprintf(connStrFmt, DbOtherGroupSqlite3Url)

	dbOtherGroup, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 执行SQLite自动清理和压缩
	dbOtherGroup.Exec("PRAGMA auto_vacuum=1;")
	dbOtherGroup.Exec("vacuum;")

	DB.DbOtherGroup = dbOtherGroup

	// 连接事件消息数据库
	DbEventMessageSqlite3Url := (*confEnv)["DbEventMessagePath"]
	dsn = fmt.Sprintf(connStrFmt, DbEventMessageSqlite3Url)

	dbEventMessage, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 执行SQLite自动清理和压缩
	dbEventMessage.Exec("PRAGMA auto_vacuum=1;")
	dbEventMessage.Exec("vacuum;")

	DB.DbEventMessage = dbEventMessage

	// 输出日志，表示数据库已设置完成
	out.Println("SetupDatabase")
}
