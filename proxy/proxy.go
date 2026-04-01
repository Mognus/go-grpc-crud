package proxy

import (
	"context"
	"strconv"

	apperrors "github.com/Mognus/go-grpc-crud/errors"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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

		items, total, err := call(c.UserContext(), page, limit, c.Query("search"), filters, c.Query("sort_by"), c.Query("sort_order"))
		if err != nil {
			return apperrors.GrpcToHTTP(err)
		}

		return c.JSON(ListResponse{Items: items, Total: total, Page: int(page), Limit: int(limit)})
	}
}

func DefaultGetProxy(call func(ctx context.Context, id uint64) (any, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return apperrors.BadRequest("invalid id")
		}
		item, err := call(c.UserContext(), id)
		if err != nil {
			return apperrors.GrpcToHTTP(err)
		}
		return c.JSON(item)
	}
}

func parseBody[T any](c *fiber.Ctx, req *T) error {
	if m, ok := any(req).(proto.Message); ok {
		return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(c.Body(), m)
	}
	return c.BodyParser(req)
}

func DefaultCreateProxy[T any](call func(ctx context.Context, req *T) (any, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T
		if err := parseBody(c, &req); err != nil {
			return apperrors.BadRequest("invalid body")
		}
		item, err := call(c.UserContext(), &req)
		if err != nil {
			return apperrors.GrpcToHTTP(err)
		}
		return c.Status(fiber.StatusCreated).JSON(item)
	}
}

func DefaultUpdateProxy[T any](call func(ctx context.Context, id uint64, req *T) (any, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return apperrors.BadRequest("invalid id")
		}
		var req T
		if err := parseBody(c, &req); err != nil {
			return apperrors.BadRequest("invalid body")
		}
		item, err := call(c.UserContext(), id, &req)
		if err != nil {
			return apperrors.GrpcToHTTP(err)
		}
		return c.JSON(item)
	}
}

func DefaultDeleteProxy(call func(ctx context.Context, id uint64) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil {
			return apperrors.BadRequest("invalid id")
		}
		if err := call(c.UserContext(), id); err != nil {
			return apperrors.GrpcToHTTP(err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
