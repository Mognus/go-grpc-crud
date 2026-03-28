package server

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type ListRequest struct {
	Page      int32
	Limit     int32
	Search    string
	Filters   map[string]string
	SortBy    string
	SortOrder string
}

type ListConfig struct {
	Preloads        []string
	Searchable      []string
	SortableColumns []string
	DefaultSort     string
}

func DefaultList[T any](db *gorm.DB, req ListRequest, cfg ListConfig) ([]T, int64, error) {
	page, limit := int(req.Page), int(req.Limit)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	var zero T
	query := db.Model(&zero)

	for _, p := range cfg.Preloads {
		query = query.Preload(p)
	}

	if req.Search != "" && len(cfg.Searchable) > 0 {
		conditions := make([]string, len(cfg.Searchable))
		values := make([]any, len(cfg.Searchable))
		for i, f := range cfg.Searchable {
			conditions[i] = f + " ILIKE ?"
			values[i] = "%" + req.Search + "%"
		}
		query = query.Where(strings.Join(conditions, " OR "), values...)
	}

	query = applyFilters(query, req.Filters)

	var total int64
	query.Count(&total)

	defaultSort := cfg.DefaultSort
	if defaultSort == "" {
		defaultSort = "id ASC"
	}
	query = applySorting(query, req.SortBy, req.SortOrder, cfg.SortableColumns, defaultSort)

	var results []T
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func DefaultGet[T any](db *gorm.DB, id uint64, preloads ...string) (T, error) {
	var result T
	if err := withPreloads(db, preloads).First(&result, id).Error; err != nil {
		return result, err
	}
	return result, nil
}

func DefaultCreate[T any](db *gorm.DB, model *T, preloads ...string) (*T, error) {
	if err := db.Create(model).Error; err != nil {
		return nil, err
	}
	if len(preloads) > 0 {
		if err := withPreloads(db, preloads).Find(model).Error; err != nil {
			return nil, err
		}
	}
	return model, nil
}

func DefaultUpdate[T any](db *gorm.DB, id uint64, updates map[string]any, preloads ...string) (*T, error) {
	var result T
	if err := db.First(&result, id).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&result).Updates(updates).Error; err != nil {
		return nil, err
	}
	if len(preloads) > 0 {
		if err := withPreloads(db, preloads).Find(&result, id).Error; err != nil {
			return nil, err
		}
	}
	return &result, nil
}

func DefaultDelete(db *gorm.DB, model any, id uint64) error {
	return db.Delete(model, id).Error
}

func withPreloads(db *gorm.DB, preloads []string) *gorm.DB {
	for _, p := range preloads {
		db = db.Preload(p)
	}
	return db
}

func applyFilters(query *gorm.DB, filters map[string]string) *gorm.DB {
	for key, value := range filters {
		if value == "" {
			continue
		}
		if idx := strings.Index(key, "__"); idx != -1 {
			field, op := key[:idx], key[idx+2:]
			switch op {
			case "contains":
				query = query.Where(fmt.Sprintf("%s ILIKE ?", field), "%"+value+"%")
			case "startswith":
				query = query.Where(fmt.Sprintf("%s ILIKE ?", field), value+"%")
			case "endswith":
				query = query.Where(fmt.Sprintf("%s ILIKE ?", field), "%"+value)
			case "gte":
				query = query.Where(fmt.Sprintf("%s >= ?", field), value)
			case "lte":
				query = query.Where(fmt.Sprintf("%s <= ?", field), value)
			case "gt":
				query = query.Where(fmt.Sprintf("%s > ?", field), value)
			case "lt":
				query = query.Where(fmt.Sprintf("%s < ?", field), value)
			case "ne":
				query = query.Where(fmt.Sprintf("%s != ?", field), value)
			}
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	return query
}

func applySorting(query *gorm.DB, sortBy, sortOrder string, allowed []string, defaultOrder string) *gorm.DB {
	if sortBy == "" {
		return query.Order(defaultOrder)
	}
	allowedSet := make(map[string]bool, len(allowed))
	for _, c := range allowed {
		allowedSet[c] = true
	}
	if len(allowed) > 0 && !allowedSet[sortBy] {
		return query.Order(defaultOrder)
	}
	order := "ASC"
	if strings.ToLower(sortOrder) == "desc" {
		order = "DESC"
	}
	return query.Order(sortBy + " " + order)
}
