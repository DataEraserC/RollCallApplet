// @Title       user.go
// @Description 放置操作用户数据库的网站入口函数以及工具函数
// @Author      DataEraserC
// @Update      DataEraserC  (2024/2/17   21:54)

package main

import (
	"errors"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// MemberOf 用户加入的组织gorm数据库对象.用于记录用户进入了什么组织
type MemberOf struct {
	GroupID     uint `gorm:"unique"`
	Permissions string
}

// 每次要对用户数据库修改时必须先动态加载数据库

// @title         InitUser
// @description   初始化用户数据的文件夹以及用户数据库
// @auth          DataEraserC                                   (2024/2/17   21:54)
// @param         GlobalPath                            string              "指定数据存放在什么地方"
// @param         UserID                                uint                "指定用户的ID"
// @param         SafeMode                              bool                "是否自动迁移数据库模型"
// @return        UserDatabase                          *gorm.DB            "用户数据库"
// @return        err                                   error               "可能存在的错误"
func InitUser(GlobalPath string, UserID uint, SafeMode bool) (*gorm.DB, error) {
	UserDatabase, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/user/%d/database.db", GlobalPath, UserID)), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	if SafeMode {
		err = UserDatabase.AutoMigrate(&MemberOf{})
		if err != nil {
			return nil, errors.New("failed to AutoMigrate database")
		}
	}
	return UserDatabase, nil
}
