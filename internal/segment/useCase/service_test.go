package useCase

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/segment/mocks"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_prepareCSVReport(t *testing.T) {
	type args struct {
		reports []segment.BetterCSVReport
		userID  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test_prepareCSVReport_1",
			args: args{
				reports: []segment.BetterCSVReport{
					{
						UserID:      "1",
						SegmentName: "test_name1",
						Active:      true,
						CreatedAt:   time.Now(),
						DeletedAt:   nil,
					},
					{
						UserID:      "2",
						SegmentName: "TESTNAME2",
						Active:      true,
						CreatedAt:   time.Now(),
						DeletedAt:   nil,
					},
				},
				userID: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			createdAt := strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "_"), ":", "_")
			fileName := fmt.Sprintf("report_userID_%s_%s.csv", tt.args.userID, createdAt)
			expectedFilePath, _ := filepath.Abs(reportPath + fileName)
			absPath, err := prepareCSVReport(tt.args.reports, tt.args.userID)

			assert.Nil(t, err)
			assert.Equal(t, expectedFilePath, absPath)
			err = os.Remove(absPath)
			assert.Nil(t, err)
		})
	}
}

func Test_prepareOriginalCSVReports(t *testing.T) {
	type args struct {
		reports []segment.CSVReport
		userID  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test_prepareCSVReport_1",
			args: args{
				reports: []segment.CSVReport{
					{
						UserID:      "1",
						SegmentName: "test_name1",
						Action:      "true",
						Date:        time.Now(),
					},
					{
						UserID:      "2",
						SegmentName: "TESTNAME2",
						Action:      "false",
						Date:        time.Now(),
					},
				},
				userID: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			createdAt := strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.Stamp), " ", "_"), ":", "_")
			fileName := fmt.Sprintf("report_userID_%s_%s.csv", tt.args.userID, createdAt)

			expectedFilePath, _ := filepath.Abs(reportPath + fileName)
			absPath, err := prepareOriginalCSVReports(tt.args.reports, tt.args.userID)

			assert.Nil(t, err)
			assert.Equal(t, expectedFilePath, absPath)
			err = os.Remove(absPath)
			assert.Nil(t, err)
		})
	}
}

func Test_service_Create(t *testing.T) {

	type args struct {
		ctx context.Context
		dto segment.ToCreateSegmentDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test_service_Create_1",
			args: args{
				ctx: context.Background(),
				dto: segment.ToCreateSegmentDTO{
					Name: "test_name1",
				},
			},
		},
		{name: "Test_service_Create_2",
			args: args{
				ctx: context.Background(),
				dto: segment.ToCreateSegmentDTO{
					Name: "TESTNAME2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewStorage(t)
			logs := logger.GetLogger()
			storage.On("Create", tt.args.ctx, segment.CreateSegmentDto(tt.args.dto)).Return("id", nil)

			s := &service{
				logs:    logs,
				storage: storage,
			}
			got, err := s.Create(tt.args.ctx, tt.args.dto)
			assert.Equal(t, "id", got)
			assert.NoError(t, err)
			storage.AssertExpectations(t)
		})
	}
}

func Test_service_Delete(t *testing.T) {

	type args struct {
		ctx context.Context
		dto segment.ToDeleteSegmentDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_service_Delete_1",
			args: args{
				ctx: context.Background(),
				dto: segment.ToDeleteSegmentDTO{
					Name: "segment_id_1",
				},
			},
		},
		{
			name: "Test_service_Delete_2",
			args: args{
				ctx: context.Background(),
				dto: segment.ToDeleteSegmentDTO{
					Name: "segment_id_2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := mocks.NewStorage(t)
			storage.On("Delete", tt.args.ctx, segment.DeleteSegmentDto(tt.args.dto)).Return(nil)
			storage.On("Delete", tt.args.ctx, segment.DeleteSegmentDto(tt.args.dto)).Return(nil)

			s := &service{
				storage: storage,
			}
			err := s.Delete(tt.args.ctx, tt.args.dto)
			assert.NoError(t, err)
			storage.AssertExpectations(t)
		})
	}
}

func Test_service_EditUserToSegments(t *testing.T) {

	type args struct {
		ctx context.Context
		dto segment.ToUpdateUsersSegmentsDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test_service_EditUserToSegments_1",
			args: args{
				ctx: context.Background(),
				dto: segment.ToUpdateUsersSegmentsDTO{
					UserID: "user_id_1",
					Add:    []string{"segment1", "segment2"},
					Delete: []string{"segment3"},
				},
			},
		},
		{
			name: "Test_service_EditUserToSegments_2",
			args: args{
				ctx: context.Background(),
				dto: segment.ToUpdateUsersSegmentsDTO{
					UserID: "user_id_2",
					Add:    []string{"segment4"},
					Delete: []string{"segment5", "segment6"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := mocks.NewStorage(t)
			s := &service{
				storage: storage,
			}
			segmentUnit := segment.UpdateUsersSegmentsDto(tt.args.dto)
			storage.On("IsUserExist", tt.args.ctx, segmentUnit).Return(nil)

			for _, segmentName := range segmentUnit.Add {
				storage.On("AddUserToSegments", tt.args.ctx, segmentUnit, segmentName).Return(nil)
			}

			for _, segmentName := range segmentUnit.Delete {
				storage.On("DeleteSegmentFromUser", tt.args.ctx, segmentUnit, segmentName).Return(nil)
			}

			err := s.EditUserToSegments(tt.args.ctx, tt.args.dto)
			assert.NoError(t, err)
			storage.AssertExpectations(t)

		})
	}
}

func Test_service_GetActiveSegments(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []segment.ActiveSegments
		wantErr bool
	}{
		{
			name: "Test_service_GetActiveSegments_1",
			args: args{
				ctx:    context.Background(),
				userID: "user_id_1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := mocks.NewStorage(t)
			s := &service{
				storage: storage,
			}
			storage.On("GetActiveSegments", tt.args.ctx, tt.args.userID).Return(tt.want, nil)

			s.GetActiveSegments(tt.args.ctx, tt.args.userID)
			got, err := s.GetActiveSegments(tt.args.ctx, tt.args.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			storage.AssertExpectations(t)
		})
	}
}

func Test_service_GetUserHistoryOptimized(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		ctx    context.Context
		userID string
		dto    segment.ReportDateDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				storage: tt.fields.storage,
			}
			got, err := s.GetUserHistoryOptimized(tt.args.ctx, tt.args.userID, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserHistoryOptimized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserHistoryOptimized() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetUserHistoryOriginal(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		ctx    context.Context
		userID string
		dto    segment.ReportDateDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				storage: tt.fields.storage,
			}
			got, err := s.GetUserHistoryOriginal(tt.args.ctx, tt.args.userID, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserHistoryOriginal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserHistoryOriginal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
