package crud

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ============================================================================
// Query Helper Functions
// ============================================================================

// ApplyFilters applies filter conditions to a GORM query
// Supports operators: contains, startswith, endswith, gte, lte, gt, lt, ne, isnull
// Usage: ?name__contains=john&age__gte=18
func ApplyFilters(query *gorm.DB, filters map[string]string) *gorm.DB {
	// Reserved keys that are not filters
	reserved := map[string]bool{
		"sort_by":    true,
		"sort_order": true,
		"page":       true,
		"limit":      true,
		"search":     true, // handled separately by ApplySearch
	}

	for key, value := range filters {
		if value == "" || reserved[key] {
			continue
		}

		// Handle special filter operators
		// e.g., "age__gte=18" -> WHERE age >= 18
		// e.g., "name__contains=john" -> WHERE name LIKE '%john%'
		if strings.Contains(key, "__") {
			parts := strings.Split(key, "__")
			field := parts[0]
			operator := parts[1]

			switch operator {
			case "gte":
				query = query.Where(fmt.Sprintf("%s >= ?", field), value)
			case "lte":
				query = query.Where(fmt.Sprintf("%s <= ?", field), value)
			case "gt":
				query = query.Where(fmt.Sprintf("%s > ?", field), value)
			case "lt":
				query = query.Where(fmt.Sprintf("%s < ?", field), value)
			case "contains":
				query = query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
			case "startswith":
				query = query.Where(fmt.Sprintf("%s LIKE ?", field), value+"%")
			case "endswith":
				query = query.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value)
			case "ne":
				query = query.Where(fmt.Sprintf("%s != ?", field), value)
			case "isnull":
				if value == "true" {
					query = query.Where(fmt.Sprintf("%s IS NULL", field))
				} else {
					query = query.Where(fmt.Sprintf("%s IS NOT NULL", field))
				}
			}
		} else {
			// Exact match
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	return query
}

// ApplySearch applies a global search across multiple fields using OR logic.
// Usage: ?search=john with searchable=["name","email"]
// → WHERE name LIKE '%john%' OR email LIKE '%john%'
func ApplySearch(query *gorm.DB, filters map[string]string, searchable []string) *gorm.DB {
	term, ok := filters["search"]
	if !ok || term == "" || len(searchable) == 0 {
		return query
	}

	conditions := make([]string, len(searchable))
	values := make([]any, len(searchable))
	for i, field := range searchable {
		conditions[i] = fmt.Sprintf("%s LIKE ?", field)
		values[i] = "%" + term + "%"
	}

	return query.Where(strings.Join(conditions, " OR "), values...)
}

// ApplySorting applies sorting to a GORM query
// Usage: ?sort_by=created_at&sort_order=desc
func ApplySorting(query *gorm.DB, filters map[string]string, defaultSort string) *gorm.DB {
	if sortBy, ok := filters["sort_by"]; ok && sortBy != "" {
		order := "ASC"
		if sortOrder, ok := filters["sort_order"]; ok && strings.ToUpper(sortOrder) == "DESC" {
			order = "DESC"
		}
		return query.Order(fmt.Sprintf("%s %s", sortBy, order))
	}

	// Default sort
	if defaultSort != "" {
		return query.Order(defaultSort)
	}
	return query.Order("id DESC")
}

// ApplyPagination applies offset and limit to a GORM query
func ApplyPagination(query *gorm.DB, page, limit int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	return query.Offset(offset).Limit(limit)
}

// CountTotal counts total records for a query (call before ApplyPagination)
func CountTotal(query *gorm.DB) (int64, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// ExecuteQuery executes the query and returns results as []any
func ExecuteQuery(query *gorm.DB, model any) ([]any, error) {
	// Create a slice to hold results
	modelType := reflect.TypeOf(model).Elem()
	results := reflect.New(reflect.SliceOf(modelType)).Interface()

	// Execute query
	if err := query.Find(results).Error; err != nil {
		return nil, err
	}

	// Convert results to []any
	resultValue := reflect.ValueOf(results).Elem()
	items := make([]any, resultValue.Len())
	for i := 0; i < resultValue.Len(); i++ {
		items[i] = resultValue.Index(i).Interface()
	}

	return items, nil
}

// BuildListResponse builds a ListResponse from items and pagination info
func BuildListResponse(items []any, total int64, page, limit int) ListResponse {
	return ListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}
}

// ============================================================================
// Utility Functions
// ============================================================================

// ParseID parses a string ID to uint64
func ParseID(id string) (uint64, error) {
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}
	return idUint, nil
}

// MapDataToModel maps a map[string]any to a model struct using reflection
func MapDataToModel(model any, data map[string]any) error {
	resultValue := reflect.ValueOf(model).Elem()

	for key, value := range data {
		// Skip ID field
		if key == "id" {
			continue
		}

		// Convert snake_case to PascalCase for struct field lookup
		fieldName := snakeToPascal(key)
		field := resultValue.FieldByName(fieldName)

		if field.IsValid() && field.CanSet() && value != nil {
			if err := setFieldValue(field, value); err != nil {
				// Skip fields that can't be set, don't fail completely
				continue
			}
		}
	}

	return nil
}

// setFieldValue sets a reflect.Value field to the given value with type conversion
func setFieldValue(field reflect.Value, value any) error {
	fieldValue := reflect.ValueOf(value)

	// Direct conversion if types are compatible
	if fieldValue.Type().ConvertibleTo(field.Type()) {
		field.Set(fieldValue.Convert(field.Type()))
		return nil
	}

	// Handle special cases
	switch field.Kind() {
	case reflect.Bool:
		if str, ok := value.(string); ok {
			field.SetBool(str == "true" || str == "1")
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if str, ok := value.(string); ok {
			if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				field.SetInt(i)
				return nil
			}
		}
		if f, ok := value.(float64); ok {
			field.SetInt(int64(f))
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if str, ok := value.(string); ok {
			if i, err := strconv.ParseUint(str, 10, 64); err == nil {
				field.SetUint(i)
				return nil
			}
		}
		if f, ok := value.(float64); ok {
			field.SetUint(uint64(f))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if str, ok := value.(string); ok {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				field.SetFloat(f)
				return nil
			}
		}
	case reflect.String:
		field.SetString(fmt.Sprintf("%v", value))
		return nil
	case reflect.Pointer:
		// Handle pointer fields (e.g. *uint for nullable foreign keys).
		// Create a new value of the pointed-to type, set it recursively, then assign the pointer.
		elemType := field.Type().Elem()
		newVal := reflect.New(elemType).Elem()
		if err := setFieldValue(newVal, value); err != nil {
			return err
		}
		ptr := reflect.New(elemType)
		ptr.Elem().Set(newVal)
		field.Set(ptr)
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", value, field.Type())
}

// GetRelatedOptions loads all records of a model and maps them to SelectOptions.
// valueField and labelField are the struct field names (PascalCase), e.g. "ID", "Name".
// Falls back to an empty string label if the label field is empty.
func GetRelatedOptions(db *gorm.DB, model any, valueField, labelField string) []SelectOption {
	modelType := reflect.TypeOf(model).Elem()
	results := reflect.New(reflect.SliceOf(modelType)).Interface()

	if err := db.Find(results).Error; err != nil {
		return nil
	}

	slice := reflect.ValueOf(results).Elem()
	options := make([]SelectOption, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		row := slice.Index(i)
		value := row.FieldByName(valueField).Interface()
		label := fmt.Sprintf("%v", row.FieldByName(labelField).Interface())
		options[i] = SelectOption{Value: value, Label: label}
	}

	return options
}

// snakeToPascal converts snake_case to PascalCase.
// "id" segments are uppercased to "ID" to match Go naming conventions,
// so that e.g. "model_asset_id" maps to "ModelAssetID" instead of "ModelAssetId".
func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			if strings.ToLower(part) == "id" {
				parts[i] = "ID"
			} else {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
	}
	return strings.Join(parts, "")
}
