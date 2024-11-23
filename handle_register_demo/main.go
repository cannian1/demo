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

	instanceA, err := manage_ins.GetInstance(manage_ins.EnumType(cast.ToInt64(id)))
	if err != nil {
		fmt.Println(err)
		return
	}
	instanceA.Handle(c.Request.Context())
}

func main() {

	r := gin.Default()
	r.GET("/instance", handle)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
