// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package router

type Router interface {
	// 增加路由
	// 使用/:id冒号+path参数名称标识path类型参数
	// '*' 表示通配后续所有的路径，如果需要单独指定，则需要放在通配符加入路由之前。
	//    AddRoute: "/hello/world"
	//    AddRoute: "/hello/*"
	// 注意'*'通配符之后不允许添加任何路径字符
	AddRoute(addr string, v interface{}) error

	// 查询router中是否有匹配路径的路由
	// Param: addr 用于匹配的路径
	// Param: m 用于存储PathParam路径参数的map，key为“:变量”，value为addr中实际的值
	// Return: interface{} 添加路由时传入的value
	// Return: error 发生错误时抛出
	Match(addr string, m *map[string]string) (interface{}, error)
}
