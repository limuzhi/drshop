package service

var (
	routerNoCheckRole = make([]func(server *AdminService), 0)
	routerCheckRole   = make([]func(server *AdminService), 0)
)

// 无需认证的路由示例
func noCheckRoleRouter(server *AdminService) {
	// 可根据业务需求来设置接口版本
	for _, f := range routerNoCheckRole {
		f(server)
	}
}

// 需要认证的路由示例
func checkRoleRouter(server *AdminService) {
	// 可根据业务需求来设置接口版本
	for _, f := range routerCheckRole {
		f(server)
	}
}
