package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email    string `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password string `json:"-" gorm:"not null"` // 加密存储

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Posts     []Post         `json:"posts,omitempty" gorm:"foreignKey:UserID"`    // 关联文章结构体，一个用户可以拥有多个文章
	Comments  []Comment      `json:"comments,omitempty" gorm:"foreignKey:UserID"` // 关联评论结构体，一个用户可以发送多个评论
}

// 对user进行加密
func (u *User) EncryptionPass() error {
	encPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encPass)
	return nil
}

// 验证密码是否正确

func (u *User) CheckPass(password string) bool {
	logrus.Info("CheckPass called - username: ", u.Username)
	logrus.Info("Input password length: ", password)
	logrus.Info("Stored hash: ", u.Password)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	logrus.Info("Compare result: ", err)
	return err == nil
}

// 让sql执行前进行密码的加密操作
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return u.EncryptionPass()
}
