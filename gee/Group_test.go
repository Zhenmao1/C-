package gee

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testMiddleWare() HandlerFunc {
	return func(c *Context) {
		now := time.Now()
		//c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("局部中间件[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(now))
	}
}
func TestMiddleWare(t *testing.T) {
	r := New()
	r.Use(Logger())
	v1 := r.NewGroup("/v1")
	{
		v1.Use(testMiddleWare())
		v1.GET("/", func(c *Context) {
			c.HTML(http.StatusOK, "hu", "<h1>Hello hu</h1>")
		})

		v1.GET("/hello/:name", func(c *Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	// 创建一个模拟的 HTTP 请求
	req, err := http.NewRequest("GET", "/v1/hello/hu", nil)
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

func TestGroup(t *testing.T) {
	r := New()
	r.GET("/index", func(c *Context) {
		c.HTML(http.StatusOK, "hu", "<h1>Index Page</h1>")
	})
	v1 := r.NewGroup("/v1")
	{
		v1.GET("/", func(c *Context) {
			c.HTML(http.StatusOK, "hu", "<h1>Hello hu</h1>")
		})

		v1.GET("/hello/hu", func(c *Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.NewGroup("/v2")
	{
		v2.GET("/hello/:name", func(c *Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *Context) {
			c.JSON(http.StatusOK, H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}
	traverseHandlers(v1.engine.router.handlers)
	traverseNode(v1.engine.router.root["GET"])
	traverseHandlers(v2.engine.router.handlers)
	traverseNode(v2.engine.router.root["GET"])
	traverseNode(v2.engine.router.root["POST"])

	n, ps := v2.engine.router.Getroute("GET", "/v2/hello/hu")

	if n == nil {
		t.Fatal("不应返回空值")
	}

	if n.pattern != "/v2/hello/:name" {
		t.Fatal("应匹配 /hello/:name")
	}

	if ps["name"] != "hu" {
		t.Fatal("name 应该等于 'hu'")
	}

	fmt.Printf("匹配的路径: %s, 参数['name']: %s\n", n.pattern, ps["name"])
}
