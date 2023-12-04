package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	root     map[string]*node
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		root:     make(map[string]*node),
	}
}

// 把字符串进行分割，且最多只保留一个*
func parsePattern(s string) []string {
	ret := make([]string, 0)
	parts := strings.Split(s, "/")

	for _, item := range parts {
		if item != "" {
			ret = append(ret, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return ret
}

// AddRouter 添加一个路由规则到路由器中
func (r *router) AddRouter(method string, pattern string, handler HandlerFunc) {
	// 打印路由信息，包括请求方法和路由模式
	log.Printf("Route %4s - %s", method, pattern)
	// 解析路由模式字符串，得到路由模式的各个部分
	parts := parsePattern(pattern)
	// 构建路由键，格式为 "请求方法-路由模式"
	key := method + "-" + pattern
	// 检查是否已存在请求方法对应的根节点，如果不存在，则创建一个
	_, ok := r.root[method]
	if !ok {
		r.root[method] = &node{}
	}
	// 在根节点中插入路由信息，包括路由模式的各个部分和处理函数
	r.root[method].insert(pattern, parts, 0)
	// 将处理函数与路由键关联，存储到路由器的 handlers 映射中
	r.handlers[key] = handler
}

func (r *router) Getroute(method string, path string) (*node, map[string]string) {
	// 解析路径，将其拆分成多个部分
	searchParts := parsePattern(path)

	// 初始化参数映射，用于存储从路径中提取的参数
	params := make(map[string]string)

	// 获取指定HTTP方法的路由树的根节点
	root, ok := r.root[method]
	if !ok {
		// 如果没有找到指定HTTP方法的根节点，返回节点和参数均为nil
		return nil, nil
	}

	// 在路由树中执行搜索，使用解析的路径部分
	n := root.search(searchParts, 0)

	if n != nil {
		// 如果找到匹配的节点，解析节点的模式
		parts := parsePattern(n.pattern)

		// 遍历模式的各个部分，提取动态参数
		for index, part := range parts {
			// 如果是动态参数（以':'开头），将其值存储在参数映射中
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}

			// 如果是通配符参数（以'*'开头且长度大于1），将其值存储在参数映射中
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				// 跳出循环，因为通配符参数后面的所有部分已经被处理
				break
			}
		}

		// 返回找到的节点和提取的参数
		return n, params
	}

	// 如果未找到匹配的节点，返回节点和参数均为nil
	return nil, nil
}

func (r *router) Handle(c *Context) {
	key := c.Method + "-" + c.Path
	n, params := r.Getroute(c.Method, c.Path)
	if handler, ok := r.handlers[key]; ok && n != nil {
		c.Params = params
		c.handlers = append(c.handlers, handler)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
