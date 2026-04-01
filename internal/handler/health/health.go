package health

import (
	"back/pkg/httputils"
	"net/http"
	healthdto "back/internal/dto/health"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	httputils.RespondJSON(w, http.StatusOK, healthdto.HealthResponse{Status: "OK"})
}
