package main

import (
	_ "demo/handle_register_demo/domain" // init 函数被统一调用
	"demo/handle_register_demo/manage_ins"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func handle(c *gin.Context) {
	id := c.Query("id")

	// 从全局实例注册器里获取实例
	instanceA, err := manage_ins.GetInstance(manage_ins.EnumType(cast.ToInt64(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	instanceA.Handle(c.Request.Context())
}

func main() {
	r := gin.Default()

	// curl 127.0.0.1:8080/instance?id=1
	// curl 127.0.0.1:8080/instance?id=2
	r.GET("/instance", handle)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
