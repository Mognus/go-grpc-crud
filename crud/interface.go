package crud

import "github.com/gofiber/fiber/v2"

// CRUDProvider is the interface that all modules must implement
// to be manageable through the admin panel
type CRUDProvider interface {
	// Model Info
	GetModelName() string // e.g. "users", "todos"
	GetSchema() Schema    // Describes fields for frontend

	// CRUD Operations
	List(filters map[string]string, page, limit int) (ListResponse, error)
	Get(id string) (any, error)
	Create(data map[string]any) (any, error)
	Update(id string, data map[string]any) (any, error)
	Delete(id string) error

	// HTTP Handlers
	ListHandler() fiber.Handler
	SchemaHandler() fiber.Handler
	GetHandler() fiber.Handler
	CreateHandler() fiber.Handler
	UpdateHandler() fiber.Handler
	DeleteHandler() fiber.Handler
}

// Schema describes the structure of a model for the frontend
type Schema struct {
	Name        string   `json:"name"`        // "users"
	DisplayName string   `json:"displayName"` // "Users"
	Fields      []Field  `json:"fields"`      // List of fields
	Searchable  []string `json:"searchable"`  // Fields that are searchable
}

// Field describes a single field in the model
type Field struct {
	Name       string   `json:"name"`       // "email"
	Type       string   `json:"type"`       // "string", "number", "boolean", "date", "enum", "relation"
	Label      string   `json:"label"`      // "Email Address"
	Required     bool `json:"required"`               // Is this field required?
	Readonly     bool `json:"readonly"`               // Always skip in forms
	TableHidden  bool `json:"tableHidden,omitempty"`  // Hide from table view
	EditHidden   bool `json:"editHidden,omitempty"`   // Hide in edit form
	CreateHidden bool `json:"createHidden,omitempty"` // Hide in create form
	EnumValues []string `json:"enumValues,omitempty"` // Possible values for enum fields
	// Relation field options (for type "relation")
	Options   []SelectOption `json:"options,omitempty"`   // Available options for relation select
	UploadURL string         `json:"uploadUrl,omitempty"` // Upload endpoint for type "file"
	Accept    string         `json:"accept,omitempty"`    // File accept filter e.g. ".glb", "image/*"
}

// SelectOption represents an option in a relation select field
type SelectOption struct {
	Value any    `json:"value"` // The ID value (usually uint)
	Label string `json:"label"` // Display text
}

// ListConfig holds optional configuration for DefaultList.
// Replaces the old variadic preloads parameter to allow passing
// both preloads and searchable fields in a single, extensible struct.
type ListConfig struct {
	Preloads   []string // GORM relations to eager-load, e.g. "ModelAsset", "Image"
	Searchable []string // Fields to search across when ?search= is provided (OR logic)
}

// ListResponse is the response format for list operations
type ListResponse struct {
	Items []any `json:"items"` // The data items
	Total int64 `json:"total"` // Total count of items
	Page  int   `json:"page"`  // Current page
	Limit int   `json:"limit"` // Items per page
}
