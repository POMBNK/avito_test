package segment

// Data Transfer Object to transport using models
type ToCreateSegmentDTO struct {
	Name    string `json:"name"`
	Percent int    `json:"percent,omitempty"`
}

type ToDeleteSegmentDTO struct {
	Name string `json:"name"`
}

type ToUpdateUsersSegmentsDTO struct {
	UserID string `json:"userID"`
	Add    []struct {
		Name    string `json:"name"`
		TtlDays int    `json:"ttl_days,omitempty"`
	} `json:"add"`
	Delete []string `json:"delete"`
}

type ReportDateDTO struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

// CreateSegmentDto map DTO fields to model
func CreateSegmentDto(dto ToCreateSegmentDTO) Segment {
	return Segment{
		Name: dto.Name,
	}
}

// CreateSegmentDto map DTO fields to model
func DeleteSegmentDto(dto ToDeleteSegmentDTO) Segment {
	return Segment{
		Name: dto.Name,
	}
}

func UpdateUsersSegmentsDto(dto ToUpdateUsersSegmentsDTO) SegmentsUsers {
	return SegmentsUsers{
		UserID: dto.UserID,
		Add:    dto.Add,
		Delete: dto.Delete,
	}
}
