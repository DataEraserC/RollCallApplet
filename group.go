// @Title       group.go
// @Description 放置操作组织数据库的网站入口函数以及工具函数
// @Author      DataEraserC
// @Update      DataEraserC  (2024/2/17   21:54)

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"go.starlark.net/lib/time"
	"gorm.io/gorm"
)

// MemberInfo 成员信息gorm对象,记录了成员权限
type MemberInfo struct {
	UserID      uint
	Permissions string
}

// MeetingInfo 会议信息gorm对象,记录了对应ID的会议的会议描述及开始结束时间
type MeetingInfo struct {
	ID                 uint
	BeginAt            time.Time
	EndAt              time.Time
	MeetingDescription string
}

// 每次要对组织数据库修改时必须先动态加载数据库

// @title         InitGroup
// @description   初始化组织数据的文件夹以及组织数据库
// @auth          DataEraserC                                   (2024/2/17   21:54)
// @param         GlobalPath                            string              "指定数据存放在什么地方"
// @param         GroupID                               uint                "指定读取的组织的ID"
// @param         SafeMode                              bool                "是否自动迁移数据库模型"
// @return        GroupDatabase                         *gorm.DB            "组织数据库"
// @return        err                                   error               "可能存在的错误"
func InitGroup(GlobalPath string, GroupID uint, SafeMode bool) (*gorm.DB, error) {
	GroupDataPath := fmt.Sprintf("%s/group/%d", GlobalPath, GroupID)
	// 初始化GroupDataPath文件夹
	if _, err := os.Stat(GroupDataPath); os.IsNotExist(err) {
		// GroupDataPath不存在，创建GroupDataPath
		err := os.MkdirAll(GroupDataPath, 0755)
		if err != nil {
			log.Printf("Failed to create global directory: %v\n", err)
		} else {
			log.Println("Global directory created successfully!")
		}
	} else if err != nil {
		log.Printf("Error checking Global directory: %v\n", err)
	} else {
		log.Println("Global directory already exists!")
	}

	GroupDatabase, err := gorm.Open(sqlite.Open(GroupDataPath+"/database.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	if SafeMode {
		err = GroupDatabase.AutoMigrate(&MemberInfo{}, &MeetingInfo{})
		if err != nil {
			return nil, errors.New("failed to AutoMigrate database")
		}
	}
	return GroupDatabase, nil
}
