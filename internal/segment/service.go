package segment

import (
	"context"
	"github.com/POMBNK/avito_test_task/pkg/logger"
)

type Storage interface {
	// Create method.
	// Create a segment database entries with "name" of segment and "active" boolean flag.
	// Active -> true means segment is active otherwise active -> false.
	Create(ctx context.Context, segment Segment) (string, error)

	// Delete method.
	// Update a segment field "active" to false (0).
	// The field is not deleted from table:
	//   - Not to corrupt the data in the user entity;
	//   - Save statistic data on future.
	Delete(ctx context.Context, segment Segment) error

	// AddUserToSegments Method for adding a user to a segment.
	//Accepts a list of (names) of segments to add a user to
	AddUserToSegments(ctx context.Context, segmentsUser SegmentsUsers, segmentName string) error

	//IsUserExist check if user already exist
	IsUserExist(ctx context.Context, segmentsUser SegmentsUsers) error
}

type service struct {
	logs    *logger.Logger
	storage Storage
}

func (s *service) Create(ctx context.Context, dto ToCreateSegmentDTO) (string, error) {
	segmentUnit := CreateSegmentDto(dto)

	ID, err := s.storage.Create(ctx, segmentUnit)
	if err != nil {
		return "", err
	}
	return ID, nil
}

func (s *service) Delete(ctx context.Context, dto ToDeleteSegmentDTO) error {
	segmentUnit := DeleteSegmentDto(dto)

	err := s.storage.Delete(ctx, segmentUnit)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddUserToSegments(ctx context.Context, dto ToUpdateUsersSegmentsDTO) error {
	segmentUnit := UpdateUsersSegmentsDto(dto)
	err := s.storage.IsUserExist(ctx, segmentUnit)
	if err != nil {
		return err
	}

	for _, segmentName := range segmentUnit.Add {
		err = s.storage.AddUserToSegments(ctx, segmentUnit, segmentName)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewService(logs *logger.Logger, storage Storage) Service {
	return &service{
		logs:    logs,
		storage: storage,
	}
}
