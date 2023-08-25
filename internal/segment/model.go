package segment

// Segment is a struct of segment model.
type Segment struct {
	ID     string `json:"ID"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type SegmentsUsers struct {
	ID     string   `json:"ID"`
	UserID string   `json:"userID"`
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}
