package crud

import "github.com/gofiber/fiber/v2"

// DefaultListHandler returns a default handler for listing items
func DefaultListHandler(provider CRUDProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get pagination parameters
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 20)

		// Get all query parameters as filters
		filters := make(map[string]string)
		c.Context().QueryArgs().VisitAll(func(key, value []byte) {
			keyStr := string(key)
			// Skip page and limit
			if keyStr != "page" && keyStr != "limit" {
				filters[keyStr] = string(value)
			}
		})

		// Get data from provider
		response, err := provider.List(filters, page, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(response)
	}
}

// DefaultSchemaHandler returns a default handler for getting schema
func DefaultSchemaHandler(provider CRUDProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schema := provider.GetSchema()
		return c.JSON(schema)
	}
}

// DefaultGetHandler returns a default handler for getting a single item
func DefaultGetHandler(provider CRUDProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		item, err := provider.Get(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(item)
	}
}

// DefaultCreateHandler returns a default handler for creating an item
func DefaultCreateHandler(provider CRUDProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse request body
		var data map[string]any
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Create item
		item, err := provider.Create(data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(item)
	}
}

// DefaultUpdateHandler returns a default handler for updating an item
func DefaultUpdateHandler(provider CRUDProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Parse request body
		var data map[string]any
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Update item
		item, err := provider.Update(id, data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(item)
	}
}

// DefaultDeleteHandler returns a default handler for deleting an item
func DefaultDeleteHandler(provider CRUDProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Delete item
		if err := provider.Delete(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
