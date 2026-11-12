package domain

// PayDay represents a payday entry.
// This minimal definition includes only fields required for the current UI form.
// Additional fields can be added as the application evolves.
type PayDay struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
