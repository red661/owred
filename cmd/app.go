package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"

	"hoyang/ownsa/database"
	"hoyang/ownsa/repository"
	"hoyang/ownsa/router"
	"hoyang/ownsa/service"
	"hoyang/ownsa/utils"
)

// 运行备份脚本
// func runShellScript() {
// 	// 获取当前 Go 程序的执行目录
// 	execDir, err := os.Getwd()
// 	if err != nil {
// 		log.Printf("无法获取当前执行目录: %s\n", err)
// 		return
// 	}

// 	// 备份脚本的路径
// 	backupScript := filepath.Join(execDir, "backup.sh") // ./backup.sh

// 	// 检查 backup.sh 是否存在
// 	if _, err := os.Stat(backupScript); os.IsNotExist(err) {
// 		log.Printf("❌ 备份脚本不存在: %s\n", backupScript)
// 		return
// 	}

//		// 执行备份脚本
//		cmd := exec.Command("/bin/bash", backupScript)
//		cmd.Stdout = os.Stdout
//		cmd.Stderr = os.Stderr
//		err = cmd.Run()
//		if err != nil {
//			log.Printf("❌ 备份脚本执行失败: %s\n", err)
//		} else {
//			log.Println("✅ 自动备份任务已启动")
//		}
//	}
func main() {
	// 加载环境配置
	confEnv := utils.GetEnvConf()

	// 获取当前进程的锁文件，避免多进程运行时产生冲突
	utils.AcquireProcessIDLock(confEnv["PIDFile"])

	// 初始化日志记录器
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)

	// 启动自动备份脚本
	setupCrontab()

	// 针对非 ARM 架构，并且存在命令行参数时的操作处理
	if runtime.GOARCH != "arm" && len(os.Args) > 1 {
		// 定义命令行标志，用于特定功能的执行
		genPtr := flag.Bool("gen", false, "generate model.gen")                 // 生成模型文件
		cleanPtr := flag.Bool("clean", false, "remove all .db files")           // 删除数据库文件
		initPtr := flag.Bool("init", false, "migrate and initialize database")  // 初始化数据库
		syncMsgPtr := flag.Bool("sync_msg", false, "synchronize event message") // 同步事件消息

		// 解析命令行标志
		flag.Parse()

		// 根据标志执行对应操作
		if *genPtr {
			GenerateModel() // 调用模型生成函数
		} else if *cleanPtr {
			// 删除配置中指定的 SQLite 数据库文件
			DbConfigSqlite3Url := confEnv["DbConfigPath"]
			log.Printf("remove %s\n", DbConfigSqlite3Url)
			os.Remove(DbConfigSqlite3Url)
			DbCredentialSqlite3Url := confEnv["DbCredentialPath"]
			log.Printf("remove %s\n", DbCredentialSqlite3Url)
			os.Remove(DbCredentialSqlite3Url)
			DbOtherGroupSqlite3Url := confEnv["DbOtherGroupPath"]
			log.Printf("remove %s\n", DbOtherGroupSqlite3Url)
			os.Remove(DbOtherGroupSqlite3Url)
			DbEventMessageSqlite3Url := confEnv["DbEventMessagePath"]
			log.Printf("remove %s\n", DbEventMessageSqlite3Url)
			os.Remove(DbEventMessageSqlite3Url)
		} else if *initPtr {
			// 初始化数据库及其表结构
			database.SetupDatabase(&confEnv, logger)
			database.Migrate()
			service.InitializeDatabase()
		} else if *syncMsgPtr {
			// 同步事件消息到数据库
			database.SetupDatabase(&confEnv, logger)
			validate := validator.New()
			eventMessageDataRepository := repository.NewEventMessageDataRepositoryImpl(database.DB.DbEventMessage)
			eventMessageDataService := service.NewEventMessageDataServiceImpl(eventMessageDataRepository, validate)
			eventMessageDataService.Sync()
		}

		// 命令行任务完成后退出
		return
	}

	// 检查是否存在 "SystemResetFactoryFile"，如果不存在则继续检查数据库初始化
	if _, err := os.Stat(confEnv["SystemResetFactoryFile"]); errors.Is(err, os.ErrNotExist) {
		if _, err := os.Stat(confEnv["DataDir"]); errors.Is(err, os.ErrNotExist) {
			// 如果数据目录不存在，则重新初始化数据库
			database.SetupDatabase(&confEnv, logger)
			database.Migrate()
			service.InitializeDatabase()
		}

		// 启动 HTTP 服务器并阻塞主进程，等待退出信号
		done := make(chan os.Signal)
		go func(confEnv *map[string]string, out *log.Logger) {
			// 创建并运行 HTTP 服务
			if err := router.CreateHttpServ(confEnv, logger); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("error: %s\n", err)
			}
		}(&confEnv, logger)

		// 阻塞直到接收到退出信号
		<-done
	} else {
		// 如果存在系统重置文件，删除所有数据库文件并重新创建
		DbConfigSqlite3Url := confEnv["DbConfigPath"]
		log.Printf("remove %s\n", DbConfigSqlite3Url)
		os.Remove(DbConfigSqlite3Url)
		DbCredentialSqlite3Url := confEnv["DbCredentialPath"]
		log.Printf("remove %s\n", DbCredentialSqlite3Url)
		os.Remove(DbCredentialSqlite3Url)
		DbOtherGroupSqlite3Url := confEnv["DbOtherGroupPath"]
		log.Printf("remove %s\n", DbOtherGroupSqlite3Url)
		os.Remove(DbOtherGroupSqlite3Url)
		DbEventMessageSqlite3Url := confEnv["DbEventMessagePath"]
		log.Printf("remove %s\n", DbEventMessageSqlite3Url)
		os.Remove(DbEventMessageSqlite3Url)

		// 删除数据目录
		dataDir := confEnv["DataDir"]
		log.Printf("remove %s\n", dataDir)
		os.Remove(dataDir)

		// 删除系统重置文件
		os.Remove(confEnv["SystemResetFactoryFile"])

		// 同步文件系统并重启系统
		syscall.Sync()
		syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	}
}

// setupCrontab 设置定时备份任务
// setupCrontab 设置定时备份任务
func setupCrontab() {
	// 获取当前执行路径
	execDir, err := os.Getwd()
	if err != nil {
		log.Printf("❌ 无法获取当前执行路径: %v", err)
		return
	}

	// 动态拼接备份脚本和日志文件路径
	backupScript := filepath.Join(execDir, "backup.sh")
	logFile := filepath.Join(execDir, "db_backup.log")

	log.Printf("🔹 备份脚本路径: %s", backupScript)
	log.Printf("🔹 日志文件路径: %s", logFile)

	// 检查备份脚本是否存在
	if _, err := os.Stat(backupScript); os.IsNotExist(err) {
		log.Printf("❌ 备份脚本不存在: %s", backupScript)
		return
	}

	// 创建 cron 实例 (支持秒级精度)
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)),
		),
	)

	// 添加定时任务 (每月1号凌晨0点0分执行)
	// cron 表达式格式: 秒 分 时 日 月 周
	_, err = c.AddFunc("0 0 0 1 * *", func() {
		log.Println("💾 开始执行数据库备份...")

		cmd := exec.Command("/bin/bash", backupScript)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		startTime := time.Now()
		if err := cmd.Run(); err != nil {
			log.Printf("❌ 备份失败: %v (耗时: %v)", err, time.Since(startTime))
		} else {
			log.Printf("✅ 备份成功 (耗时: %v)", time.Since(startTime))
		}
	})

	if err != nil {
		log.Printf("❌ 添加定时任务失败: %v", err)
		return
	}

	// 启动 cron 服务
	c.Start()
	log.Println("✅ 定时备份任务已启动")

	// 可选: 添加优雅停止逻辑
	// 可以在 main 函数中捕获退出信号时调用 c.Stop()
}
