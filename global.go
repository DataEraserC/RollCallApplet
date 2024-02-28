// @Title       global.go
// @Description 放置操作全局数据库的网站入口函数以及工具函数
// @Author      DataEraserC
// @Update      DataEraserC  (2024/2/17   21:54)

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserInfo 用户gorm对象，定义了数据库内用户的信息
type UserInfo struct {
	// ID字段就是UserID
	ID                 uint `gorm:"primaryKey;AUTO_INCREMENT"`
	Avatar             string
	Name               string
	NickName           string
	Gender             string
	Collage            string
	Majar              string
	Grade              uint
	PhoneNumber        string
	RegistrationNumber string
}

/*
所有用户的登陆信息必须存在同一个数据库 以便登陆时验证身份
*/

// Login 用户登陆信息gorm对象,定义了用户名密码等用于登陆的信息
type Login struct {
	UserID   uint   `gorm:"unique;not null"`
	Username string `gorm:"unique"`
	// Password字段可以在后端加盐并hash
	// 但后端必须保证盐会被长久保存
	Password string
	OpenID   string `gorm:"unique"`
	// 修改login表时撤销所有token即可 无需为此添加UpdateAt字段
	// UpdateAt int64
}

/*
如果要求"能只通过Token就查询到UserID"则所有用户的Token必须存放在同一个数据库
如果要每个用户一个Token数据库则需要在前端无法获取到UserID时要求重新登陆
*/

// Token 令牌gorm对象,定义了用户的Token,用于操作时(比如加入组织等)鉴权
type Token struct {
	UserID    uint
	Token     string `gorm:"unique"`
	CreatedAt int64
}

// WXLoginResp 微信登陆返回值json对象,用于接收微信登陆函数的返回值,(不重要)
type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// GroupInfo 部门信息gorm对象,记录了最基础的部门id和部门Code的对应关系
type GroupInfo struct {
	ID        uint
	GroupCode string `gorm:"unique"`
}

// CreateGroupRequest 创建部门请求gorm对象,记录了创建部门的申请
type CreateGroupRequest struct {
	ID               uint
	UserID           uint
	Reason           string
	GroupName        string
	GroupCode        string `gorm:"unique"`
	GroupDescription string
}

// @title         InitGlobal
// @description   初始化全局数据的文件夹以及全局数据库
// @auth          DataEraserC                                   (2024/2/17   21:54)
// @param         GlobalPath                            string              "指定数据存放在什么地方"
// @param         SafeMode                              bool                "是否自动迁移数据库模型"
// @return        GlobalDatabase                        *gorm.DB            "全局数据库"
// @return        err                                   error               "可能存在的错误"
func InitGlobal(GlobalPath string, SafeMode bool) (*gorm.DB, error) {
	// 初始化GlobalPath文件夹
	if _, err := os.Stat(GlobalPath); os.IsNotExist(err) {
		// GlobalPath不存在，创建GlobalPath
		err := os.MkdirAll(GlobalPath, 0755)
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

	// 连接数据库
	GlobalDatabase, err := gorm.Open(sqlite.Open(GlobalPath+"/database.db"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")

	}
	if SafeMode {
		// AutoMigrate 自动迁移数据库
		err = GlobalDatabase.AutoMigrate(&UserInfo{}, &Login{}, &Token{}, &CreateGroupRequest{}, &GroupInfo{})
		if err != nil {
			return nil, errors.New("failed to AutoMigrate database")
		}
	}
	return GlobalDatabase, nil
}

// @title         generateToken
// @description   生成token的函数
// @auth          DataEraserC              (2024/2/17   21:54)
// @param         UserID           uint                "指定需要生成token的UserID"
// @param         secretKey        string              "指定用于生成token的密钥"
// @return        tokenString      string              "令牌字符串"
func generateToken(UserID uint, secretKey string) string {
	// 创建一个Token对象
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置Token的自定义声明
	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = UserID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 设置Token的过期时间

	// 使用密钥对Token进行签名，生成最终的Token字符串
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}

// @title         parseToken
// @description   解析token的函数
// @auth          DataEraserC              (2024/2/17   21:54)
// @param         tokenString      string              "指定需要解析的Token"
// @param         secretKey        string              "指定用于解析token的密钥"
// @return        claims           jwt.MapClaims       "解析获得的声明对象键值对"
// @return        err              error               "可能存在的错误"
func parseToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	// 解析Token字符串
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token的签名方法是否有效
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("无效的签名方法：%v", token.Header["alg"])
	}

	// 返回Token中的声明部分
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的Token")
}

// @title         Login_account_password
// @description   处理用户名密码登陆入口的函数
// @auth          DataEraserC                           (2024/2/17   21:54)
// @param         GlobalDatabase                *gorm.DB            "全局数据库"
// @return        匿名函数                      gin.HandlerFunc     "gin消息中间件"
func Login_account_password(GlobalDatabase *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Username string `json:"Username"`
			Password string `json:"Password"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "参数错误"})
			return
		}

		var login Login
		if err := GlobalDatabase.Where("username = ?", request.Username).First(&login).Error; err != nil {
			c.JSON(400, gin.H{"code": 2, "message": "用户名或密码错误"})
			return
		}

		if login.Password != request.Password {
			c.JSON(400, gin.H{"code": 2, "message": "用户名或密码错误"})
			return
		}

		token := Token{
			UserID: login.UserID,
			Token:  generateToken(login.UserID, JWTSecretKey),
		}
		if err := GlobalDatabase.Create(&token).Error; err != nil {
			c.JSON(400, gin.H{"code": 4, "message": "无法生成token"})
			return
		}

		c.JSON(200, gin.H{"code": 0, "message": "登录成功", "Token": token.Token, "UserID": token.UserID})
	}
}

// unfinished

// @title         Login_wx
// @description   处理微信登陆入口的函数
// @auth          DataEraserC                    (2024/2/17   21:54)
// @param         GlobalDatabase         *gorm.DB            "全局数据库"
// @param         WXAppID                string              "微信小程序AppID"
// @param         WXAppSecret            string              "微信小程序AppSecret"
// @return        匿名函数               gin.HandlerFunc     "gin消息中间件"
func Login_wx(GlobalDatabase *gorm.DB, WXAppID string, WXAppSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		var request struct {
			JsCode string `json:"code"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "参数错误"})
			return
		}
		wxLoginResp, err := _WXLogin(request.JsCode, WXAppID, WXAppSecret)
		if err != nil {
			c.JSON(400, gin.H{"code": 2, "message": "内部错误"})
			return
		}

		var login Login
		if err := GlobalDatabase.Where("openid = ?", wxLoginResp.OpenId).First(&login).Error; err != nil {
			//未注册
			var user UserInfo
			if err := GlobalDatabase.Create(&user).Error; err != nil {
				c.JSON(400, gin.H{"code": 2, "message": "内部错误"})
				return
			}
			// username password unfinished
			login = Login{OpenID: wxLoginResp.OpenId, UserID: user.ID}
		}
		token := Token{
			UserID: login.UserID,
			Token:  generateToken(login.UserID, JWTSecretKey),
		}
		c.JSON(200, gin.H{"code": 0, "message": "登录成功", "Token": token.Token, "UserID": token.UserID})
	}
}

// @title         _WXLogin
// @description   处理微信登陆的函数
// @auth          DataEraserC               (2024/2/17   21:54)
// @param         code              string              "微信小程序前端获得的jscode"
// @param         AppID             string              "微信小程序AppID"
// @param         AppSecret         string              "微信小程序AppSecret"
// @return        wxResp            *WXLoginResp        "微信登陆返回值json对象"
// @return        err               error               "可能存在的错误"
func _WXLogin(code string, AppID string, AppSecret string) (*WXLoginResp, error) {

	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	url = fmt.Sprintf(url, AppID, AppSecret, code)

	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return nil, err
	}

	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("ErrCode:%d  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
	}

	return &wxResp, nil
}

func Userinfo(GlobalDatabase *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Token  string `json:"Token"`
			UserID string `json:"UserID"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "参数错误"})
			return
		}

		var tokenRecord Token
		if err := GlobalDatabase.Where("token = ?", request.Token).First(&tokenRecord).Error; err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "身份验证失败"})
			return
		}

		var user UserInfo
		if err := GlobalDatabase.First(&user, tokenRecord.UserID).Error; err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "获取用户信息失败"})
			return
		}

		c.JSON(200, gin.H{"code": 0, "message": "获取个人信息成功", "data": user})

	}
}

// Updateuserinfo 更新用户信息，仅更新请求中包含的数据，不更新为空的字段
func Updateuserinfo(GlobalDatabase *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Token              string
			UserID             uint
			Avatar             *string
			Name               *string
			NickName           *string
			Gender             *string
			Collage            *string
			Majar              *string
			Grade              *uint
			PhoneNumber        *string
			RegistrationNumber *string
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "参数错误"})
			return
		}

		var tokenData Token
		if err := GlobalDatabase.Model(&tokenData).Where("token = ?", request.Token).First(&tokenData).Error; err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "身份验证失败"})
			return
		}

		if tokenData.UserID != request.UserID {
			c.JSON(400, gin.H{"code": 3, "message": "无权限修改别人的信息"})
			return
		}

		updateData := make(map[string]interface{})
		if request.Avatar != nil {
			updateData["Avatar"] = *request.Avatar
		}
		if request.Name != nil {
			updateData["Name"] = *request.Name
		}
		if request.NickName != nil {
			updateData["NickName"] = *request.NickName
		}
		if request.Gender != nil {
			updateData["Gender"] = *request.Gender
		}
		if request.Collage != nil {
			updateData["Collage"] = *request.Collage
		}
		if request.Majar != nil {
			updateData["Majar"] = *request.Majar
		}
		if request.Grade != nil {
			updateData["Grade"] = *request.Grade
		}
		if request.PhoneNumber != nil {
			updateData["PhoneNumber"] = *request.PhoneNumber
		}
		if request.RegistrationNumber != nil {
			updateData["RegistrationNumber"] = *request.RegistrationNumber
		}

		var user UserInfo
		if err := GlobalDatabase.Model(&user).Where("ID = ?", tokenData.UserID).Updates(updateData).Error; err != nil {
			c.JSON(500, gin.H{"code": 2, "message": "修改个人信息失败"})
			return
		}

		c.JSON(200, gin.H{"code": 0, "message": "修改个人信息成功"})
	}
}

func Logout(GlobalDatabase *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Token string
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "参数错误"})
			return
		}

		var tokenData Token
		if err := GlobalDatabase.Model(&tokenData).Where("token = ?", request.Token).First(&tokenData).Error; err != nil {
			c.JSON(400, gin.H{"code": 1, "message": "身份验证失败"})
			return
		}

		_DeleteTokenByToken(GlobalDatabase, tokenData.Token)

	}
}

// 删除某用户的所有Token (在修改密码时需要用到)
func _DeleteTokensByUserID(GlobalDatabase *gorm.DB, userID uint) error {
	result := GlobalDatabase.Where("user_id = ?", userID).Delete(&Token{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 需要的参数TokenDatabase指存放Token的Database
// 用单个token搜索并删除该token
func _DeleteTokenByToken(TokenDatabase *gorm.DB, token string) error {
	result := TokenDatabase.Where("token = ?", token).Delete(&Token{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 需要的参数TokenDatabase指存放Token的Database
func _GetUserIDByToken(TokenDatabase *gorm.DB, token string) (uint, error) {
	var tokenData Token
	if err := TokenDatabase.Model(&tokenData).Where("token = ?", token).First(&tokenData).Error; err != nil {
		return 0, errors.New("无法找到token对应的UserID")
	}
	return tokenData.UserID, nil
}
