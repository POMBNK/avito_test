package useCase

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/segment/delivery/http"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/POMBNK/avito_test_task/pkg/utils"
	"github.com/jszwec/csvutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const reportPath = "reports/csv/"

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name Storage
type Storage interface {
	// Create method.
	// Create a segment database entries with "name" of segment and "active" boolean flag.
	// Active -> true means segment is active otherwise active -> false.
	Create(ctx context.Context, segment segment.Segment) (string, error)

	// Delete method.
	// Update a segment field "active" to false (0).
	// The field is not deleted from table:
	//   - Not to corrupt the data in the user entity;
	//   - Save statistic data on future.
	Delete(ctx context.Context, segment segment.Segment) error

	// AddUserToSegments Method for adding a user to a segment.
	//Accepts a list of (names) of segments to add a user to
	AddUserToSegments(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName, deleteAfter string) error

	// DeleteSegmentFromUser Method for removing a user from segment.
	//Accepts a list of (names) of segments to delete from user
	DeleteSegmentFromUser(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName string) error

	//GetActiveSegments Method to get active segments from all users
	GetActiveSegments(ctx context.Context, userID string) ([]segment.ActiveSegments, error)

	//IsUserExist check if user already exist
	IsUserExist(ctx context.Context, segmentsUser segment.SegmentsUsers) error

	GetUserHistoryOptimized(ctx context.Context, userID, timestampz string) ([]segment.BetterCSVReport, error)

	GetUserHistoryOriginal(ctx context.Context, userID string, timestampz string) ([]segment.CSVReport, error)

	CheckSegmentsTTL(ctx context.Context) error
}

type service struct {
	logs    *logger.Logger
	storage Storage
}

func (s *service) Create(ctx context.Context, dto segment.ToCreateSegmentDTO) (string, error) {
	segmentUnit := segment.CreateSegmentDto(dto)

	ID, err := s.storage.Create(ctx, segmentUnit)
	if err != nil {
		return "", err
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

func (s *service) GetUserHistoryOptimized(ctx context.Context, userID string, dto segment.ReportDateDTO) (string, error) {
	intYear, err := strconv.Atoi(dto.Year)
	if err != nil {
		return "", err
	}
	timestampz, err := utils.MapToTimestampz(dto.Month, intYear)

	reports, err := s.storage.GetUserHistoryOptimized(ctx, userID, timestampz)
	if err != nil {
		return "", err
	}
	if len(reports) == 0 {
		return "", fmt.Errorf("empty report")
	}

	link, err := prepareCSVReport(reports, userID)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (s *service) GetUserHistoryOriginal(ctx context.Context, userID string, dto segment.ReportDateDTO) (string, error) {
	intYear, err := strconv.Atoi(dto.Year)
	if err != nil {
		return "", err
	}
	timestampz, err := utils.MapToTimestampz(dto.Month, intYear)

	reports, err := s.storage.GetUserHistoryOriginal(ctx, userID, timestampz)
	if err != nil {
		return "", err
	}
	if len(reports) == 0 {
		return "", fmt.Errorf("empty report")
	}

	link, err := prepareOriginalCSVReports(reports, userID)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (s *service) CheckSegmentsTTL(ctx context.Context) error {
	return s.storage.CheckSegmentsTTL(ctx)
}

func prepareCSVReport(reports []segment.BetterCSVReport, userID string) (string, error) {
	b, err := csvutil.Marshal(reports)
	if err != nil {
		return "", err
	}

	createdAt := strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "_"), ":", "_")
	fileName := fmt.Sprintf("report_userID_%s_%s.csv", userID, createdAt)
	err = os.MkdirAll(reportPath, 0744)
	if err != nil && os.IsExist(err) {
		return "", err
	}
	err = os.WriteFile(reportPath+fileName, b, 0744)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(reportPath + fileName)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func prepareOriginalCSVReports(reports []segment.CSVReport, userID string) (string, error) {
	b, err := csvutil.Marshal(reports)
	if err != nil {
		return "", err
	}

	createdAt := strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "_"), ":", "_")
	fileName := fmt.Sprintf("report_userID_%s_%s.csv", userID, createdAt)
	err = os.MkdirAll(reportPath, 0744)
	if err != nil && os.IsExist(err) {
		return "", err
	}
	err = os.WriteFile(reportPath+fileName, b, 0744)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(reportPath + fileName)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func NewService(logs *logger.Logger, storage Storage) http.Service {
	return &service{
		logs:    logs,
		storage: storage,
	}
}
