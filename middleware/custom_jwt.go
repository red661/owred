package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"hoyang/ownsa/data/response"
	"hoyang/ownsa/utils"
)

// CreateToken 用于创建 JWT token，包含用户的 ID、用户名和过期时间
func CreateToken(id uint, username string, tokentime uint) (string, error) {
	// 读取环境变量配置
	confEnv, err := godotenv.Read() // 读取项目根目录下的 .env 文件
	utils.ErrorPanic(err)

	// 获取 SecretKey 用于签名
	secretKey := []byte(confEnv["SecretKey"])
	tokentimeOut := time.Now().Add(time.Minute * time.Duration(tokentime)).Unix()

	// 创建 JWT token 并设置 Claims（即 payload）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       fmt.Sprintf("%d", id), // 用户 ID
		"username": username,              // 用户名
		"exp":      tokentimeOut,          // 过期时间为 24 小时后
	})

	// 签名生成 token 字符串
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err // 返回错误
	}

	return tokenString, nil // 返回生成的 token 字符串
}

// VerifyToken 用于验证 JWT token 的合法性，并返回其中的用户 ID
func VerifyToken(secretKey []byte, tokenString string) (id uint, err error) {
	// 解析 token，获取其中的 Claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil // 使用密钥验证 token
	})

	if err != nil {
		return 0, err // 如果解析失败，返回错误
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token") // 如果 token 不合法，返回错误
	}

	// 获取 token 中的 Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("wrong JWT Claims") // 如果 Claims 格式不正确，返回错误
	}

	// 解析 id 字段，并转换为 uint 类型
	uid, err := strconv.ParseUint(claims["id"].(string), 10, 32)
	// 打开 SQLite 数据库
	db, err := sql.Open("sqlite3", "./appdata/db/config.db")
	// db, err := sql.Open("sqlite3", "./ownsadb/config.db") //部署
	if err != nil {
		return 0, fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// 查询数据库中的 token
	var dbToken string
	query := "SELECT token FROM ownsa_controller_user WHERE id = ?"
	err = db.QueryRow(query, uid).Scan(&dbToken)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("user not found in database")
	} else if err != nil {
		return 0, fmt.Errorf("database error: %v", err)
	}

	// 对比 token
	if dbToken != tokenString {
		return 0, fmt.Errorf("token expired or logged in from another device")
	}
	return uint(uid), err // 返回解析后的用户 ID
}

// TokenAuthMiddleware 是一个 Gin 中间件，用于验证请求中的 JWT token
func TokenAuthMiddleware(confEnv *map[string]string) gin.HandlerFunc {
	// 获取 SecretKey，用于验证 token
	secretKey := []byte((*confEnv)["SecretKey"])

	return func(ctx *gin.Context) {
		// 从请求头中获取 Authorization token
		token := ctx.Request.Header.Get("Authorization")
		// 如果 Authorization 头部没有 token，尝试从 Cookie 中获取
		if token == "" {
			if cookie, err := ctx.Request.Cookie("X-Authorization"); err == nil {
				token = cookie.Value
			}
		}

		// 验证 token 是否有效，并获取用户 ID
		id, err := VerifyToken(secretKey, token)

		// 如果验证失败，返回 Unauthorized 错误
		if err != nil {
			webResponse := response.Response{}
			webResponse.Code = http.StatusUnauthorized
			webResponse.Success = false
			webResponse.Message = "Unauthorized"
			ctx.AbortWithStatusJSON(http.StatusOK, webResponse)
		} else {
			// 如果验证通过，将用户 ID 存储在上下文中，并继续处理请求
			ctx.Set("id", id)
			ctx.Next()
		}
	}
}
