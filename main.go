// @Title       main.go
// @Description 放置主函数及调用gin等
// @Author      DataEraserC
// @Update      DataEraserC  (2024/2/17   21:54)

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var (
	DataPath = "data"
	LogPath  = "logs"
	GinPort  = ":8080"
)

var (
	GlobalDatabase *gorm.DB = nil
)

func main() {
	// 检测 LogPath是否存在 不存在需先创建
	// LogPath 目前只能先进行检测
	if _, err := os.Stat(LogPath); os.IsNotExist(err) {
		// LogPath不存在，创建LogPath
		err := os.MkdirAll(LogPath, 0755)
		if err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
		} else {
			fmt.Println("Log directory created successfully!")
		}
	} else if err != nil {
		fmt.Printf("Error checking log directory: %v\n", err)
	} else {
		fmt.Println("Log directory already exists!")
	}

	f, err := os.OpenFile(LogPath+"/log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()

	// 设置log输出为同时文件及输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SeDataEraserCags(log.Ldate | log.Ltime | log.Lshortfile)

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 用户登录接口(账号密码)
	r.POST("/login_account_password", Login_account_password(GlobalDatabase))

	// 用户登录接口(微信)
	r.POST("/login_wx", Login_wx(GlobalDatabase, WXAppID, WXAppSecret))

	// 用户获取个人信息接口
	r.POST("/userinfo", Userinfo(GlobalDatabase))

	// 用户修改个人信息接口
	r.POST("/updateuserinfo", Updateuserinfo(GlobalDatabase))

	// 用户注销登陆接口
	r.POST("/logout", Logout(GlobalDatabase))

	err = r.Run(GinPort)
	if err != nil {
		panic("failed at r.Run()")
	}
}

func init() {
	// 导入环境变量覆盖敏感密钥
	if envWXAppID := os.Getenv("WXAppID"); envWXAppID != "" {
		WXAppID = envWXAppID
	}

	if envWXAppSecret := os.Getenv("WXAppSecret"); envWXAppSecret != "" {
		WXAppSecret = envWXAppSecret
	}

	if envJWTSecretKey := os.Getenv("JWTSecretKey"); envJWTSecretKey != "" {
		JWTSecretKey = envJWTSecretKey
	}

	// 导入环境变量
	if envDataPath := os.Getenv("DataPath"); envDataPath != "" {
		DataPath = envDataPath
	}

	if envLogPath := os.Getenv("LogPath"); envLogPath != "" {
		LogPath = envLogPath
	}

	if envGinPort := os.Getenv("GinPort"); envGinPort != "" {
		GinPort = envGinPort
	}

	// 调用子模块函数初始化

	// 全局唯一的资源(必须加载)
	log.Println("Initializing global Resource......")
	var err error
	GlobalDatabase, err = InitGlobal(DataPath, true)
	if err != nil {
		panic("failed at init()")
	}

	log.Println("Initialize global Resource successfully......")

}
