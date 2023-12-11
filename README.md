
## 开发安装
* git clone 
* 切换到dev开发分支，进入项目根目录
* 安装依赖包，执行 ``go mod tidy``
* 找到 etc/appbuilder-example.yaml 配置文件，复制一份为etc/appbuilder-api.yaml, 修改配置项
* 创建数据库 appbuilder


## 启动API服务
* 进入项目根目录
* 运行 `go run appbuilder.go` 或者 `go run appbuilder.go -f etc/appbuilder-api.yaml`
* 服务地址:http://{ip}:8884


