package crud

import "github.com/gofiber/fiber/v2"

// GRPCProvider is the minimal interface for the admin panel.
// Implement this for gRPC-backed providers — no DB operations needed.
// CRUDProvider satisfies this automatically, so existing DB-backed
// providers work unchanged with RegisterCRUD.
type GRPCProvider interface {
	GetModelName() string
	GetSchema() Schema

	HandleList(c *fiber.Ctx) error
	HandleSchema(c *fiber.Ctx) error
	HandleGet(c *fiber.Ctx) error
	HandleCreate(c *fiber.Ctx) error
	HandleUpdate(c *fiber.Ctx) error
	HandleDelete(c *fiber.Ctx) error
}
