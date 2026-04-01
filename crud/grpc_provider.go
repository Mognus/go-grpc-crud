package crud

import "github.com/gofiber/fiber/v2"

// GRPCProvider is the minimal interface for the admin panel.
// Implement this for gRPC-backed providers — no DB operations needed.
// CRUDProvider satisfies this automatically, so existing DB-backed
// providers work unchanged with RegisterCRUD.
type GRPCProvider interface {
	GetModelName() string
	GetSchema() Schema

	ListHandler() fiber.Handler
	SchemaHandler() fiber.Handler
	GetHandler() fiber.Handler
	CreateHandler() fiber.Handler
	UpdateHandler() fiber.Handler
	DeleteHandler() fiber.Handler
}
