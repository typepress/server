# server

HTTP Static Server module base on Martini.

基于 Martini 的 HTTP 静态文件服务器模块, 不包含路由. 只支持以下功能:

 - 最基础的命令行参数
 - os.Getenv 获取命令行参数
 - 支持 TOML 配置文件
 - 安全关闭机制
 - 静态文件输出, 预压缩
 - 预置 i18n 接口

虽然只是最基本的功能, 但确实是一个完整的架构.

# Usage

Simple:

```go
package main

import "github.com/typepress/server"

func main() {
	server.Simple()
}

```

Run:

```go
package main

import (
	"github.com/typepress/core"
	"github.com/typepress/server"
)

func main() {

	// *martini.Martini, martini.Router
	m, r := core.Martini()

	// stopSignal for stop server safe, Usage:
	// core.FireSignal(stopSignal)

	stopSignal := server.StopSignal()

	// something

	err := server.Run(m, r)
	println(err)
}
```

License
=======
BSD-2-Clause