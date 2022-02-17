package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rama-kairi/fiber-api/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api")

	// Users
	user := api.Group("/users")
	user.Get("/", GetAllUsers)
	user.Get("/:id", GetUser)
	user.Put("/:id", middleware.Auth, UpdateUser)
	user.Delete("/:id", middleware.Auth, DeleteUser)

	// Authentication
	auth := api.Group("/auth")
	auth.Post("/signup", Signup)
	auth.Post("/login", Login)

	product := api.Group("/products")
	product.Post("/", middleware.Auth, CreateProduct)
	product.Get("/", GetAllProducts)
	product.Get("/:id", GetProduct)
	product.Put("/:id", middleware.Auth, UpdateProduct)
	product.Delete("/:id", middleware.Auth, DeleteProduct)

	orderItem := api.Group("/orderitems")
	orderItem.Post("/", middleware.Auth, CreateOrderItem)
	orderItem.Get("/", GetAllOrderItems)
	orderItem.Get("/:id", GetOrderItem)
	orderItem.Put("/:id", middleware.Auth, UpdateOrderItem)
	orderItem.Get("/order/:id", GetAllOrderItemsByOrderID)
	orderItem.Delete("/:id", middleware.Auth, DeleteOrderItem)

	order := api.Group("/orders")
	order.Post("/", middleware.Auth, CreateOrder)
	order.Get("/", GetAllOrders)
	order.Get("/:id", GetOrderByID)
	order.Put("/:id", middleware.Auth, UpdateOrder)
	order.Get("/user/:id", GetOrdersByUserID)
	order.Delete("/:id", middleware.Auth, DeleteOrder)
}
