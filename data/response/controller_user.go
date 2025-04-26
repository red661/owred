package response

// 用户管理接口中返回的用户信息
type ControllerUserResponse struct {
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	UserType      uint   `json:"usertype"`
	Permission1   uint   `json:"permission1"` // 权限1 系统设置
	Permission2   uint   `json:"permission2"` // 权限2 设备管理
	Permission3   uint   `json:"permission3"` // 权限3 设备维护
	Permission4   uint   `json:"permission4"` // 权限4 人员管理
	Permission5   uint   `json:"permission5"` // 权限5 统计分析
	Permission6   uint   `json:"permission6"` // 权限6 备用
	Permission7   uint   `json:"permission7"` // 权限7 备用
	Permission8   uint   `json:"permission8"` // 权限8 备用
	LastLoginTime uint   `json:"last_login_time"`
}

// 用户认证成功后返回的令牌信息
type TokenResponse struct {
	ID       uint   `json:"id"`
	UserType uint   `json:"user_type"`
	Token    string `json:"token"`
}
