package useCase

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/apierror"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/segment/delivery/http"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/POMBNK/avito_test_task/pkg/utils"
	"github.com/jszwec/csvutil"
	"strconv"
	"strings"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name Storage
type Storage interface {
	// Create method.
	// Create a segment database entries with "name" of segment and "active" boolean flag.
	// Active -> true means segment is active otherwise active -> false.
	Create(ctx context.Context, segment segment.Segment) (string, error)

	// Delete method.
	// Update a segment field "active" to false.
	// The field is not deleted from table:
	//   - Not to corrupt the data in the user entity;
	//   - Save statistic data on future.
	Delete(ctx context.Context, segment segment.Segment) error

	// AddUserToSegments Method for adding a user to a segment.
	//Accepts a list of segments names to add a user to
	AddUserToSegments(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName, deleteAfter string) error

	// DeleteSegmentFromUser Method for removing a user from segment.
	//Accepts a list of (list of segments names to delete from user
	DeleteSegmentFromUser(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName string) error

	//GetActiveSegments Method to get all active segments belongs to user by UserID
	GetActiveSegments(ctx context.Context, userID string) ([]segment.ActiveSegments, error)

	//IsUserExist check if user already exist
	IsUserExist(ctx context.Context, segmentsUser segment.SegmentsUsers) error

	// GetUserHistoryOptimized retrieves the user history from the PostgreSQL database in an optimized way
	GetUserHistoryOptimized(ctx context.Context, userID, timestampz string) ([]segment.BetterCSVReport, error)

	// GetUserHistoryOptimized retrieves the user history from the PostgreSQL database in an original way
	GetUserHistoryOriginal(ctx context.Context, userID string, timestampz string) ([]segment.CSVReport, error)

	// CheckSegmentsTTL updates the active state of user segments based on their TTL (time-to-live)
	CheckSegmentsTTL(ctx context.Context) error

	//AddToRandomUsers adds random users to a segment in the database.
	AddToRandomUsers(ctx context.Context, segment segment.Segment, percent int) error
}

type service struct {
	logs    *logger.Logger
	storage Storage
}

func (s *service) Create(ctx context.Context, dto segment.ToCreateSegmentDTO) (string, error) {
	segmentUnit := segment.CreateSegmentDto(dto)
	if dto.Name == "" {
		return "", apierror.New("segment name is empty", "validation err", "Avito_Segment_Service-000401")
	}
	ID, err := s.storage.Create(ctx, segmentUnit)
	if err != nil {
		return "", err
	}
	if dto.Percent < 0 {
		return "", apierror.New("percent must be positive", "validation err", "Avito_Segment_Service-000402")
	}
	if dto.Percent != 0 {
		err = s.storage.AddToRandomUsers(ctx, segmentUnit, dto.Percent)
		if err != nil {
			return "", err
		}
	}

	return ID, nil
}

func (s *service) Delete(ctx context.Context, dto segment.ToDeleteSegmentDTO) error {
	segmentUnit := segment.DeleteSegmentDto(dto)

	err := s.storage.Delete(ctx, segmentUnit)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) EditUserToSegments(ctx context.Context, dto segment.ToUpdateUsersSegmentsDTO) error {
	segmentUnit := segment.UpdateUsersSegmentsDto(dto)
	err := s.storage.IsUserExist(ctx, segmentUnit)
	if err != nil {
		return err
	}

	for _, segmentField := range segmentUnit.Add {
		var deleteAfter = ""
		if !(segmentField.TtlDays == 0) {
			deleteAfter = time.Now().AddDate(0, 0, segmentField.TtlDays).Format("2006-01-02 15:04:05-07")
		}
		err = s.storage.AddUserToSegments(ctx, segmentUnit, segmentField.Name, deleteAfter)
		if err != nil {
			return err
		}
	}

	for _, segmentName := range segmentUnit.Delete {
		err = s.storage.DeleteSegmentFromUser(ctx, segmentUnit, segmentName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetActiveSegments(ctx context.Context, userID string) ([]segment.ActiveSegments, error) {
	activeSegments, err := s.storage.GetActiveSegments(ctx, userID)
	if err != nil {
		return nil, err
	}

	return activeSegments, err
}

func (s *service) MakeCSVUserReportOptimized(ctx context.Context, userID string, dto segment.ReportDateDTO) (segment.ReportFile, error) {
	var newReport segment.ReportFile

	intYear, err := strconv.Atoi(dto.Year)
	if err != nil {
		return segment.ReportFile{}, err
	}
	timestampz, err := utils.MapToTimestampz(dto.Month, intYear)

	reports, err := s.storage.GetUserHistoryOptimized(ctx, userID, timestampz)
	if err != nil {
		return segment.ReportFile{}, err
	}
	if len(reports) == 0 {
		return segment.ReportFile{}, apierror.New("empty report", "empty report", "Avito_Segment_Service-000403")
	}

	b, err := csvutil.Marshal(reports)
	if err != nil {
		return segment.ReportFile{}, err
	}
	createdAt := strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "_"), ":", "_")
	fileName := fmt.Sprintf("report_userID_%s_%s.csv", userID, createdAt)

	newReport.Data = b
	newReport.Name = fileName

	return newReport, nil
}

func (s *service) MakeCSVUserReport(ctx context.Context, userID string, dto segment.ReportDateDTO) (segment.ReportFile, error) {
	var newReport segment.ReportFile

	intYear, err := strconv.Atoi(dto.Year)
	if err != nil {
		return segment.ReportFile{}, err
	}
	timestampz, err := utils.MapToTimestampz(dto.Month, intYear)

	reports, err := s.storage.GetUserHistoryOriginal(ctx, userID, timestampz)
	if err != nil {
		return segment.ReportFile{}, err
	}
	if len(reports) == 0 {
		return segment.ReportFile{}, apierror.New("empty report", "empty report", "Avito_Segment_Service-000403")
	}

	b, err := csvutil.Marshal(reports)
	if err != nil {
		return segment.ReportFile{}, err
	}
	createdAt := strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "_"), ":", "_")
	fileName := fmt.Sprintf("report_userID_%s_%s.csv", userID, createdAt)

	newReport.Data = b
	newReport.Name = fileName

	return newReport, nil
}

func (s *service) CheckSegmentsTTL(ctx context.Context) error {
	return s.storage.CheckSegmentsTTL(ctx)
}

func NewService(logs *logger.Logger, storage Storage) http.Service {
	return &service{
		logs:    logs,
		storage: storage,
	}
}
