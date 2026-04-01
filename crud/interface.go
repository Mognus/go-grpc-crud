package crud

type Schema struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Fields      []Field  `json:"fields"`
	Searchable  []string `json:"searchable"`
}

type Field struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Label        string         `json:"label"`
	Required     bool           `json:"required"`
	Readonly     bool           `json:"readonly"`
	TableHidden  bool           `json:"tableHidden,omitempty"`
	EditHidden   bool           `json:"editHidden,omitempty"`
	CreateHidden bool           `json:"createHidden,omitempty"`
	EnumValues   []string       `json:"enumValues,omitempty"`
	Options      []SelectOption `json:"options,omitempty"`
	UploadURL    string         `json:"uploadUrl,omitempty"`
	Accept       string         `json:"accept,omitempty"`
}

type SelectOption struct {
	Value any    `json:"value"`
	Label string `json:"label"`
}
