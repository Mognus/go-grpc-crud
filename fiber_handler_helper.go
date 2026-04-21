package fiberhandler

import (
	"strconv"

	apperrors "github.com/Mognus/go-grpc-crud/errors"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ListParams struct {
	Page      int32
	Limit     int32
	Search    string
	Filters   map[string]string
	SortBy    string
	SortOrder string
}

type ListResponse struct {
	Items []any `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func ParseListParams(c *fiber.Ctx) ListParams {
	params := ListParams{
		Page:      int32(c.QueryInt("page", 1)),
		Limit:     int32(c.QueryInt("limit", 20)),
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Filters:   map[string]string{},
	}

	reserved := map[string]bool{
		"page": true, "limit": true, "search": true,
		"sort_by": true, "sort_order": true,
	}

	c.Context().QueryArgs().VisitAll(func(k, v []byte) {
		if key := string(k); !reserved[key] {
			params.Filters[key] = string(v)
		}
	})

	return params
}

func ParseID(c *fiber.Ctx) (uint64, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, apperrors.BadRequest("invalid id")
	}
	return id, nil
}

func DecodeBody[T any](c *fiber.Ctx, req *T) error {
	if m, ok := any(req).(proto.Message); ok {
		decoder := protojson.UnmarshalOptions{DiscardUnknown: true}
		if err := decoder.Unmarshal(c.Body(), m); err != nil {
			return apperrors.BadRequest("invalid body")
		}
		return nil
	}

	if err := c.BodyParser(req); err != nil {
		return apperrors.BadRequest("invalid body")
	}

	return nil
}

func WriteList(c *fiber.Ctx, items []any, total int64, page, limit int32) error {
	return c.JSON(ListResponse{
		Items: items,
		Total: total,
		Page:  int(page),
		Limit: int(limit),
	})
}

func WriteCreated(c *fiber.Ctx, item any) error {
	return c.Status(fiber.StatusCreated).JSON(item)
}
