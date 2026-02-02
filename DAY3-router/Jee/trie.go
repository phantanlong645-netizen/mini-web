package Jee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  //完整的路由路径（如 /p/:lang/doc）。只有叶子节点（一条路径的终点）才会存这个值，中间节点为空。
	part     string  //当前节点的路径片段（如 p、:lang 或 doc）
	children []*node //子节点列表
	isWild   bool    //是否是通配符节点。如果 part 以 : 或 * 开头，则为 true
}

func (n *node) String() string {
	return fmt.Sprintf("node{paatern:%s,part:%s,iswild:%t}", n.pattern, n.part, n.isWild)
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	if child != nil {
		child.insert(pattern, parts, height+1)
	}
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil

}
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if part == child.part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
