package segment

// Segment is a struct of segment model.
type Segment struct {
	ID     string `json:"ID"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
