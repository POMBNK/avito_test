package http

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/POMBNK/avito_test_task/docs"
	"github.com/POMBNK/avito_test_task/internal/apierror"
	"github.com/POMBNK/avito_test_task/internal/handlers"
	"github.com/POMBNK/avito_test_task/internal/responses"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

const (
	segmentsURL                = "/api/segments"
	usersToSegments            = "/api/segments/"
	id                         = "user_id"
	activeSegmentsURL          = "/api/segments/:user_id"
	swagger                    = "/swagger/*filepath"
	cronURL                    = "/api/segments/ttl"
	csvReportDownload          = "/api/reports/download/:user_id"
	csvReportOptimizedDownload = "/api/reports/optimized/download/:user_id"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.0 --name Service
type Service interface {
	Create(ctx context.Context, dto segment.ToCreateSegmentDTO) (string, error)
	Delete(ctx context.Context, dto segment.ToDeleteSegmentDTO) error
	EditUserToSegments(ctx context.Context, dto segment.ToUpdateUsersSegmentsDTO) error
	GetActiveSegments(ctx context.Context, userID string) ([]segment.ActiveSegments, error)
	CheckSegmentsTTL(ctx context.Context) error
	MakeCSVUserReport(ctx context.Context, userID string, dto segment.ReportDateDTO) (segment.ReportFile, error)
	MakeCSVUserReportOptimized(ctx context.Context, userID string, dto segment.ReportDateDTO) (segment.ReportFile, error)
}

type handler struct {
	logs    *logger.Logger
	service Service
}

func (h *handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, segmentsURL, apierror.Middleware(h.CreateSegment))
	r.HandlerFunc(http.MethodDelete, segmentsURL, apierror.Middleware(h.DeleteSegment))
	r.HandlerFunc(http.MethodPut, usersToSegments, apierror.Middleware(h.EditUserSegments))
	r.HandlerFunc(http.MethodGet, activeSegmentsURL, apierror.Middleware(h.GetActiveSegmentFromUser))
	r.HandlerFunc(http.MethodGet, csvReportDownload, apierror.Middleware(h.DownloadCSVUserReport))
	r.HandlerFunc(http.MethodGet, csvReportOptimizedDownload, apierror.Middleware(h.DownloadCSVUserReportOptimized))
	r.HandlerFunc(http.MethodPost, cronURL, apierror.Middleware(h.CronJobSegments))
	r.HandlerFunc(http.MethodGet, swagger, httpSwagger.WrapHandler)
}

// @Summary Create segment
// @Tags segments
// @Description create segment
// @ID create-segment
// @Accept  json
// @Produce  json
// @Param input body segment.ToCreateSegmentDTO true "segment info"
// @Success 201 {object} responses.ApiResponse
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments [post]
func (h *handler) CreateSegment(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Create segment")

	w.Header().Set("Content-Type", "application/json")

	var segmentDTO segment.ToCreateSegmentDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	if err := json.NewDecoder(r.Body).Decode(&segmentDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	segmentID, err := h.service.Create(r.Context(), segmentDTO)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("%s/%s", segmentsURL, segmentID))
	w.Write(responses.Created.Marshal())

	return nil
}

// @Summary Delete segment
// @Tags segments
// @Description delete segment
// @ID delete-segment
// @Accept  json
// @Produce  json
// @Param input body segment.ToDeleteSegmentDTO true "segment info"
// @Success 200 {object} responses.ApiResponse
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments [delete]
func (h *handler) DeleteSegment(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Delete segment")
	w.Header().Set("Content-Type", "application/json")
	
	var segmentDTO segment.ToDeleteSegmentDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	if err := json.NewDecoder(r.Body).Decode(&segmentDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	err := h.service.Delete(r.Context(), segmentDTO)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responses.Deleted.Marshal())

	return nil
}

// @Summary Edit segments to user
// @Tags segments
// @Description Adds already created segments to users and removes segments from the user
// @Description Adds ttl if it needs.
// @ID edit-segment
// @Accept  json
// @Produce  json
// @Param input body segment.ToUpdateUsersSegmentsDTO true "segment to add/delete info"
// @Success 200 {object} responses.ApiResponse
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments [put]
func (h *handler) EditUserSegments(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Add segments to user")
	w.Header().Set("Content-Type", "application/json")

	var segmentsDTO segment.ToUpdateUsersSegmentsDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&segmentsDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	err := h.service.EditUserToSegments(r.Context(), segmentsDTO)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responses.Updated.Marshal())

	return nil
}

// @Summary Get active segments by id
// @Tags segments
// @Description Get all active user's segments by userID
// @ID get-active-segments
// @Produce  json
// @Param userID path int true "userID"
// @Success 200 {array} segment.ActiveSegments
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments/{userID} [get]
func (h *handler) GetActiveSegmentFromUser(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Get active segments from user")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName(id)

	activeSegments, err := h.service.GetActiveSegments(r.Context(), userID)
	if err != nil {
		return err
	}

	activeSegmentsBytes, err := json.Marshal(activeSegments)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(activeSegmentsBytes)

	return nil

}

// @Summary Get user's report v2
// @Tags reports
// @Description Receiving CSV report with all user actions with segments by user ID, year and month.
// @Description Report v2 has better query and format to read.
// @ID get-report-v2
// @Produce octet-stream
// @Param userID path int true "userID"
// @Param month query string true "Month as the beginning of the time period for the CSV report"
// @Param year query string true "Year as the beginning of the time period for the CSV report"
// @Success 200
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/reports/optimized/download/{userID} [get]
func (h *handler) DownloadCSVUserReportOptimized(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Get CSV report v2")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName(id)

	var dateDTO segment.ReportDateDTO
	dateDTO.Month = r.URL.Query().Get("month")
	dateDTO.Year = r.URL.Query().Get("year")
	file, err := h.service.MakeCSVUserReportOptimized(r.Context(), userID, dateDTO)

	if err != nil {
		return err
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(file.Data)
	w.WriteHeader(http.StatusOK)
	return nil
}

// @Summary Get user's report v1
// @Tags reports
// @Description Receiving CSV report with all user actions with segments by user ID, year and month.
// @Description Report v1 created by original requirements.
// @ID get-report-v1
// @Produce octet-stream
// @Param userID path int true "userID"
// @Param month query string true "Month as the beginning of the time period for the CSV report"
// @Param year query string true "Year as the beginning of the time period for the CSV report"
// @Success 200
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/reports/download/{userID} [get]
func (h *handler) DownloadCSVUserReport(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Get CSV report v1")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName(id)

	var dateDTO segment.ReportDateDTO
	dateDTO.Month = r.URL.Query().Get("month")
	dateDTO.Year = r.URL.Query().Get("year")
	file, err := h.service.MakeCSVUserReport(r.Context(), userID, dateDTO)

	if err != nil {
		return err
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(file.Data)
	w.WriteHeader(http.StatusOK)

	return nil
}

// @Summary Cron job to check ttl segments
// @Tags ttl
// @Description CronJobSegments is a function that handles cron job requests for checking segment TTL.
// @ID ttl
// @Produce  json
// @Success 200
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments/ttl [post]
func (h *handler) CronJobSegments(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Cron job running...")
	w.Header().Set("Content-Type", "application/json")

	err := h.service.CheckSegmentsTTL(context.Background())
	if err != nil {
		return err
	}
	h.logs.Info("Cron job done")
	w.Write(responses.OK.Marshal())
	w.WriteHeader(http.StatusOK)
	return nil
}

func NewHandler(logs *logger.Logger, service Service) handlers.Handler {
	return &handler{
		logs:    logs,
		service: service,
	}
}
