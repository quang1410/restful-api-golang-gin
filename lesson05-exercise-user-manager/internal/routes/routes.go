package routes

import "github.com/gin-gonic/gin"

// Route là một interface (hợp đồng).
// Bất kỳ struct nào có method Register() đều được coi là một Route.
// Mục đích: giúp RegisterRoutes không cần biết cụ thể là UserRoutes hay ProductRoutes,
// chỉ cần biết "cái này có Register()" là đủ để gọi.
type Route interface {
	Register(r *gin.RouterGroup)
}

// RegisterRoutes nhận vào gin.Engine và nhiều Route (dùng ...Route để nhận nhiều tham số).
// Tạo một group chung /api/v1, rồi lần lượt gọi Register() của từng route.
// Nhờ dùng interface, sau này thêm ProductRoutes, OrderRoutes... chỉ cần truyền vào đây,
// không cần sửa gì trong hàm này.
func RegisterRoutes(r *gin.Engine, routes ...Route) {
	// Tạo route group /api/v1 — tất cả các route con đều có prefix này
	// Ví dụ: /api/v1/users, /api/v1/products,...
	apiGroup := r.Group("/api/v1")

	// Duyệt qua từng route và gọi Register() của nó,
	// truyền vào apiGroup để các route con đăng ký vào đúng prefix /api/v1
	for _, route := range routes {
		route.Register(apiGroup)
	}
}
