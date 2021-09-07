// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package utils

import "strings"

type RouteConverter struct {
	src  string
	dest string
}

func NewRouteConverter(src, dest string) *RouteConverter {
	if src[len(src)-1] == '*' {
		src = src[:len(src)-1]
	}
	if dest[len(dest)-1] == '*' {
		dest = dest[:len(dest)-1]
	}
	ret := &RouteConverter{
		src:  src,
		dest: dest,
	}
	return ret
}

func (r *RouteConverter) ConvertAddress(addr string, m map[string]string) string {
	src, dest := r.src, r.dest
	for k, v := range m {
		src = strings.Replace(src, k, v, -1)
		dest = strings.Replace(dest, k, v, -1)
	}
	return strings.Replace(addr, src, dest, 1)
}
