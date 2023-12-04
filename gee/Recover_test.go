package gee

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecover(t *testing.T) {
	r := New()
	r.Use(Recover())
	r.GET("/", func(c *Context) {
		c.String(http.StatusOK, "Hello hu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *Context) {
		names := []string{"hu"}
		c.String(http.StatusOK, names[100])
	})
	// 创建一个模拟的 HTTP 请求
	req, err := http.NewRequest("GET", "/panic", nil)
	if err != nil {
		t.Fatal(err)
	}
	// 创建一个 ResponseRecorder，用于记录响应
	w := httptest.NewRecorder()
	// 调用路由的 ServeHTTP 方法，传入模拟的请求和响应
	r.ServeHTTP(w, req)
	// 检查响应结果
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}
