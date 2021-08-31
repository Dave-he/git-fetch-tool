
# git-fetch-tool
## 工程简介
    使用go搭建的异步快速拉取git工程的工具，
    用于快速更新当前目录下所有git目录
## 安装步骤
- 编译main.go生成二进制文件 `go build -o [filename]`
- 将二进制文件移入`$GOPATH`下或者配置环境变量 
- 到工作目录打开控制台运行: `[filename]` 或者 `[filename] /目录`
