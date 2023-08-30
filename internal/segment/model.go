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
	ID     string `json:"ID"`
	UserID string `json:"userID"`
	Add    []struct {
		Name    string `json:"name"`
		TtlDays int    `json:"ttl_days,omitempty"`
	} `json:"add"`
	Delete []string `json:"delete"`
}

type BetterCSVReport struct {
	UserID      string
	SegmentName string
	Active      bool
	CreatedAt   time.Time
	DeletedAt   interface{}
}

type CSVReport struct {
	UserID      string    `json:"user_id"`
	SegmentName string    `json:"segment_name"`
	Action      string    `json:"action"`
	Date        time.Time `json:"date"`
}

type ReportFile struct {
	Data []byte
	Name string
}
