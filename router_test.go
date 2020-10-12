// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package router

import "testing"

func TestRouter(t *testing.T) {
	r := newRouter()
	r.addRoute("/", "/")
	r.addRoute("/hello", "/hello")
	r.addRoute("/hello/:id", "/hello/:id")
	r.addRoute("/hello/:id/world", "/hello/:id/world")
	r.addRoute("/hello/:id/world/:name", "/hello/:id/world/:name")
	r.addRoute("/test", "/test")
	r.addRoute("/hello/:id/world", "/hello/:id/world2")

	for _, v := range r.nodes.children {
		t.Log(v.path, v.v)
		if (v.path != "hello" || v.v.(string) != "/hello") &&
			(v.path != "test" || v.v.(string) != "/test") {
			t.Fatal("not match")
		}
		for _, v := range v.children {
			t.Log(v.path, v.v)
			if v.path != ":id" || v.v.(string) != "/hello/:id" {
				t.Fatal("not match")
			}
			for _, v := range v.children {
				t.Log(v.path, v.v)
				if v.path != "world" || v.v.(string) != "/hello/:id/world2" {
					t.Fatal("not match")
				}
				for _, v := range v.children {
					t.Log(v.path, v.v)
					if v.path != ":name" || v.v.(string) != "/hello/:id/world/:name" {
						t.Fatal("not match")
					}
				}
			}
		}
	}
}
