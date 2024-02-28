# 接口文档

> 大部分请求必须提供UserID和Token,若前端未能保存UserID信息可向后端发起请求获取并存到前端的LocalStorage

## 用户登录接口(账号密码)

接口地址：/login_account_password

请求方法：POST

请求参数：

- Username：用户名称，类型为字符串
- Password：用户密码，类型为字符串

请求示例：

```http
POST /login_account_password
Content-Type: application/json

{
    "Username": "testuser",
    "Password": "123456"
}
```

返回数据：

- code：返回状态码，0 表示成功，非0 表示失败
- message：返回信息，登录成功或失败的提示信息
- Token：用户登录后生成的令牌，类型为字符串
- UserID：用户ID，类型为integer

成功返回示例：

```json
{
  "code": 0,
  "message": "登录成功",
  "Token": "abcd1234",
  "UserID": 1
}
```

## 用户登录接口(微信接口api) (第一次登陆需要使用这个接口)

接口地址：/login_wx

请求方法：POST

请求参数：

- code ：wx.login()得到的code

请求示例：

```http
POST /login_wx
Content-Type: application/json

{
    "code": "123456",
}
```

返回数据：

- code：返回状态码，0 表示成功，非0 表示失败
- message：返回信息，登录成功或失败的提示信息
- token：用户登录后生成的令牌，类型为字符串
- UserID：用户ID，类型为integer

成功返回示例：

```json
{
  "code": 0,
  "message": "登录成功",
  "Token": "abcd1234",
  "UserID": 1
}
```

## 用户获取个人信息接口

接口地址：/userinfo

请求方法：POST

请求参数：

- Token：用户登录后生成的令牌，类型为字符串
- UserID：用户ID，类型为integer

请求示例：

```http
POST /userinfo
Content-Type: application/json

{
    "Token": "abcd1234",
    "UserID" : 1
}
```

返回数据：

- code：返回状态码，0 表示成功，非0 表示失败
- message：返回信息，获取个人信息成功或失败的提示信息
- data：用户个人信息的数据，包含用户名

成功返回示例：

```json
{
  "code": 0,
  "message": "获取个人信息成功",
  "data": {
    "ID": 1,
    "Username": "username1",
    "Password": "",
    "Avatar": "/default_avatar.png",
    "Name": "Name1",
    "PhoneNumber": "123456789",
    "RegistrationNumber": "20221002122",
    "Permission": "老师"
  }
}
```

## 用户修改个人信息接口

接口地址：/updateuserinfo

请求方法：POST

请求参数：

- Token：用户登录后生成的令牌，类型为字符串
- ID : 需要修改的用户,正常情况只能修改自己的信息,有权限时可以修改他人信息(不给出ID默认修改自己)
- Username：修改后的用户名，类型为字符串
- Password：修改后的密码，类型为字符串
- Avatar：修改后的头像链接，类型为字符串
- Name：正常用户不提供姓名修改,教师管理时可以修改其他人姓名
- PhoneNumber：修改后的手机号，类型为字符串
- RegistrationNumber : 正常用户不提供学号修改,教师管理时可以修改其他人学号
- Permission : 管理员可修改他人权限

请求示例：

```http
POST /updateuserinfo
Content-Type: application/json

{
    "Token": "abcd1234",
    "ID" : 1
    "Username": "newusername",
    "Password": "newpassword",
    "Avatar": "new_avatar.jpg",
    "PhoneNumber": "1234567890",
    "RegistrationNumber": "20221002111",
    "Permission": "老师"
}
```

返回数据：

- code：返回状态码，0 表示成功，非0 表示失败(后端应该先确认身份是否足够修改全部 否则一个都不能修改)
- message：返回信息，修改个人信息成功或失败的提示信息

成功返回示例：

```json
{
  "code": 0,
  "message": "修改个人信息成功"
}
```

## 用户注销登陆

接口地址：/logout

请求方法：POST

请求参数：

- Token：用户登录后生成的令牌，类型为字符串
- UserID：用户ID，类型为integer

请求示例：

```http
POST /updateuserinfo
Content-Type: application/json

{
    "Token": "abcd1234",
    "UserID": 1
}
```

返回数据：

- code：返回状态码，0 表示成功，非0 表示失败
- message：返回信息，注销成功或失败的提示信息

成功返回示例：

```json
{
  "code": 0,
  "message": "注销成功"
}
```
