# 项目结构

```plain
├── data                             # [运行时]生成的全局数据文件夹
│   ├── database.db                  # [运行时]生成的全局数据库
│   ├── group                        # [运行时]生成的组织数据文件夹
│   │   ├── 1                        # [运行时]生成的组织group1的数据文件夹
│   │   │   ├── database.db          # [运行时]生成的组织group1的数据库
│   │   │   └── meeting              # [运行时]生成的组织group1的会议数据文件夹
│   │   │       ├── 1                # [运行时]生成的组织group1的会议meeting1数据文件夹
│   │   │       │   └── database.db  # [运行时]生成的组织group1的会议meeting1数据库
│   │   │       ├── 2                # [运行时]生成的组织group1的会议meeting1数据文件夹
│   │   │       │   └── database.db  # [运行时]生成的组织group1的会议meeting2数据库
│   │   │       └── 3                # [运行时]生成的组织group1的会议meeting3数据文件夹
│   │   │           └── database.db  # [运行时]生成的组织group1的会议meeting3数据库
│   │   └── 2                        # [运行时]生成的组织group2的数据文件夹
│   │       ├── database.db          # [运行时]生成的组织group2的数据库
│   │       └── meeting              # [运行时]生成的组织group2的会议数据文件夹
│   │           └── 1                # [运行时]生成的组织group2的会议meeting1数据文件夹
│   │               └── database.db  # [运行时]生成的组织group2的会议meeting1数据库
│   └── user                         # [运行时]生成的用户数据文件夹
│       ├── 1                        # [运行时]生成的用户user1的数据文件夹
│       │   └── database.db          # [运行时]生成的用户user1的数据库
│       └── 2                        # [运行时]生成的用户user2的数据文件夹
│           └── database.db          # [运行时]生成的用户user1的数据库
├── logs                             # [运行时]生成的日志目录
│   └── log.log                      # [运行时]生成的日志文件
├── docs                             #* 文档
│   ├── BuildInstructions.md         # 编译教程
│   ├── CodeDesign.md                # 代码设计
│   ├── database.md                  # 数据库设计
│   ├── Interface.md                 #* 接口文档
│   └── struct.md                    #* 项目结构
├── default.nix                      # nix 项目管理 (存放项目名称等)
├── flake.nix                        # nix 表达式 用于生成合适的环境shell
├── shell.nix                        # nix 表达式 用于生成合适的环境shell
├── flake.lock                       # nix 版本管理
├── gomod2nix.toml                   # nix下类似go.mod的东西
├── go.sum                           #* 依赖的 module 的校验信息
├── go.mod                           #* 依赖库以及依赖库的版本
├── main.go                          * 主程序
├── global.go                        # global子模块的代码
├── group.go                         # group子模块的代码
├── meeting.go                       # meeting子模块的代码
├── secrets.go                       # 密钥变量存储
├── user.go                          # user子模块的代码
├── README.md                        #* 项目说明文件
└── Makefile                         # 编译项目用的脚本等
```
> 带`*`表示需要关注的部分(随项目推进需要修改的部分)(重要程度上升) 带`#`表示重要程度下降
