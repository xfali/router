// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package router

import (
	"errors"
	"strings"
)

const (
	normalNode = iota
	pathParamNode
	wildcardNode
)

type node struct {
	fullPath string
	nType    int8
	children []*node
	path     string

	v interface{}
}

type router struct {
	nodes *node
}

func newRouter() *router {
	n := &node{
		path:  "/",
		nType: normalNode,
	}
	return &router{
		nodes: n,
	}
}

// 增加路由
// 使用/:id冒号+path参数名称标识path类型参数
// /* 表示通配后续所有的路径，如果需要单独指定，则需要放在通配符加入路由之前。
//    addRoute: "/hello/world"
//    addRoute: "/hello/*"
func (r *router) addRoute(addr string, v interface{}) {
	if addr == "" || addr[0] != '/' {
		panic("invalid address")
	}
	r.nodes.parseNode(addr[1:], v)
}

func (r *router) match(addr string) (interface{}, error) {
	v, err := r.matchNode(addr)
	if v != nil {
		return v.Get(true), err
	}
	return v, err
}

func (r *router) matchAddress(addr string, m *map[string]string) (interface{}, error) {
	return r.nodes.matchString(addr, m)
}

func (r *router) matchNode(addr string) (*node, error) {
	node := parseNode(addr)
	if node == nil {
		return nil, errors.New("cannot parse address")
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
			return nil, errors.New("path not match")
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
		return nil, errors.New("other's children must be 1")
	}
	return nil, errors.New("not found")
}

// 提取符合paths路径的node
func (n *node) matchString(addr string, m *map[string]string) (interface{}, error) {
	if addr == "" || addr[0] != '/' {
		return nil, errors.New("invalid address")
	}
	paths := strings.Split(addr[1:], "/")
	p := make([]string, 0, len(paths)+1)
	p = append(p, "/")
	p = append(p, paths...)
	return n.matchPaths(p, m)
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
			return nil, errors.New("path not match")
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

	return nil, errors.New("not found")
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

func (n *node) parseNode(addr string, value interface{}) {
	if addr == "" {
		n.v = value
		return
	}
	i, start := 0, 0
	nType := normalNode
	if addr[0] == '*' {
		i = 1
		nType = wildcardNode
	} else {
		for ; i < len(addr); i++ {
			if addr[i] == ':' {
				nType = pathParamNode
			} else if addr[i] == '/' {
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

	if i < len(addr) && tmp.nType != wildcardNode {
		tmp.parseNode(addr[i+1:], value)
	} else {
		tmp.v = value
	}
}
