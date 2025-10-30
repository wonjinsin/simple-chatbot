package http

import (
	"net/http"

	"github.com/wonjinsin/simple-chatbot/pkg/utils"
)

// HealthController handles health check endpoints
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check handles health check requests
func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	utils.WriteStandardJSON(w, r, http.StatusOK, nil)
}
