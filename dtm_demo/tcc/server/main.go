package main

import (
	"demo/dtm_demo/tcc/dto"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/try", apiTry)

	r.POST("/confirm", apiConfirm)

	r.POST("/cancel", apiCancel)

	r.Run(":8080")
}

var (
	Balance           int64 = 8000
	LockBalance       int64 = 0
	mu                sync.Mutex
	transactionStatus = make(map[string]string) // 记录每个GID的状态：try/confirm/cancel
)

func apiTry(c *gin.Context) {
	gid := c.Query("gid")
	var req dto.ConsumeReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// 幂等检查：如果之前已处理过相同GID的Try，直接返回之前的结果
	if status, exists := transactionStatus[gid]; exists {
		if status == "try" {
			c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
		} else {
			c.JSON(409, gin.H{"dtm_result": "FAILURE"})
		}
		return
	}

	if Balance < req.Amount {
		c.JSON(409, gin.H{"dtm_result": "FAILURE", "message": "余额不足"})
		return
	}

	Balance -= req.Amount
	LockBalance += req.Amount
	transactionStatus[gid] = "try" // 记录状态

	fmt.Printf("[apiTry] GID=%s, Balance=%d, LockBalance=%d\n", gid, Balance, LockBalance)
	c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
}

func apiConfirm(c *gin.Context) {
	gid := c.Query("gid")
	var req dto.ConsumeReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// 幂等检查：如果已Confirm，直接返回成功
	if status, exists := transactionStatus[gid]; exists {
		if status == "confirm" {
			c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
			return
		}
	} else {
		// 没有Try记录，无法Confirm
		c.JSON(409, gin.H{"dtm_result": "FAILURE"})
		return
	}

	LockBalance -= req.Amount
	transactionStatus[gid] = "confirm"
	fmt.Printf("[apiConfirm] GID=%s, Balance=%d, LockBalance=%d\n", gid, Balance, LockBalance)
	c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
}

func apiCancel(c *gin.Context) {
	gid := c.Query("gid")
	var req dto.ConsumeReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// 幂等检查：如果已Cancel，直接返回成功
	if status, exists := transactionStatus[gid]; exists {
		if status == "cancel" {
			c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
			return
		}
	} else {
		// 没有Try记录，无法Cancel
		c.JSON(409, gin.H{"dtm_result": "FAILURE"})
		return
	}

	// 只有Try状态才需要回滚余额
	if transactionStatus[gid] == "try" {
		LockBalance -= req.Amount
		Balance += req.Amount
	}
	transactionStatus[gid] = "cancel"
	fmt.Printf("[apiCancel] GID=%s, Balance=%d, LockBalance=%d\n", gid, Balance, LockBalance)
	c.JSON(200, gin.H{"dtm_result": "SUCCESS"})
}
