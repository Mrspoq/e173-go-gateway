package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/e173-gateway/e173_go_gateway/internal/service"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type AuthHandlers struct {
	authService     service.AuthService
	customerService service.CustomerService
}

func NewAuthHandlers(authService service.AuthService, customerService service.CustomerService) *AuthHandlers {
	return &AuthHandlers{
		authService:     authService,
		customerService: customerService,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	User         *models.User  `json:"user,omitempty"`
	SessionToken string        `json:"session_token,omitempty"`
	RedirectURL  string        `json:"redirect_url,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get client info
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Attempt login
	user, session, err := h.authService.Login(req.Username, req.Password, ipAddress, userAgent)
	if err != nil {
		response := LoginResponse{
			Success: false,
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// Determine redirect URL based on user role
	redirectURL := "/dashboard"
	if user.HasRole(models.RoleSuperAdmin) || user.HasRole(models.RoleAdmin) {
		redirectURL = "/admin/dashboard"
	}

	response := LoginResponse{
		Success:      true,
		Message:      "Login successful",
		User:         user,
		SessionToken: session.SessionToken,
		RedirectURL:  redirectURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get session token from cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session found", http.StatusUnauthorized)
		return
	}

	// Logout
	err = h.authService.Logout(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logout successful",
	})
}

func (h *AuthHandlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Change password
	err = h.authService.ChangePassword(user.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})
}

func (h *AuthHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandlers) GetSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user sessions
	sessions, err := h.authService.GetUserSessions(user.ID)
	if err != nil {
		http.Error(w, "Failed to get sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func (h *AuthHandlers) RevokeSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current user from session
	user, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Revoke all user sessions
	err = h.authService.RevokeUserSessions(user.ID)
	if err != nil {
		http.Error(w, "Failed to revoke sessions", http.StatusInternalServerError)
		return
	}

	// Clear current session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "All sessions revoked",
	})
}

// AuthMiddleware validates session and adds user to request context
func (h *AuthHandlers) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate session
		user, err := h.authService.ValidateSession(cookie.Value)
		if err != nil {
			// Clear invalid session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				Path:     "/",
			})
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		r = r.WithContext(WithUser(r.Context(), user))
		next(w, r)
	}
}

// RoleMiddleware checks if user has required role
func (h *AuthHandlers) RoleMiddleware(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !user.HasRole(requiredRole) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}

// Helper function to get current user from request
func (h *AuthHandlers) getCurrentUser(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	return h.authService.ValidateSession(cookie.Value)
}
