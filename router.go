// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package router

const (
	normalNode = iota
	pathParamNode
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

func (r *router) addRoute(addr string, v interface{}) {
	if addr == "" || addr[0] != '/' {
		panic("invalid address")
	}
	r.nodes.parseNode(addr[1:], v)
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
		return
	}
	i, start, parCur := 0, 0, 0
	for ; i < len(addr); i++ {
		if addr[i] == ':' {
			parCur = i
		} else if addr[i] == '/' {
			break
		}
	}

	nType := normalNode
	if parCur > 0 {
		nType = pathParamNode
	}

	var tmp *node = nil
	if nType == normalNode {
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

	if i < len(addr) {
		tmp.parseNode(addr[i+1:], value)
	} else {
		tmp.v = value
	}
}
