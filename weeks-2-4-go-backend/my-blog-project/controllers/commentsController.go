package controllers

import (
	"my-blog-project/database"
	"my-blog-project/model"
	"my-blog-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct{}

type CreateCommentsReq struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	// 先获取入参的文章id，并解析
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "postID转换为uint失败：")
		return
	}

	// 然后转请求参数为json
	var req CreateCommentsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 从上下文获取用户id
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "用户未授权")
		return
	}

	// 检查文章是否存在
	var post model.Post

	if err := database.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "未找到改文章")
		return
	}

	// 生成评论准备插入
	comment := model.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  uint(postID),
	}

	// 插入评论
	if err := database.DB.Create(&comment).Error; err != nil {
		utils.InternalServerError(c, "创建评论失败")
		return
	}

	database.DB.Preload("User").Preload("Post").Preload("Post.User").First(&comment, comment.ID)

	utils.Success(c, comment)
}

// 获取文章的评论列表

func (cc *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章id")
		return
	}

	// 检查文章是否存在
	var post model.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		utils.NotFound(c, "未找到该文章")
		return
	}

	var comments []model.Comment

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 查询评论列表，预加载用户信息
	if err := database.DB.Preload("User").
		Preload("Post").
		Preload("Post.User").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&comments).Error; err != nil {
		utils.InternalServerError(c, "获取评论失败")
		return
	}

	// 获取总数
	var total int64
	database.DB.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&total)

	utils.Success(c, gin.H{
		"comments":  comments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
