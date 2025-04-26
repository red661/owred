package model

import (
	"gorm.io/gorm"
)

const (
	UserTypeInternal = iota // 0
	UserTypeFactory  = iota // 1：出厂设置
	UserTypeManager  = iota // 2：经销服务商
	UserTypeOwnsa    = iota // 3：Ownsa用户
)

const (
	BPTypeDefault = iota
	BPTypeBIS     = 99
)

// 控制器用户
type ControllerUser struct {
	gorm.Model

	Username      string `gorm:"type:varchar(50);unique;not null"` // 用户名
	Password      string `gorm:"type:varchar(50);not null"`        // 密码
	Token         string `gorm:"not null"`                         // Token
	UserType      uint   `gorm:"not null"`                         // 类型 1：出厂设置 2：经销服务商 3：Ownsa用户 类型为3时 权限字段才生效
	Permission1   uint   `gorm:"not null"`                         // 权限1 系统设置
	Permission2   uint   `gorm:"not null"`                         // 权限2 设备管理
	Permission3   uint   `gorm:"not null"`                         // 权限3 设备维护
	Permission4   uint   `gorm:"not null"`                         // 权限4 人员管理
	Permission5   uint   `gorm:"not null"`                         // 权限5 统计分析
	Permission6   uint   `gorm:"not null"`                         // 权限6 备用
	Permission7   uint   `gorm:"not null"`                         // 权限7 备用
	Permission8   uint   `gorm:"not null"`                         // 权限8 备用
	LastLoginTime uint   // 最近登录时间 UNIX时间戳
}

// TableName 返回 ControllerUser 类型的表名。
// 该方法主要用于ORM（对象关系映射），告知框架数据库中对应的表。
// 没有参数。
// 返回值为字符串，表示数据库中的表名。
func (ControllerUser) TableName() string {
	return "red_controller_user"
}
