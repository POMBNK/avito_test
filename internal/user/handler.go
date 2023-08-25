package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (h *handler) CreateSegment(w http.ResponseWriter, r *http.Request) error {
	h.logs.Info("Create segment")

	path := r.URL.Path[len(segmentsURL):]
	slugs := strings.Split(path, "/")

	var segmentDTOs ToCreateSegmentDTO
	defer r.Body.Close()
	h.logs.Debug("mapping json to DTO")
	if err := json.NewDecoder(r.Body).Decode(&segmentDTO); err != nil {
		return fmt.Errorf("failled to decode body from json body due error:%w", err)
	}

	segmentID, err := h.service.Create(r.Context(), segmentDTO)
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", segmentURL, segmentID))
	w.WriteHeader(http.StatusCreated)
	return nil
}
