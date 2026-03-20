package controllers

import (
	"my-blog-project/database"
	"my-blog-project/model"
	"my-blog-project/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// 登录
func (uc *UserController) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	var user model.User
	// 先查询用户
	if err := database.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.BadRequest(c, "用户不存在")
		return
	}

	// 验证密码
	pass := user.CheckPass(req.Password)
	if !pass {
		utils.BadRequest(c, "密码错误")
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, req.Username)
	if err != nil {
		utils.BadRequest(c, "token生成失败")
		return
	}

	utils.Success(c, AuthResponse{
		Token: token,
		User:  user,
	})
}

// 注册
func (uc *UserController) Register(c *gin.Context) {
	// 先解析请求参数
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 真正的业务逻辑：先检查用户名是否已存在
	var existingUser model.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "用户名已存在")
		return
	}

	// 检查邮箱是否已存在

	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.BadRequest(c, "邮箱已存在")
		return
	}

	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	// 创建新用户
	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "创建用户失败")
	}
	// 返回token等信息
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalServerError(c, "生成token失败")
		return
	}

	utils.Success(c, AuthResponse{
		Token: token,
		User:  user,
	})

}

// 查看个人信息
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "用户无权限")
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		utils.NotFound(c, "未查询到该用户")
		return
	}

	utils.Success(c, user)
}
