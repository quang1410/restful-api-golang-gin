package routes

import (
	"galvin/lession05-exercise-user-management/internal/handler"

	"github.com/gin-gonic/gin"
)

// UserRoutes chịu trách nhiệm đăng ký tất cả các route liên quan đến User.
// Nó giữ một tham chiếu đến UserHandler để gọi các hàm xử lý tương ứng.
type UserRoutes struct {
	handler *handler.UserHandler
}

// NewUserRoutes là constructor — nhận vào UserHandler và trả về *UserRoutes.
// Được gọi trong main.go khi khởi tạo ứng dụng.
func NewUserRoutes(handler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}

// Register là method thực thi interface Route.
// Nhận vào một RouterGroup (ở đây là /api/v1) và đăng ký các route con vào đó.
// Kết quả: các route bên dưới sẽ có đường dẫn đầy đủ là /api/v1/users/...
func (ur *UserRoutes) Register(r *gin.RouterGroup) {
	// Tạo group con /users bên trong /api/v1 → full path: /api/v1/users
	users := r.Group("/users")
	{
		// GET /api/v1/users → lấy danh sách tất cả users
		users.GET("/", ur.handler.GetAllUsers)

		// GET /api/v1/users/:id → lấy user theo ID
		// users.GET("/:id", ur.handler.GetUserByID)

		// POST /api/v1/users → tạo user mới
		// users.POST("/", ur.handler.CreateUser)

		// PUT /api/v1/users/:id → cập nhật user
		// users.PUT("/:id", ur.handler.UpdateUser)

		// DELETE /api/v1/users/:id → xóa user
		// users.DELETE("/:id", ur.handler.DeleteUser)
	}
}
