package gee

import (
	"fmt"
	"reflect"
	"testing"
)

// 创建测试路由
func newTestRouter() *router {
	r := newRouter()
	r.AddRouter("GET", "/", nil)
	r.AddRouter("GET", "/hello/:name", nil)
	r.AddRouter("GET", "/hello/b/c", nil)
	r.AddRouter("GET", "/hi/:name", nil)
	r.AddRouter("GET", "/assets/*filepath", nil)
	return r
}

func traverseHandlers(handlers map[string]HandlerFunc) {
	for key, value := range handlers {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}
}
func traverseNode(root *node) {
	fmt.Printf("node %v\n", root)
	// 如果 node 结构有其他字段，也可以在这里进行遍历
	for _, child := range root.children {
		traverseNode(child)
	}

}
func TestAddRouter(t *testing.T) {
	r := newTestRouter()
	traverseHandlers(r.handlers)
	traverseNode(r.root["GET"])
}

// 测试解析模式函数
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("测试解析模式函数失败")
	}
}

// 测试获取路由函数
func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.Getroute("GET", "/hello/hu")

	if n == nil {
		t.Fatal("不应返回空值")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("应匹配 /hello/:name")
	}

	if ps["name"] != "hu" {
		t.Fatal("name 应该等于 'hu'")
	}

	fmt.Printf("匹配的路径: %s, 参数['name']: %s\n", n.pattern, ps["name"])
}
