package mock_demo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi")
	}))
	defer ts.Close()
	addr := ts.URL
	GetUser(addr)

	// GetUser("https://www.baidu.com")
}
