<div align=left>
<img src="https://img.shields.io/badge/golang-1.16-blue"/>
<img src="https://img.shields.io/badge/gin-1.7.2-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-1.21.12-red"/>
</div>


# 项目文档

## 1. 基本介绍

### 1.1 项目介绍

> gofast是一个基于 [GVA](https://github.com/flipped-aurora/gin-vue-admin) 的后端部分进行修改的代码生成项目，可以通过sql文件生成相应的crud的基本业务代码，减少开发量。
### 1.2 代码结构
```
.
├── cmd
│   ├── gf
│   │   ├── initFile.go
│   │   └── root.go
│   └── main.go
├── core
│   └── zap.go
├── gen_code.go
├── global
│   └── global.go
├── go.mod
├── go.sum
├── initialize
│   ├── gorm.go
│   └── internal
│       └── logger.go
├── middleware
│   ├── cors.go
│   ├── error.go
│   └── loadtls.go
├── README.md
├── resource
│   ├── static
│   │   ├── api_health.static
│   │   ├── request.static
│   │   └── response.static
│   └── template
│       ├── api.go.tpl
│       ├── health.go.tpl
│       ├── initRouter.go.tpl
│       ├── main.go.tpl
│       ├── model.go.tpl
│       ├── request.go.tpl
│       ├── router.go.tpl
│       └── service.go.tpl
└── utils
    ├── auto_code.go
    ├── directory.go
    ├── file_operation.go
    └── rotatelogs.go
```
| 文件夹       | 说明                    | 描述                        |
| ------------ | ----------------------- | --------------------------- |
| `cmd`        | 命令行工具层                   | 用来生成代码的命令行工具 |
| `api`        | api层                   | api层 |
| `--v1`       | v1版本接口              | v1版本接口                  |
| `core`       | 核心文件                | 核心组件zap的初始化 |
| `docs`       | swagger文档目录         | swagger文档目录 |
| `global`     | 全局对象                | 全局对象 |
| `initialize` | 初始化 | gorm的初始化 |
| `--internal` | 初始化内部函数 | gorm 的 longger 自定义,在此文件夹的函数只能由 `initialize` 层进行调用 |
| `middleware` | 中间件层 | 用于存放 `gin` 中间件代码 |
| `resource`   | 静态资源文件夹          | 负责存放静态文件                |
| `--static` | 静态文件 | 不需要模板渲染的静态文件复制 |
| `--template` | 模板 | 模板文件夹,存放的是代码生成器的模板 |
| `utils`      | 工具包                  | 工具函数封装，包括自动生成代码的逻辑，文件操作的逻辑            |
## 2. 使用说明

```
- golang版本 >= v1.14
- IDE推荐：vscode
- 暂时只支持linux
- 只支持module模式开发
- 只支持mysql数据库
```

### 2.1 初始化

确保GO111MODULE参数是ON状态，或者是auto状态
确保已经配置正确的GOPROXY地址

```bash
# 创建项目文件夹(以demoProject为例)
mkdir demoProject
```
```bash
# 进入创建好的项目
cd demoProject
```
```bash
# 新建main.go文件
touch main.go
# 写入以下代码：(项目路径需要替换为自己的)
package main

import "github.com/pyx-py/gofast"

func main() {
	gofast.NewCodeTool("{{项目的绝对路径，比如/root/user/demoProject}}", "")
}
```
```bash
# 在项目下初始化module, 例子中以demo作为module名称(以下的demo都是代值module)，使用时也需要更改为自己的module
go mod init demo
```
```bash
# 下载gofast包, 执行：
go mod tidy
```
```bash
# 生成命令行工具，执行:
go run main.go
# 此时项目目录下多了一个gf的可执行文件，初始化完毕
```


### 2.2 代码生成工具的使用

```bash
# 此时项目文件夹下已经有了gf文件，可查看帮助文档
./gf -h
# gf工具的命令只有一个init，可查看帮助文档
./gf init -h
# 此命令下的flag含义如下：
-m, --module : module名称，必传
-p, --path : 项目目录路径，必传
-s, --sql : sql文件路径，非必传
-l, --log : 日志存放目录，非必传
-f, --gofast: 下载的gofast路径，非必传
-c, --column : 生成代码的列表接口的搜索字段，非必传

# 此时可执行
./gf init -m demo -p /root/user/demoProject
# 即可生成最基础的代码框架(不含crud代码和日志包)
# 再次拉取依赖包
go mod tidy

# 可再次执行含有sql的命令，例如./gf init -m demo -p /root/user/demoProject -s ./t_band.sql就能生成业务代码
# sql文件格式示例：
CREATE TABLE `t_band` (
                           `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                           `created_at` datetime DEFAULT NULL,
                           `updated_at` datetime DEFAULT NULL,
                           `deleted_at` datetime DEFAULT NULL,
                           `band` varchar(128) NOT NULL COMMENT '',
                           `band_pinyin` varchar(128) NOT NULL COMMENT '',
                           `band_name` varchar(128) NOT NULL COMMENT '',
                           PRIMARY KEY (`id`),
                           KEY `idx_band` (`band`),
                           KEY `idx_band_py` (`band_py`),
                           KEY `idx_band_name` (`band_name`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb4
# 此时main.go文件已经被覆盖，只需要手动添加48行的数据库Dsn地址，再执行
go run main.go
# 项目成功运行
```

### 2.3 swagger自动化API文档

#### 2.3.1 安装 swagger

````
go get -u github.com/swaggo/swag/cmd/swag
````




#### 2.3.2 生成API文档
确保swag已经加入到全局的可执行环境变量，否则会报找不到命令
> 比如在～/.bashrc文件中添加golang的bin文件路径,写入一行(具体路径根据个人修改)

> export PATH="/root/go/bin:$PATH"

```` bash
# 执行swag初始化
swag init
````

> 执行上面的命令后，项目目录下会出现docs文件夹里 `docs.go`, `swagger.json`, `swagger.yaml` 三个文件，再次启动go服务之后, 在浏览器输入 [http://localhost:8888/swagger/index.html](http://localhost:8888/swagger/index.html) 即可查看swagger文档。

> 如果出现页面错误`Failed to load spec`,后台也出现panic的话，需要在initializa/router.go的import中添加导入 `_ "demo/docs"`，再次go run main.go即可成功。

