package httputils

import (
	userdto "back/internal/dto/user"
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}
}

func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, userdto.ErrorResponse{
		Error: message,
	})
}

func RespondErrorWithField(w http.ResponseWriter, status int, message, field string) {
	RespondJSON(w, status, userdto.ErrorResponse{
		Error: message,
		Field: field,
	})
}

func RespondDecodeError(w http.ResponseWriter) {
	RespondError(w, http.StatusBadRequest, "Ошибка получения тела запроса")
}
