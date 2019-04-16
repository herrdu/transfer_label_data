#### 标注数据转换工具

###### Go 环境安装

1. brew install go

###### Go 依赖

1. go get github.com/json-iterator/go

###### 修改目标目录

打开 main.go 文件，修改 filePath 的值为导出标注结果的文件路径

###### 执行代码

在 main.go 目录里 ，执行 go run main.go 命令

###### 生成结果

程序执行完毕之后，生成 new_result.json 文件为处理后的新结果
