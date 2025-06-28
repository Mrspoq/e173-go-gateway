package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware creates a middleware that validates JWT tokens
func JWTMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from Authorization header first
		authHeader := c.GetHeader("Authorization")
		var tokenString string
		
		if authHeader != "" {
			tokenString = ExtractTokenFromHeader(authHeader)
		}
		
		// If no token in header, try to get from cookie
		if tokenString == "" {
			cookie, err := c.Cookie("auth_token")
			if err == nil {
				tokenString = cookie
			}
		}
		
		// If still no token, check for session-based auth (for HTMX requests)
		if tokenString == "" {
			sessionToken, exists := c.Get("session_token")
			if exists {
				if token, ok := sessionToken.(string); ok {
					tokenString = token
				}
			}
		}

		if tokenString == "" {
			// For HTMX requests, redirect to login
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			// For regular requests, return JSON error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			// For HTMX requests, redirect to login
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user info in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// RoleMiddleware creates a middleware that checks user roles
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role information not found"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role information"})
			c.Abort()
			return
		}

		// Admin role has access to everything
		if role == "admin" {
			c.Next()
			return
		}

		// Check if user has the required role
		if !hasRole(role, requiredRole) {
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRole checks if the user role satisfies the required role
func hasRole(userRole, requiredRole string) bool {
	// Define role hierarchy
	roleHierarchy := map[string]int{
		"user":     1,
		"operator": 2,
		"admin":    3,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return strings.EqualFold(userRole, requiredRole)
	}

	return userLevel >= requiredLevel
}
