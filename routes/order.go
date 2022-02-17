package routes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rama-kairi/fiber-api/database"
	"github.com/rama-kairi/fiber-api/models"
)

func OrderItemsResponse(orderItem models.OrderItem, product models.Product) map[string]interface{} {
	return map[string]interface{}{
		"id":         orderItem.ID,
		"created_at": orderItem.CreatedAt,
		"updated_at": orderItem.UpdatedAt,
		"quantity":   orderItem.Quantity,
		"price":      orderItem.Price,
		"product":    product,
		"product_id": orderItem.ProductID,
		"order_id":   orderItem.OrderID,
	}
}

func OrderItemsAllResponse(orderItems []models.OrderItem) []map[string]interface{} {
	responseOrderItems := make([]map[string]interface{}, len(orderItems))

	for i, orderItem := range orderItems {
		responseOrderItems[i] = OrderItemsResponse(orderItem, orderItem.Product)
	}

	return responseOrderItems
}

func OrderResponse(order models.Order, OrderItems []models.OrderItem, user models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         order.ID,
		"created_at": order.CreatedAt,
		"updated_at": order.UpdatedAt,
		"quantity":   order.Quantity,
		"price":      order.Price,
		"user":       ResponseUser(user),
		"userID":     order.UserID,
		"orderItems": OrderItemsAllResponse(OrderItems),
	}
}

// CreateOrderItem - add new order item
func CreateOrderItem(c *fiber.Ctx) error {
	// Schema for order item Create
	type orderItemCreate struct {
		Quantity  int `json:"quantity"`
		ProductID int `json:"product_id"`
	}
	// Declaring the DB variable
	db := database.Database.Db

	// Converting the OrderItemCreate struct to JSON
	orderItemJson := new(orderItemCreate)

	// Parsing the body of the request into the order item variable
	if err := c.BodyParser(&orderItemJson); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Declaring the Product variable
	product := models.Product{}

	// Getting the product by the product id
	db.First(&product, orderItemJson.ProductID)

	// println(fmt.Sprintf("%T", strconv.Itoa(order_item.ProductID)))

	// Checking if the product exists
	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found with id " + strconv.Itoa(orderItemJson.ProductID),
		})
	}

	// Declaring the Order variable for Creating the OrderItem with custom Price
	orderItemInstance := models.OrderItem{
		Quantity:  orderItemJson.Quantity,
		ProductID: orderItemJson.ProductID,
		Price:     product.Price * float64(orderItemJson.Quantity),
	}

	// Creating the OrderItem
	db.Create(&orderItemInstance)

	// Returning the OrderItem
	return c.Status(fiber.StatusCreated).JSON(OrderItemsResponse(orderItemInstance, product))
}

// GetAllOrderItems - get all order items
func GetAllOrderItems(c *fiber.Ctx) error {
	var orderItems []models.OrderItem

	database.Database.Db.Find(&orderItems)

	responseOrderItems := make([]map[string]interface{}, len(orderItems))

	for i, orderItem := range orderItems {
		product := models.Product{}
		database.Database.Db.First(&product, orderItem.ProductID)
		responseOrderItems[i] = OrderItemsResponse(orderItem, product)
	}

	return c.JSON(responseOrderItems)
}

// GetOrderItem - get order item by id
func GetOrderItem(c *fiber.Ctx) error {
	var orderItem models.OrderItem

	id := c.Params("id")

	database.Database.Db.First(&orderItem, id)

	if orderItem.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order Item not found with id " + id,
		})
	}

	product := models.Product{}
	database.Database.Db.First(&product, orderItem.ProductID)

	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found with id " + strconv.Itoa(orderItem.ProductID),
		})
	}

	return c.JSON(OrderItemsResponse(orderItem, product))
}

// GetAllOrderItemsByOrderID - get all order items by order id
func GetAllOrderItemsByOrderID(c *fiber.Ctx) error {
	var orderItems []models.OrderItem

	orderID := c.Params("id")

	database.Database.Db.Where("order_id = ?", orderID).Find(&orderItems)

	responseOrderItems := make([]map[string]interface{}, len(orderItems))

	for i, orderItem := range orderItems {
		product := models.Product{}
		database.Database.Db.First(&product, orderItem.ProductID)
		responseOrderItems[i] = OrderItemsResponse(orderItem, product)
	}

	return c.JSON(responseOrderItems)
}

// UpdateOrderItem - update order item by id
func UpdateOrderItem(c *fiber.Ctx) error {
	type orderItemUpdate struct {
		Quantity  int `json:"quantity"`
		ProductID int `json:"product_id"`
	}
	db := database.Database.Db
	orderItemJson := new(orderItemUpdate)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid id" + err.Error(),
		})
	}

	if err := c.BodyParser(&orderItemJson); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	orderItem := models.OrderItem{}
	db.First(&orderItem, id)

	if orderItemJson.ProductID != 0 {
		orderItem.ProductID = orderItemJson.ProductID
	}

	if orderItemJson.Quantity != 0 {
		orderItem.Quantity = orderItemJson.Quantity
	}

	product := models.Product{}
	database.Database.Db.First(&product, orderItem.ProductID)

	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found with id " + strconv.Itoa(orderItem.ProductID),
		})
	}
	orderItem.Price = product.Price * float64(orderItem.Quantity)

	db.Save(&orderItem)

	return c.JSON(OrderItemsResponse(orderItem, product))
}

// DeleteOrderItem - Delete Order Item

func DeleteOrderItem(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid id" + err.Error(),
		})
	}

	orderItem := models.OrderItem{}
	database.Database.Db.First(&orderItem, id)

	if orderItem.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order Item not found with id " + strconv.Itoa(id),
		})
	}

	database.Database.Db.Delete(&orderItem)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

// Order Routes

// CreateOrder - Create Order
func CreateOrder(c *fiber.Ctx) error {
	type OrderCreate struct {
		OrderItemIds []int `json:"order_item_ids"`
	}

	claims := c.Locals("user")
	userID := claims.(jwt.MapClaims)["user_id"].(float64)

	db := database.Database.Db
	orderJson := new(OrderCreate)

	if err := c.BodyParser(&orderJson); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	orderItems_all := make([]models.OrderItem, len(orderJson.OrderItemIds))
	price := 0.0
	quantity := len(orderJson.OrderItemIds)

	for i, orderItemID := range orderJson.OrderItemIds {
		orderItem := models.OrderItem{}
		db.First(&orderItem, orderItemID)

		if orderItem.ID == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order Item not found with id " + strconv.Itoa(orderItemID),
			})
		}

		price += orderItem.Price

		orderItems_all[i] = orderItem
	}

	order := models.Order{
		Price:    price,
		Quantity: quantity,
		UserID:   int(userID),
	}
	db.Create(&order)

	db.Model(&order).Association("OrderItems").Append(orderItems_all)

	user := models.User{}
	database.Database.Db.First(&user, order.UserID)

	return c.Status(fiber.StatusCreated).JSON(OrderResponse(order, orderItems_all, user))
}

// GetAllOrders - Get all orders
func GetAllOrders(c *fiber.Ctx) error {
	var orders []models.Order
	db := database.Database.Db

	db.Find(&orders)

	responseOrders := make([]map[string]interface{}, len(orders))

	for i, order := range orders {
		orderItems := make([]models.OrderItem, len(order.OrderItems))
		db.Where("order_id = ?", order.ID).Find(&orderItems)

		user := models.User{}
		db.First(&user, order.UserID)

		responseOrders[i] = OrderResponse(order, orderItems, user)
	}

	return c.JSON(responseOrders)
}

// GetOrderByID - Get order by id
func GetOrderByID(c *fiber.Ctx) error {
	var order models.Order
	db := database.Database.Db

	id := c.Params("id")

	db.First(&order, id)

	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found with id " + id,
		})
	}

	orderItems := make([]models.OrderItem, len(order.OrderItems))
	db.Where("order_id = ?", order.ID).Find(&orderItems)

	user := models.User{}
	database.Database.Db.First(&user, order.UserID)

	return c.JSON(OrderResponse(order, orderItems, user))
}

// UpdateOrder - Update Order
func UpdateOrder(c *fiber.Ctx) error {
	type OrderUpdate struct {
		OrderItemIds []int `json:"order_item_ids"`
	}

	db := database.Database.Db
	orderJson := new(OrderUpdate)

	if err := c.BodyParser(&orderJson); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	order := models.Order{}
	db.First(&order, c.Params("id"))

	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found with id " + c.Params("id"),
		})
	}

	orderItems_all := make([]models.OrderItem, len(orderJson.OrderItemIds))
	price := 0.0
	quantity := len(orderJson.OrderItemIds)

	for i, orderItemID := range orderJson.OrderItemIds {
		orderItem := models.OrderItem{}
		db.First(&orderItem, orderItemID)

		price += orderItem.Price
		orderItems_all[i] = orderItem
	}

	order.Price = price
	order.Quantity = quantity

	db.Save(&order)

	db.Model(&order).Association("OrderItems").Replace(orderItems_all)

	user := models.User{}
	database.Database.Db.First(&user, order.UserID)

	return c.JSON(OrderResponse(order, orderItems_all, user))
}

// GetOrdersByUserID - Get orders by user id
func GetOrdersByUserID(c *fiber.Ctx) error {
	var orders []models.Order
	db := database.Database.Db

	db.Where("user_id = ?", c.Params("id")).Find(&orders)

	responseOrders := make([]map[string]interface{}, len(orders))

	for i, order := range orders {
		orderItems := make([]models.OrderItem, len(order.OrderItems))
		db.Where("order_id = ?", order.ID).Find(&orderItems)

		user := models.User{}
		db.First(&user, order.UserID)

		responseOrders[i] = OrderResponse(order, orderItems, user)
	}

	return c.Status(fiber.StatusOK).JSON(responseOrders)
}

// DeleteOrder - Delete Order
func DeleteOrder(c *fiber.Ctx) error {
	var order models.Order
	db := database.Database.Db

	id := c.Params("id")

	db.First(&order, id)

	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found with id " + id,
		})
	}

	db.Delete(&order)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
