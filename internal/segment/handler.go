package segment

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/POMBNK/avito_test_task/docs"
	"github.com/POMBNK/avito_test_task/internal/apierror"
	"github.com/POMBNK/avito_test_task/internal/handlers"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

const (
	segmentsURL       = "/api/segments"
	usersToSegments   = "/api/segments/"
	id                = "user_id"
	activeSegmentsURL = "/api/segments/:user_id"
	csvReport         = "/api/reports/:user_id"
	csvReportOriginal = "/api/original_reports/:user_id"
	swagger           = "/swagger/*filepath"
)

type Service interface {
	Create(ctx context.Context, dto ToCreateSegmentDTO) (string, error)
	Delete(ctx context.Context, dto ToDeleteSegmentDTO) error
	EditUserToSegments(ctx context.Context, dto ToUpdateUsersSegmentsDTO) error
	GetActiveSegments(ctx context.Context, userID string) ([]ActiveSegments, error)
	GetUserHistoryOptimized(ctx context.Context, userID string, dto ReportDateDTO) (string, error)
	GetUserHistoryOriginal(ctx context.Context, userID string, dto ReportDateDTO) (string, error)
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
	r.HandlerFunc(http.MethodPost, csvReport, apierror.Middleware(h.GetCSVReport))
	r.HandlerFunc(http.MethodPost, csvReportOriginal, apierror.Middleware(h.GetOriginalCSVReport))
	r.HandlerFunc(http.MethodGet, swagger, httpSwagger.WrapHandler)
}

// @Summary Create segment
// @Tags segments
// @Description create segment
// @ID create-segment
// @Accept  json
// @Produce  json
// @Param input body segment.ToCreateSegmentDTO true "segment info"
// @Success 201
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments [post]
func (h *handler) CreateSegment(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Create segment")

	w.Header().Set("Content-Type", "application/json")

	var segmentDTO ToCreateSegmentDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	if err := json.NewDecoder(r.Body).Decode(&segmentDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	segmentID, err := h.service.Create(r.Context(), segmentDTO)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", segmentsURL, segmentID))
	w.WriteHeader(http.StatusCreated)
	return nil
}

// @Summary Delete segment
// @Tags segments
// @Description delete segment
// @ID delete-segment
// @Accept  json
// @Produce  json
// @Param input body segment.ToDeleteSegmentDTO true "segment info"
// @Success 204
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments [delete]
func (h *handler) DeleteSegment(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Delete segment")
	w.Header().Set("Content-Type", "application/json")

	var segmentDTO ToDeleteSegmentDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	if err := json.NewDecoder(r.Body).Decode(&segmentDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	err := h.service.Delete(r.Context(), segmentDTO)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// @Summary Edit segments to user
// @Tags segments
// @Description Adds and removes tags to the user
// @ID edit-segment
// @Accept  json
// @Produce  json
// @Param input body segment.SegmentsUsers true "User's segment info"
// @Success 200
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/segments [put]
func (h *handler) EditUserSegments(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Add segments to user")
	w.Header().Set("Content-Type", "application/json")
	//TODO: parse userID from URL or JSON?
	//params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	//userID := params.ByName(id)

	var segmentsDTO ToUpdateUsersSegmentsDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&segmentsDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	err := h.service.EditUserToSegments(r.Context(), segmentsDTO)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// @Summary Get active segments by id
// @Tags segments
// @Description Get all active user's segments by userID
// @ID get-active-segments
// @Accept  json
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

	w.Write(activeSegmentsBytes)
	w.WriteHeader(http.StatusOK)

	return nil

}

// @Summary Get user's report v2
// @Tags reports
// @Description Receiving CSV report with all user actions with segments by user ID, year and month.
// @Description Report v2 has better query and format to read.
// @ID get-report-v2
// @Accept  json
// @Produce  json
// @Param userID path int true "userID"
// @Param input body segment.ReportDateDTO true "Month and year as the beginning of the time period for the CSV report"
// @Success 200 {string} linkToFile
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/reports/{userID} [post]
func (h *handler) GetCSVReport(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Get CSV report")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName(id)

	var dateDTO ReportDateDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	if err := json.NewDecoder(r.Body).Decode(&dateDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	reportLink, err := h.service.GetUserHistoryOptimized(r.Context(), userID, dateDTO)
	if err != nil {
		return err
	}

	w.Write([]byte(reportLink))
	w.WriteHeader(http.StatusOK)
	return nil
}

// @Summary Get user's report v1
// @Tags reports
// @Description Receiving CSV report with all user actions with segments by user ID, year and month.
// @Description Report v1 exist according to the terms of reference about CSV reports
// @ID get-report-v1
// @Accept  json
// @Produce  json
// @Param userID path int true "userID"
// @Param input body segment.ReportDateDTO true "Month and year as the beginning of the time period for the CSV report"
// @Success 200 {string} linkToFile
// @Failure 400,404 {object} apierror.ApiError
// @Failure 500 {object} apierror.ApiError
// @Failure default {object} apierror.ApiError
// @Router /api/original_reports/{userID} [post]
func (h *handler) GetOriginalCSVReport(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Get CSV report")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName(id)

	var dateDTO ReportDateDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	if err := json.NewDecoder(r.Body).Decode(&dateDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	reportLink, err := h.service.GetUserHistoryOriginal(r.Context(), userID, dateDTO)
	if err != nil {
		return err
	}

	w.Write([]byte(reportLink))
	w.WriteHeader(http.StatusOK)
	return nil
}

func NewHandler(logs *logger.Logger, service Service) handlers.Handler {
	return &handler{
		logs:    logs,
		service: service,
	}
}
