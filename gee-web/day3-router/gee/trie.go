package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 根据提供的路由部分从当前节点的子节点中找出所有匹配的节点。
//
// 参数:
//
//	part - 当前正在处理的URL路由部分。
//
// 返回:
//
//	返回一个包含所有匹配的子节点的切片。如果没有匹配的节点，返回一个空切片。
//
// 函数逻辑:
//  1. 初始化一个空的节点切片用于存储匹配的子节点。
//  2. 遍历当前节点的所有子节点。
//  3. 对每个子节点，检查其part属性是否与给定的part相匹配，或者该节点是否被标记为通配（isWild）。
//  4. 如果一个子节点满足匹配条件，将其添加到结果切片中。
//  5. 返回包含所有匹配子节点的切片。
func (n *node) matchChildren(part string) []*node {
	// 初始化一个空切片用于存储匹配的节点
	nodes := make([]*node, 0)
	for _, child := range n.children {
		// 检查当前子节点的部分是否与提供的部分匹配，或该子节点是否为通配节点
		if child.part == part || child.isWild {
			// 如果匹配，添加到nodes切片中
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 这个函数将遍历或扩展现有的路由树，以确保完整的路由模式被正确添加。
// 如果到达了parts数组的末尾（height == len(parts)），则将当前节点的pattern属性设置为完整的路由模式。
// 如果在当前深度找不到匹配的子节点，则会创建一个新的子节点，并继续递归处理。
func (n *node) insert(pattern string, parts []string, height int) {
	// 检查是否达到了parts数组的末端
	if len(parts) == height {
		// 设置当前节点的pattern属性为完整路由模式
		n.pattern = pattern
		return
	}
	// 获取当前部分的路由
	part := parts[height]
	// 尝试匹配当前节点的子节点中是否存在与当前部分匹配的节点
	child := n.matchChild(part)
	// 如果没有找到匹配的子节点，创建一个新的子节点
	if child == nil {
		// 新节点的isWild属性根据part是否包含':'或'*'来设置
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 将新创建的节点添加到当前节点的children列表中
		n.children = append(n.children, child)
	}
	// 对找到或新创建的子节点进行递归处理，处理下一个路由部分
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// 检查是否达到parts数组的末端或遇到通配符节点"*"
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 如果当前节点的pattern为空，说明不是有效的终止节点，返回nil
		if n.pattern == "" {
			return nil
		}
		// 返回当前节点，因为它匹配了全部的路由片段
		return n
	}

	// 获取当前处理的路由部分
	part := parts[height]
	// 使用matchChildren而不是matchChild，假设它返回与当前部分匹配的所有子节点
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
