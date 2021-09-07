// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package router

import "testing"

func TestRouter(t *testing.T) {
	r := New()
	r.AddRoute("/", "/")
	r.AddRoute("/hello", "/hello")
	r.AddRoute("/hello/:id", "/hello/:id")
	r.AddRoute("/hello/:id/world", "/hello/:id/world")
	r.AddRoute("/hello/:id/world/:name", "/hello/:id/world/:name")
	r.AddRoute("/test", "/test")
	r.AddRoute("/hello/:id/world", "/hello/:id/world2")
	r.AddRoute("/hello/:id/world/:name/*", "/hello/:id/world/:name/*")
	r.AddRoute("/hello/:id/world/:name/test", "/hello/:id/world/:name/test")
	//在/hello/:id/world/:name/*不允许添加其他的路径
	err := r.AddRoute("/hello/:id/world/:name/*/test", "/hello/:id/world/:name/*/test")
	if err == nil {
		t.Fatal(err)
	} else {
		t.Log(err)
	}

	t.Log(r.nodes.path, r.nodes.v)
	if r.nodes.v.(string) != "/" {
		t.Fatal("not match")
	}
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
					for _, v := range v.children {
						t.Log(v.path, v.v)
						if (v.path != "*" || v.v.(string) != "/hello/:id/world/:name/*") &&
							(v.path != "test" || v.v.(string) != "/hello/:id/world/:name/test") {
							t.Fatal("not match")
						}
						//cannot be here
						for _, v := range v.children {
							t.Log(v.path, v.v)
							if v.path == "test" || v.v.(string) == "/hello/:id/world/:name/*/test" {
								t.Fatal("not match")
							}
						}
					}
				}
			}
		}
	}
}

func TestMatch(t *testing.T) {
	r := New()
	r.AddRoute("/", "/")
	r.AddRoute("/hello", "/hello")
	r.AddRoute("/hello/:id", "/hello/:id")
	r.AddRoute("/hello/:id/world", "/hello/:id/world")
	r.AddRoute("/hello/:id/world/:name", "/hello/:id/world/:name")
	r.AddRoute("/test", "/test")
	r.AddRoute("/hello/:id/world/:name/*", "/hello/:id/world/:name/*")
	r.AddRoute("/hello/:id/world/:name/test", "/hello/:id/world/:name/test")

	v, err := r.Find("/hello/12/world")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(v.(string))
		if v.(string) != "/hello/:id/world" {
			t.Fatal("not match")
		}
	}

	v, err = r.Find("/hello/12/xx")
	if err == nil {
		t.Fatal(err)
	} else {
		t.Log(v, err)
	}

	v, err = r.Find("/hello/12/world/user/wwxx")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(v.(string))
		if v.(string) != "/hello/:id/world/:name/*" {
			t.Fatal("not match")
		}
	}

	v, err = r.Find("/hello/12/world/user/wwxx/dasdas")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(v.(string))
		if v.(string) != "/hello/:id/world/:name/*" {
			t.Fatal("not match")
		}
	}
}

func TestPathMap(t *testing.T) {
	r := New()
	r.AddRoute("/", "/")
	r.AddRoute("/hello", "/hello")
	r.AddRoute("/hello/:id", "/hello/:id")
	r.AddRoute("/hello/:id/world", "/hello/:id/world")
	r.AddRoute("/hello/:id/world/:name", "/hello/:id/world/:name")
	r.AddRoute("/test", "/test")
	r.AddRoute("/hello/:id/world/:name/*", "/hello/:id/world/:name/*")
	r.AddRoute("/hello/:id/world/:name/test", "/hello/:id/world/:name/test")

	v, err := r.matchNode("/hello/12/world/user/wwxx")
	if err != nil {
		t.Fatal(err)
	}
	node := parseNode("/hello/12/world/user/wwxx")
	ret := map[string]string{}
	v.pathMap(node, &ret)
	for k, v := range ret {
		t.Log(k, v)
	}
	if ret[":id"] != "12" || ret[":name"] != "user" {
		t.Fatal("not match")
	}
}

func TestMatchString(t *testing.T) {
	r := New()
	r.AddRoute("/", "/")
	r.AddRoute("/hello", "/hello")
	r.AddRoute("/hello/:id", "/hello/:id")
	r.AddRoute("/hello/:id/world", "/hello/:id/world")
	r.AddRoute("/hello/:id/world/:name", "/hello/:id/world/:name")
	r.AddRoute("/test", "/test")
	r.AddRoute("/hello/:id/world/:name/*", "/hello/:id/world/:name/*")
	r.AddRoute("/hello/:id/world/:name/test", "/hello/:id/world/:name/test")

	v, err := r.Match("/", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if v.(string) != "/" {
		t.Fatal("not match")
	}

	v, err = r.Match("/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if v.(string) != "/hello" {
		t.Fatal("not match")
	}

	ret := map[string]string{}
	v, err = r.Match("/hello/12", &ret)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if v.(string) != "/hello/:id" {
		t.Fatal("not match")
	}
	for k, v := range ret {
		t.Log(k, v)
	}
	if ret[":id"] != "12" {
		t.Fatal("not match")
	}

	ret = map[string]string{}
	v, err = r.Match("/hello/12/world/user/wwxx?xx", &ret)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	if v.(string) != "/hello/:id/world/:name/*" {
		t.Fatal("not match")
	}
	for k, v := range ret {
		t.Log(k, v)
	}
	if ret[":id"] != "12" || ret[":name"] != "user" {
		t.Fatal("not match")
	}

	ret = map[string]string{}
	v, err = r.Match("/hello/12/world/user/test", &ret)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	//注意：通配符在前，则覆盖其他的
	if v.(string) != "/hello/:id/world/:name/*" {
		t.Fatal("not match")
	}
	for k, v := range ret {
		t.Log(k, v)
	}
	if ret[":id"] != "12" || ret[":name"] != "user" {
		t.Fatal("not match")
	}

	ret = map[string]string{}
	v, err = r.Match("/hello/12/xxx/user/test", &ret)
	if err == nil {
		t.Fatal(err)
	} else {
		t.Log(err)
	}
}

func TestUse(t *testing.T) {
	r := New()
	r.AddRoute("/host", "/api/v1/test/host")
	r.AddRoute("/host/:id", "/api/v1/test/host/:id")
	r.AddRoute("/all/pass/*", "/api/v1/*")
}
