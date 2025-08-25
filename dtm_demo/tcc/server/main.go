package main

import (
	"fmt"

	"demo/dtm_demo/tcc/dto"

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
	if Balance < req.Amount {
		c.JSON(200, gin.H{"dtm_result": "FAILURE", "message": "余额不足"})
		return
	}

	Balance -= req.Amount
	LockBalance += req.Amount
	fmt.Println("apiTry")
	c.JSON(200, gin.H{"dtm_result": "try SUCCESS"})
}

func apiConfirm(c *gin.Context) {
	var req dto.ConsumeReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(409, "")
	}
	LockBalance -= req.Amount
	fmt.Println("apiConfirm")
	c.JSON(200, gin.H{"dtm_result": "confirm SUCCESS"})
}

func apiCancel(c *gin.Context) {
	var req dto.ConsumeReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(409, "")
	}
	Balance += req.Amount
	LockBalance -= req.Amount
	fmt.Println("apiCancel")
	c.JSON(200, gin.H{"dtm_result": "cancel SUCCESS"})
}
