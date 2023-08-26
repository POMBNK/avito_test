package segment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/apierror"
	"github.com/POMBNK/avito_test_task/internal/handlers"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	segmentsURL       = "/api/segments"
	usersToSegments   = "/api/segments/"
	id                = "user_id"
	activeSegmentsURL = "/api/segments/:user_id"
	csvReport         = "/api/reports/:user_id"
)

type Service interface {
	Create(ctx context.Context, dto ToCreateSegmentDTO) (string, error)
	Delete(ctx context.Context, dto ToDeleteSegmentDTO) error
	EditUserToSegments(ctx context.Context, dto ToUpdateUsersSegmentsDTO) error
	GetActiveSegments(ctx context.Context, userID string) ([]ActiveSegments, error)
	GetUserHistoryOptimized(ctx context.Context, userID string, dto ReportDateDTO) (string, error)
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
}

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

	w.WriteHeader(http.StatusCreated)
	return nil
}

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

	return nil
}

func NewHandler(logs *logger.Logger, service Service) handlers.Handler {
	return &handler{
		logs:    logs,
		service: service,
	}
}
