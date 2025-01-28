package auth

import (
	"net/http"
	"strings"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if CurrentUser.Email == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			userRole := CurrentUser.Role
			for _, role := range allowedRoles {
				if strings.EqualFold(userRole, role) {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
