package segment

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/apierror"
	"github.com/POMBNK/avito_test_task/internal/handlers"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	segmentURL  = "/api/segments/:slug"
	segmentsURL = "/api/segments/"
)

type Service interface {
	Create(ctx context.Context, dto ToCreateSegmentDTO) (string, error)
	Delete(ctx context.Context, dto ToDeleteSegmentDTO) error
}

type handler struct {
	logs    *logger.Logger
	service Service
}

func (h *handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, segmentURL, apierror.Middleware(h.CreateSegment))
}

func (h *handler) CreateSegment(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Create segment")

	var segmentDTO ToCreateSegmentDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	segmentName := r.URL.Path[len(segmentsURL):]
	segmentDTO.Name = segmentName

	segmentID, err := h.service.Create(r.Context(), segmentDTO)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", segmentURL, segmentID))
	w.WriteHeader(http.StatusCreated)
	return nil
}

func NewHandler(logs *logger.Logger, service Service) handlers.Handler {
	return &handler{
		logs:    logs,
		service: service,
	}
}
