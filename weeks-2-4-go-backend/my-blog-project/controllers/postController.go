package controllers

import (
	"my-blog-project/database"
	"my-blog-project/model"
	"my-blog-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct{}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=1000"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=1000"`
	Content string `json:"content" binding:"required,min=1"`
}

func (pc *PostController) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权的用户")
		return
	}

	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}
	if err := database.DB.Create(&post).Error; err != nil {
		utils.InternalServerError(c, "创建文章失败")
		return
	}

	database.DB.Preload("User").First(&post, post.ID)

	utils.Success(c, post)
}

func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []model.Post
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	if err := database.DB.Preload("User").Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts).Error; err != nil {
		utils.InternalServerError(c, "内部错误：获取文章信息失败")
		return
	}

	// 获取总数
	var total int64
	database.DB.Model(&model.Post{}).Count(&total)

	utils.Success(c, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})

}

func (pc *PostController) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "解析postID失败")
		return
	}

	var post model.Post
	if err := database.DB.Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		utils.NotFound(c, "未找到文章")
	}

	utils.Success(c, post)
}

func (pc *PostController) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "解析postID失败")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权的用户")
		return
	}

	var post model.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "未找到文章")
		return
	}

	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "只有本人可以更新自己的文章")
		return
	}

	// 更新文章
	post.Title = req.Title
	post.Content = req.Content

	if err := database.DB.Save(&post).Error; err != nil {
		utils.InternalServerError(c, "更新文章失败")
		return
	}

	// 预加载用户信息
	database.DB.Preload("User").First(&post, post.ID)

	utils.Success(c, post)
}

func (pc *PostController) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "解析postID失败")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权的用户")
		return
	}

	var post model.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "未找到文章")
		return
	}

	if post.UserID != userID.(uint) {
		utils.Forbidden(c, "只有本人可以删除自己的文章")
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		utils.InternalServerError(c, "删除文章失败")
		return
	}

	utils.Success(c, gin.H{
		"message": "文章删除成功",
	})
}
