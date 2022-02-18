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
	user.Put("/:id", middleware.IsAuthenticated, UpdateUser)
	user.Delete("/:id", middleware.IsAuthenticated, DeleteUser)

	// Authentication
	auth := api.Group("/auth")
	auth.Post("/signup", Signup)
	auth.Post("/login", Login)
	auth.Get("/me", middleware.IsAuthenticated, UserMe)
	auth.Get("/refresh", middleware.IsAuthenticatedRefresh, Refresh)

	product := api.Group("/products")
	product.Post("/", middleware.IsAuthenticated, CreateProduct)
	product.Get("/", GetAllProducts)
	product.Get("/:id", GetProduct)
	product.Put("/:id", middleware.IsAuthenticated, UpdateProduct)
	product.Delete("/:id", middleware.IsAuthenticated, DeleteProduct)

	orderItem := api.Group("/orderitems")
	orderItem.Post("/", middleware.IsAuthenticated, CreateOrderItem)
	orderItem.Get("/", GetAllOrderItems)
	orderItem.Get("/:id", GetOrderItem)
	orderItem.Put("/:id", middleware.IsAuthenticated, UpdateOrderItem)
	orderItem.Get("/order/:id", GetAllOrderItemsByOrderID)
	orderItem.Delete("/:id", middleware.IsAuthenticated, DeleteOrderItem)

	order := api.Group("/orders")
	order.Post("/", middleware.IsAuthenticated, CreateOrder)
	order.Get("/", GetAllOrders)
	order.Get("/:id", GetOrderByID)
	order.Put("/:id", middleware.IsAuthenticated, UpdateOrder)
	order.Get("/user/:id", GetOrdersByUserID)
	order.Delete("/:id", middleware.IsAuthenticated, DeleteOrder)
}
