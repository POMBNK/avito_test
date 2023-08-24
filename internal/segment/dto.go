package segment

// Data Transfer Object to transport using models
type ToCreateSegmentDTO struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type ToDeleteSegmentDTO struct {
	Name string `json:"name"`
}

// CreateSegmentDto map DTO fields to model
func CreateSegmentDto(dto ToCreateSegmentDTO) Segment {
	return Segment{
		Name:   dto.Name,
		Active: dto.Active,
	}
}

// CreateSegmentDto map DTO fields to model
func DeleteSegmentDto(dto ToDeleteSegmentDTO) Segment {
	return Segment{
		Name: dto.Name,
	}
}
