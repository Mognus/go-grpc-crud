package crud

import (
	"reflect"

	"gorm.io/gorm"
)

// ============================================================================
// Default CRUD Operations
// ============================================================================

// DefaultList provides a default implementation for listing items.
// Pass a ListConfig to configure eager-loaded relations and searchable fields.
// Usage: DefaultList(db, &User{}, filters, 1, 20, ListConfig{Preloads: []string{"Role"}, Searchable: []string{"name","email"}})
func DefaultList(db *gorm.DB, model any, filters map[string]string, page, limit int, config ...ListConfig) (ListResponse, error) {
	cfg := ListConfig{}
	if len(config) > 0 {
		cfg = config[0]
	}

	query := db.Model(model)

	// Apply preloads for relations
	for _, preload := range cfg.Preloads {
		query = query.Preload(preload)
	}

	// Apply global search (OR across searchable fields) and column filters
	query = ApplySearch(query, filters, cfg.Searchable)
	query = ApplyFilters(query, filters)
	query = ApplySorting(query, filters, "id DESC")

	// Count total before pagination
	total, err := CountTotal(query)
	if err != nil {
		return ListResponse{}, err
	}

	// Apply pagination and execute
	query = ApplyPagination(query, page, limit)
	items, err := ExecuteQuery(query, model)
	if err != nil {
		return ListResponse{}, err
	}

	return BuildListResponse(items, total, page, limit), nil
}

// DefaultGet provides a default implementation for getting a single item by ID.
// Optional preloads parameter allows eager loading of relations.
func DefaultGet(db *gorm.DB, model any, id string, preloads ...string) (any, error) {
	idUint, err := ParseID(id)
	if err != nil {
		return nil, err
	}

	query := db.Model(model)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	result := reflect.New(reflect.TypeOf(model).Elem()).Interface()
	if err := query.First(result, idUint).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// DefaultCreate provides a default implementation for creating a new item
func DefaultCreate(db *gorm.DB, model any, data map[string]any) (any, error) {
	result := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	// Map data to model fields
	if err := MapDataToModel(result, data); err != nil {
		return nil, err
	}

	if err := db.Create(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// DefaultUpdate provides a default implementation for updating an item
func DefaultUpdate(db *gorm.DB, model any, id string, data map[string]any) (any, error) {
	idUint, err := ParseID(id)
	if err != nil {
		return nil, err
	}

	// Get existing item
	result := reflect.New(reflect.TypeOf(model).Elem()).Interface()
	if err := db.First(result, idUint).Error; err != nil {
		return nil, err
	}

	// Update fields
	if err := MapDataToModel(result, data); err != nil {
		return nil, err
	}

	if err := db.Save(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// DefaultDelete provides a default implementation for deleting an item
func DefaultDelete(db *gorm.DB, model any, id string) error {
	idUint, err := ParseID(id)
	if err != nil {
		return err
	}

	if err := db.Delete(model, idUint).Error; err != nil {
		return err
	}

	return nil
}
