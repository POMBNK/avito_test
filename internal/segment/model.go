package segment

import "time"

// Segment is a struct of segment model.
type Segment struct {
	ID     string `json:"ID"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type ActiveSegments struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

type SegmentsUsers struct {
	ID     string   `json:"ID"`
	UserID string   `json:"userID"`
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}

type BetterCSVReport struct {
	UserID      string
	SegmentName string
	Active      bool
	CreatedAt   time.Time
	DeletedAt   interface{}
}
