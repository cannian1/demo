package main

import (
	"demo/dtm_demo/tcc/dto"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/try", apiTry)

	r.POST("/confirm", apiConfirm)

	r.POST("/cancel", apiCancel)

	r.Run(":8080")
}

var Balance int64 = 8000
var LockBalance int64 = 0

func apiTry(c *gin.Context) {
	var req dto.ConsumeReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(409, "")
	}

	Balance -= req.Amount
	LockBalance += req.Amount

	// 校验要后置，否则会影响回滚补偿的正确性
	if Balance < 0 {
		c.JSON(409, gin.H{"dtm_result": "FAILURE", "message": "余额不足"})
		return
	}

	fmt.Println("[apiTry] Balance:", Balance, " LockBalance:", LockBalance)
	c.JSON(200, gin.H{"dtm_result": "try SUCCESS"})
}

func apiConfirm(c *gin.Context) {
	var req dto.ConsumeReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(409, "")
	}

	// confirm 阶段报错会被 dtm 重试
	//c.JSON(409, "我是故意报错的")
	//return

	LockBalance -= req.Amount
	fmt.Println("[apiConfirm] Balance:", Balance, " LockBalance:", LockBalance)
	c.JSON(200, gin.H{"dtm_result": "confirm SUCCESS"})
}

func apiCancel(c *gin.Context) {
	var req dto.ConsumeReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(409, "")
	}
	LockBalance -= req.Amount
	Balance += req.Amount
	fmt.Println("[apiCancel] Balance:", Balance, " LockBalance:", LockBalance)
	c.JSON(200, gin.H{"dtm_result": "cancel SUCCESS"})
}
