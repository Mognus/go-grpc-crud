package proxy

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ListResponse struct {
	Items []any `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

type ListCallFunc func(ctx context.Context, page, limit int32, search string, filters map[string]string, sortBy, sortOrder string) ([]any, int64, error)

func DefaultListProxy(call ListCallFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := int32(c.QueryInt("page", 1))
		limit := int32(c.QueryInt("limit", 20))

		reserved := map[string]bool{
			"page": true, "limit": true, "search": true,
			"sort_by": true, "sort_order": true,
		}
		filters := map[string]string{}
		c.Context().QueryArgs().VisitAll(func(k, v []byte) {
			if key := string(k); !reserved[key] {
				filters[key] = string(v)
			}
		})

		items, total, err := call(context.Background(), page, limit, c.Query("search"), filters, c.Query("sort_by"), c.Query("sort_order"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(ListResponse{Items: items, Total: total, Page: int(page), Limit: int(limit)})
	}
}

func DefaultGetProxy(call func(ctx context.Context, id uint64) (any, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		item, err := call(context.Background(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(item)
	}
}

func DefaultCreateProxy(call func(ctx context.Context, data map[string]any) (any, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]any
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		item, err := call(context.Background(), data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(item)
	}
}

func DefaultUpdateProxy(call func(ctx context.Context, id uint64, data map[string]any) (any, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		var data map[string]any
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		item, err := call(context.Background(), id, data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(item)
	}
}

func DefaultDeleteProxy(call func(ctx context.Context, id uint64) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		if err := call(context.Background(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
