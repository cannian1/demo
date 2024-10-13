package mock_demo

import (
	"fmt"
	"io"
	"net/http"
)

// GetUser 对 addr 进行 get 请求
func GetUser(addr string) {
	resp, err := http.Get(addr)
	if err != nil {
		fmt.Println("http get error:", err)
		return
	}
	defer resp.Body.Close()
	all, _ := io.ReadAll(resp.Body)
	fmt.Println(string(all))
}
