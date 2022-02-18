package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rama-kairi/fiber-api/database"
	"github.com/rama-kairi/fiber-api/models"
)

func ProductResponse(product models.Product) map[string]interface{} {
	response := map[string]interface{}{
		"id":         product.ID,
		"created_at": product.CreatedAt,
		"updated_at": product.UpdatedAt,
		"name":       product.Name,
		"price":      product.Price,
		"quantity":   product.Quantity,
	}
	return response
}

// CreateProduct creates a new product
func CreateProduct(c *fiber.Ctx) error {
	product := models.Product{}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.Database.Db.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responseProduct := ProductResponse(product)

	return c.Status(fiber.StatusCreated).JSON(responseProduct)
}

// GetAllProducts returns all products
func GetAllProducts(c *fiber.Ctx) error {
	var products []models.Product

	database.Database.Db.Find(&products)

	responseProducts := make([]map[string]interface{}, len(products))

	for i, product := range products {
		responseProducts[i] = ProductResponse(product)
	}

	return c.JSON(responseProducts)
}

// GetProduct returns a product by id
func GetProduct(c *fiber.Ctx) error {
	var product models.Product

	id := c.Params("id")

	database.Database.Db.First(&product, id)

	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found with id " + id,
		})
	}

	responseProduct := ProductResponse(product)

	return c.JSON(responseProduct)
}

// UpdateProduct updates a product by id
func UpdateProduct(c *fiber.Ctx) error {
	var product models.Product

	id := c.Params("id")

	database.Database.Db.First(&product, id)

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found with id " + id,
		})
	}

	database.Database.Db.Save(&product)

	responseProduct := ProductResponse(product)

	return c.JSON(responseProduct)
}

// DeleteProduct deletes a product by id
func DeleteProduct(c *fiber.Ctx) error {
	var product models.Product

	id := c.Params("id")

	database.Database.Db.First(&product, id)

	if product.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found with id " + id,
		})
	}

	database.Database.Db.Delete(&product)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
