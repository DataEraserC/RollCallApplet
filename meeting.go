// @Title       meeting.go
// @Description 放置操作会议数据库的网站入口函数以及工具函数
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

// MettingParticipants 记录用户是否参与会议的信息的gorm数据库对象
type MettingParticipants struct {
	UserID            uint
	ParticipationTime time.Time
}

// Sign 签到gorm数据库对象,记录签到的开始及结束时间
type Sign struct {
	ID      uint
	BeginAt time.Time
	EndAt   time.Time
}

// SignatureBook 用户签到数据库对象,记录用户是否签到
type SignatureBook struct {
	UserID uint
	SignID uint
}

// 每次要对会议数据库修改时必须先动态加载数据库

// @title         InitMeeting
// @description   初始化会议数据的文件夹以及会议数据库
// @auth          DataEraserC                                   (2024/2/17   21:54)
// @param         GlobalPath                            string              "指定数据存放在什么地方"
// @param         GroupID                               uint                "指定会议所属组织的ID"
// @param         MeetingID                             uint                "指定会议的ID"
// @param         SafeMode                              bool                "是否自动迁移数据库模型"
// @return        GroupDatabase                         *gorm.DB            "组织数据库"
// @return        err                                   error               "可能存在的错误"
func InitMeeting(GlobalPath string, GroupID uint, MeetingID uint, SafeMode bool) (*gorm.DB, error) {
	MeetingDataPath := fmt.Sprintf("%s/group/%d/meeting/%d", GlobalPath, GroupID, MeetingID)
	// 初始化MeetingDataPath文件夹
	if _, err := os.Stat(MeetingDataPath); os.IsNotExist(err) {
		// MeetingDataPath不存在，创建MeetingDataPath
		err := os.MkdirAll(MeetingDataPath, 0755)
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

	MeetingDatabase, err := gorm.Open(sqlite.Open(fmt.Sprintf(MeetingDataPath+"/database.db", GlobalPath, GroupID, MeetingID)), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	if SafeMode {
		err = MeetingDatabase.AutoMigrate(&MettingParticipants{}, &Sign{}, &SignatureBook{})
		if err != nil {
			return nil, errors.New("failed to AutoMigrate database")
		}
	}
	return MeetingDatabase, nil
}
