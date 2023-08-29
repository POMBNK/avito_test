package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/segment/mocks"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_CreateSegment(t *testing.T) {

	type args struct {
		body *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "created",
			args: args{
				body: bytes.NewBufferString(`{"name": "Test Segment","percent": 10}`),
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/segments", tt.args.body)
		w := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {

			service := mocks.NewService(t)
			service.On("Create", r.Context(), mock.Anything).Return("segmentID", nil)

			h := &handler{
				logs:    logger.GetLogger(),
				service: service,
			}

			err := h.CreateSegment(w, r)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, w.Code)
			expectedLocation := fmt.Sprintf("%s/%s", segmentsURL, "segmentID")
			assert.EqualValues(t, expectedLocation, w.Header().Get("Location"))
		})
	}
}

func Test_handler_CronJobSegments(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/segments/ttl", nil)
		w := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {

			service := mocks.NewService(t)
			service.On("CheckSegmentsTTL", r.Context()).Return(nil)

			h := &handler{
				logs:    logger.GetLogger(),
				service: service,
			}
			err := h.CronJobSegments(w, r)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func Test_handler_DeleteSegment(t *testing.T) {
	type args struct {
		body *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "deleted",
			args: args{
				body: bytes.NewBufferString(`{"name": "Test Segment"}`),
			},
		},
	}

	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodDelete, "/segments", tt.args.body)
		w := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewService(t)
			service.On("Delete", r.Context(), mock.Anything).Return(nil)

			h := &handler{
				logs:    logger.GetLogger(),
				service: service,
			}

			err := h.DeleteSegment(w, r)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, w.Code)
		})
	}
}
func Test_handler_EditUserSegments(t *testing.T) {

	type args struct {
		body *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "updated",
			args: args{
				body: bytes.NewBufferString(`{
    "userID":"2",
    "add":[
        {
			"name":"discount80"
        }
    ],
    "delete":["discount80"]}`),
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/segments", tt.args.body)
		w := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {

			service := mocks.NewService(t)
			service.On("EditUserToSegments", r.Context(), mock.Anything).Return(nil)

			h := &handler{
				logs:    logger.GetLogger(),
				service: service,
			}
			err := h.EditUserSegments(w, r)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func Test_handler_GetActiveSegmentFromUser(t *testing.T) {

	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				id: "1",
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, "/segments/"+tt.args.id, nil)
		w := httptest.NewRecorder()
		params := httprouter.Params{httprouter.Param{
			Key:   "user_id",
			Value: tt.args.id,
		}}
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, params)

		t.Run(tt.name, func(t *testing.T) {

			service := mocks.NewService(t)
			service.On("GetActiveSegments", mock.Anything, params.ByName("user_id")).Return([]segment.ActiveSegments{}, nil)

			h := &handler{
				logs:    logger.GetLogger(),
				service: service,
			}
			err := h.GetActiveSegmentFromUser(w, r.WithContext(ctx))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func Test_handler_GetCSVReport(t *testing.T) {

	type args struct {
		id   string
		body *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				id:   "1",
				body: bytes.NewBufferString(`{"month": "August","year": "2023"}`),
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/reports"+tt.args.id, tt.args.body)
		w := httptest.NewRecorder()
		params := httprouter.Params{httprouter.Param{
			Key:   "user_id",
			Value: tt.args.id,
		}}
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, params)

		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewService(t)
			service.On("GetUserHistoryOptimized", mock.Anything, tt.args.id, mock.Anything).Return("linkToFile", nil)

			h := &handler{
				logs:    logger.GetLogger(),
				service: service,
			}

			err := h.GetCSVReport(w, r.WithContext(ctx))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func Test_handler_GetOriginalCSVReport(t *testing.T) {
	type args struct {
		id   string
		body *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				id:   "1",
				body: bytes.NewBufferString(`{"month": "August","year": "2023"}`),
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/reports"+tt.args.id, tt.args.body)
		w := httptest.NewRecorder()
		params := httprouter.Params{httprouter.Param{
			Key:   "user_id",
			Value: tt.args.id,
		}}
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, params)

		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewService(t)
			service.On("GetUserHistoryOriginal", mock.Anything, tt.args.id, mock.Anything).Return("linkToFile", nil)

			h := &handler{
				logs:    logger.GetLogger(),
				service: service,
			}

			err := h.GetOriginalCSVReport(w, r.WithContext(ctx))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
