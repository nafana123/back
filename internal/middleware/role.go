package middleware

import (
	"back/internal/model"
	"back/pkg/httputils"
	"context"
	"net/http"
	"slices"
)

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

func RequireRole(allowed ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetRole(r.Context())
			if !ok || role == "" {
				httputils.RespondError(w, http.StatusForbidden, "Недостаточно прав")
				return
			}
			if !slices.Contains(allowed, role) {
				httputils.RespondError(w, http.StatusForbidden, "Недостаточно прав")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRole(model.RoleAdmin)
}
