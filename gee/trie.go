package gee

import (
	"strings"
)

type node struct {
	pattern  string  //完整路径
	part     string  //当前待匹配字符
	children []*node //该节点的所有子节点
	isWild   bool    //是否精确匹配
}

// 插入操作的查找函数
// 匹配节点，只要有一个孩子节点匹配，迅速知道该路由存在
// 不再需要向下查找了
func (n *node) matchInsert(part string) *node {
	for _, cur := range n.children {
		if cur.part == part || cur.isWild {
			return cur
		}
	}
	return nil
}

// 匹配操作的查找函数
// 匹配节点，需要返回所有匹配的节点
func (n *node) matchAll(part string) []*node {
	childs := make([]*node, 0)
	for _, cur := range n.children {
		if cur.part == part || cur.isWild {
			childs = append(childs, cur)
		}
	}
	return childs
}

// 插入函数
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	match := n.matchInsert(part)
	if match == nil {
		match = &node{
			part:   part,
			isWild: (part[0] == '*') || (part[0] == ':')}
		n.children = append(n.children, match)
	}
	match.insert(pattern, parts, height+1)
}

// 查找函数
func (n *node) search(parts []string, height int) *node {
	//所有的part都匹配，或者最后一个是*
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchAll(part)
	for _, cur := range children {
		res := cur.search(parts, height+1)
		if res != nil {
			return res
		}
	}
	return nil
}
