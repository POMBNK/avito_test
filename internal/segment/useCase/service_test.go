package useCase

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/segment/mocks"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_service_Create(t *testing.T) {

	type args struct {
		ctx context.Context
		dto segment.ToCreateSegmentDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Test_service_Create_1",
			args: args{
				ctx: context.Background(),
				dto: segment.ToCreateSegmentDTO{
					Name: "test_name1",
				},
			},
			want: "id",
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
			storage.On("Create", tt.args.ctx, segment.CreateSegmentDto(tt.args.dto)).Return(tt.want, nil)

			s := &service{
				logs:    logs,
				storage: storage,
			}
			got, err := s.Create(tt.args.ctx, tt.args.dto)
			assert.Equal(t, tt.want, got)
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
		name string
		args args
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
		name string
		args args
	}{
		{
			name: "Test_service_EditUserToSegments_1",
			args: args{
				ctx: context.Background(),
				dto: segment.ToUpdateUsersSegmentsDTO{
					UserID: "user_id_1",
					Add: []struct {
						Name    string `json:"name"`
						TtlDays int    `json:"ttl_days,omitempty"`
					}([]struct {
						Name    string
						TtlDays int
					}{
						{
							Name:    "segment1",
							TtlDays: 1,
						},
					}),
					Delete: []string{"segment2", "segment3"},
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
			deleteAfter := time.Now().AddDate(0, 0, segmentUnit.Add[0].TtlDays).Format("2006-01-02 15:04:05-07")
			for _, segmentName := range segmentUnit.Add {
				storage.On("AddUserToSegments", tt.args.ctx, segmentUnit, segmentName.Name, deleteAfter).Return(nil)
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

func Test_service_EditUserToSegmentsFailed(t *testing.T) {

	type args struct {
		ctx context.Context
		dto segment.ToUpdateUsersSegmentsDTO
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_service_EditUserToSegments_1",
			args: args{
				ctx: context.Background(),
				dto: segment.ToUpdateUsersSegmentsDTO{
					UserID: "user_id_1",
					Add: []struct {
						Name    string `json:"name"`
						TtlDays int    `json:"ttl_days,omitempty"`
					}([]struct {
						Name    string
						TtlDays int
					}{
						{
							Name: "segment1",
						},
					}),
					Delete: []string{"segment2", "segment3"},
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
			deleteAfter := ""
			for _, segmentName := range segmentUnit.Add {
				storage.On("AddUserToSegments", tt.args.ctx, segmentUnit, segmentName.Name, deleteAfter).Return(nil)
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
		name string
		args args
		want []segment.ActiveSegments
	}{
		{
			name: "Test_service_GetActiveSegments_1",
			args: args{
				ctx:    context.Background(),
				userID: "user_id_1",
			},
			want: []segment.ActiveSegments{
				{
					ID:   "segment_id_1",
					Name: "segment_name_1",
				},
				{
					ID:   "segment_id_2",
					Name: "segment_name_2",
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
			storage.On("GetActiveSegments", tt.args.ctx, tt.args.userID).Return(tt.want, nil)

			got, err := s.GetActiveSegments(tt.args.ctx, tt.args.userID)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			storage.AssertExpectations(t)
		})
	}
}

func Test_service_MakeCSVUserReportOptimized(t *testing.T) {

	type args struct {
		ctx    context.Context
		userID string
		dto    segment.ReportDateDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_service_GetUserHistoryOptimized_1",
			args: args{
				ctx:    context.Background(),
				userID: "1",
				dto: segment.ReportDateDTO{
					Month: "august",
					Year:  "2023",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := mocks.NewStorage(t)
			timestampzMock := "2023-08-01T00:00:00Z"
			storage.On("GetUserHistoryOptimized", tt.args.ctx, tt.args.userID, timestampzMock).Return([]segment.BetterCSVReport{
				{
					UserID:      "1",
					SegmentName: "test_name1",
					Active:      true,
					CreatedAt:   time.Now(),
					DeletedAt:   nil,
				},
			}, nil)

			s := &service{
				storage: storage,
			}

			_, err := s.MakeCSVUserReportOptimized(tt.args.ctx, tt.args.userID, tt.args.dto)
			assert.NoError(t, err)
		})
	}
}

func Test_service_MakeCSVUserReportOptimized_Failed(t *testing.T) {

	type args struct {
		ctx    context.Context
		userID string
		dto    segment.ReportDateDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_service_GetUserHistoryOptimized_1",
			args: args{
				ctx:    context.Background(),
				userID: "1",
				dto: segment.ReportDateDTO{
					Month: "august",
					Year:  "2023",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := mocks.NewStorage(t)
			timestampz := "2023-08-01T00:00:00Z"
			storage.On("GetUserHistoryOptimized", tt.args.ctx, tt.args.userID, timestampz).Return([]segment.BetterCSVReport{}, nil)

			s := &service{
				storage: storage,
			}

			_, err := s.MakeCSVUserReportOptimized(tt.args.ctx, tt.args.userID, tt.args.dto)
			assert.Error(t, err, fmt.Errorf("empty report"))
		})
	}
}

func Test_service_MakeCSVUserReport(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		dto    segment.ReportDateDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_service_GetUserHistoryOriginal_1",
			args: args{
				ctx:    context.Background(),
				userID: "1",
				dto: segment.ReportDateDTO{
					Month: "august",
					Year:  "2023",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestampzMock := "2023-08-01T00:00:00Z"
			storage := mocks.NewStorage(t)
			storage.On("GetUserHistoryOriginal", tt.args.ctx, tt.args.userID, timestampzMock).Return([]segment.CSVReport{
				{
					UserID:      "1",
					SegmentName: "test_name1",
					Action:      "created",
					Date:        time.Now(),
				},
			}, nil)

			s := &service{
				storage: storage,
			}
			_, err := s.MakeCSVUserReport(tt.args.ctx, tt.args.userID, tt.args.dto)
			assert.NoError(t, err)
		})
	}
}

func Test_service_MakeCSVUserReportFailed(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		dto    segment.ReportDateDTO
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_service_GetUserHistoryOriginal_1",
			args: args{
				ctx:    context.Background(),
				userID: "1",
				dto: segment.ReportDateDTO{
					Month: "august",
					Year:  "2023",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestampzMock := "2023-08-01T00:00:00Z"
			storage := mocks.NewStorage(t)
			storage.On("GetUserHistoryOriginal", tt.args.ctx, tt.args.userID, timestampzMock).Return([]segment.CSVReport{}, nil)

			s := &service{
				storage: storage,
			}
			_, err := s.MakeCSVUserReport(tt.args.ctx, tt.args.userID, tt.args.dto)
			assert.Error(t, err, fmt.Errorf("empty report"))
		})
	}
}

func Test_service_CheckSegmentsTTL(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test_service_CheckSegmentsTTL_1",
			args: args{
				ctx: context.Background(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {

		storage := mocks.NewStorage(t)
		storage.On("CheckSegmentsTTL", tt.args.ctx).Return(nil)

		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				storage: storage,
			}
			err := s.CheckSegmentsTTL(tt.args.ctx)
			assert.NoError(t, err)
			storage.AssertExpectations(t)
		})
	}
}
