// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package router

import (
	"errors"
)

const (
	normalNode = iota
	pathParamNode
	wildcardNode
)

var (
	InvalidAddressError = errors.New("Address is invalid. ")
	ParseAddressError   = errors.New("Cannot parse address. ")
	NotMatchError       = errors.New("Not match. ")
)

type node struct {
	fullPath string
	nType    int8
	children []*node
	path     string

	v interface{}
}

type defaultRouter struct {
	nodes *node
}

func New() *defaultRouter {
	n := &node{
		path:  "/",
		nType: normalNode,
	}
	return &defaultRouter{
		nodes: n,
	}
}

// 增加路由
// 使用/:id冒号+path参数名称标识path类型参数
// '*' 表示通配后续所有的路径，如果需要单独指定，则需要放在通配符加入路由之前。
//    AddRoute: "/hello/world"
//    AddRoute: "/hello/*"
// 注意'*'通配符之后不允许添加任何路径字符
func (r *defaultRouter) AddRoute(addr string, v interface{}) error {
	if addr == "" || addr[0] != '/' {
		panic(InvalidAddressError)
	}
	return r.nodes.parseNode(addr[1:], v)
}

func (r *defaultRouter) Find(addr string) (interface{}, error) {
	v, err := r.matchNode(addr)
	if v != nil {
		return v.Get(true), err
	}
	return v, err
}

// 查询router中是否有匹配路径的路由
// Param: addr 用于匹配的路径
// Param: m 用于存储PathParam路径参数的map，key为“:变量”，value为addr中实际的值
// Return: interface{} 添加路由时传入的value
// Return: error 发生错误时抛出
func (r *defaultRouter) Match(addr string, m *map[string]string) (interface{}, error) {
	return r.nodes.matchString(addr, m)
}

func (r *defaultRouter) matchNode(addr string) (*node, error) {
	node := parseNode(addr)
	if node == nil {
		return nil, ParseAddressError
	}
	return r.nodes.match(node)
}

func parseNode(addr string) *node {
	if addr == "" || addr[0] != '/' {
		return nil
	}
	n := &node{
		fullPath: addr,
		path:     "/",
		nType:    normalNode,
	}
	n.parseNode(addr[1:], nil)
	return n
}

func (n *node) clone(all bool) *node {
	if !all {
		ret := &node{
			fullPath: n.fullPath,
			nType:    n.nType,
			path:     n.path,
			v:        n.v,
		}
		return ret
	} else {
		ret := &node{
			fullPath: n.fullPath,
			nType:    n.nType,
			path:     n.path,
			children: make([]*node, 0, len(n.children)),
			v:        n.v,
		}
		for _, v := range n.children {
			tmp := v.clone(all)
			ret.children = append(ret.children, tmp)
		}
		return ret
	}
}

func (n *node) Get(last bool) interface{} {
	if !last {
		return n.v
	} else {
		if len(n.children) > 0 {
			return n.children[0].Get(last)
		} else {
			return n.v
		}
	}
}

// 比较node，key为原node的pathParam，value为other的实际值
// 注意n和other如果存在children必须为1个
func (n *node) pathMap(other *node, m *map[string]string) {
	if n.nType == wildcardNode {
		return
	}

	if n.nType == pathParamNode {
		(*m)[n.path] = other.path
	}

	if len(n.children) == 1 && len(other.children) == 1 {
		n.children[0].pathMap(other.children[0], m)
	}
}

// 提取符合other路径的node
// 注意other如果存在children必须为1个
func (n *node) match(other *node) (*node, error) {
	if n.nType == wildcardNode {
		return n.clone(false), nil
	}

	if n.nType == normalNode && other.nType == normalNode {
		if n.path != other.path {
			return nil, NotMatchError
		}
	}

	if len(other.children) == 0 {
		return n.clone(false), nil
	} else if len(other.children) == 1 {
		for _, v := range n.children {
			if ret, err := v.match(other.children[0]); err == nil {
				retNode := n.clone(false)
				retNode.children = append(retNode.children, ret)
				return retNode, nil
			}
		}
	} else {
		// other's children must be 1
		return nil, NotMatchError
	}
	return nil, NotMatchError
}

// 提取符合paths路径的node
func (n *node) matchString(addr string, m *map[string]string) (interface{}, error) {
	if addr == "" || addr[0] != '/' {
		return nil, InvalidAddressError
	}
	if addr == "/" {
		return n.v, nil
	}
	paths := []string{"/"}
	paths = split(addr, paths)
	return n.matchPaths(paths, m)
}

// 提取符合paths路径的node
func (n *node) matchPaths(paths []string, m *map[string]string) (interface{}, error) {
	if len(paths) == 0 || paths[0] == "" {
		return n.v, nil
	}
	if n.nType == wildcardNode {
		return n.v, nil
	}

	if n.nType == normalNode {
		if n.path != paths[0] {
			return nil, NotMatchError
		}
	}

	if n.nType == pathParamNode && m != nil {
		(*m)[n.path] = paths[0]
	}

	if len(paths) > 1 {
		for _, v := range n.children {
			if ret, err := v.matchPaths(paths[1:], m); err == nil {
				return ret, nil
			}
		}
	} else {
		return n.v, nil
	}

	return nil, NotMatchError
}

func (n *node) equal(other *node) bool {
	if n.nType != other.nType {
		return false
	}
	if n.nType != pathParamNode && (n.path != other.path) {
		return false
	}
	if len(n.children) != len(other.children) {
		return false
	}
	for i := range n.children {
		if !n.children[i].equal(other.children[i]) {
			return false
		}
	}
	return true
}

func (n *node) parseNode(addr string, value interface{}) error {
	if addr == "" {
		n.v = value
		return nil
	}
	finished := false
	i, start := 0, 0
	nType := normalNode
	if addr[0] == '*' {
		if len(addr) > 1 {
			// Cannot have other characters after '*'
			return ParseAddressError
		}
		i = 1
		nType = wildcardNode
	} else {
		for ; i < len(addr); i++ {
			if addr[i] == ':' {
				nType = pathParamNode
			} else if addr[i] == '/' {
				break
			} else if addr[i] == '?' {
				finished = true
				break
			}
		}
	}

	var tmp *node = nil
	if nType == wildcardNode {
		for _, v := range n.children {
			if v.nType == wildcardNode {
				tmp = v
				break
			}
		}
	} else if nType == normalNode {
		path := addr[start:i]
		for _, v := range n.children {
			if v.path == path {
				tmp = v
				break
			}
		}
	} else if nType == pathParamNode {
		for _, v := range n.children {
			if v.nType == pathParamNode {
				tmp = v
				break
			}
		}
	}

	if tmp == nil {
		tmp = &node{
			path:     addr[start:i],
			fullPath: addr,
			nType:    int8(nType),
		}
		n.children = append(n.children, tmp)
	}

	if !finished && i < len(addr) && tmp.nType != wildcardNode {
		return tmp.parseNode(addr[i+1:], value)
	} else {
		tmp.v = value
	}
	return nil
}

func split(addr string, origin []string) []string {
	finished := false
	start, i := 0, 0
	for ; i < len(addr); i++ {
		if addr[i] == '/' || addr[i] == '?' {
			if i == start {
				start++
				continue
			}
			origin = append(origin, addr[start:i])
			start = i + 1
			if addr[i] == '?' {
				finished = true
				break
			}
		}
	}
	if !finished && start < i {
		origin = append(origin, addr[start:i])
	}
	return origin
}
