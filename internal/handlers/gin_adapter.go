package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// WrapHandler converts a standard http.HandlerFunc to gin.HandlerFunc
func WrapHandler(h http.HandlerFunc) gin.HandlerFunc {
	return gin.WrapH(h)
}

// WrapMiddleware converts a standard HTTP middleware to Gin middleware
func WrapMiddleware(middleware func(http.HandlerFunc) http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a wrapped handler that calls c.Next()
		wrappedHandler := middleware(func(w http.ResponseWriter, r *http.Request) {
			// Update the context with any modifications from middleware
			c.Request = r
			c.Next()
		})
		
		// Call the middleware with our wrapped handler
		wrappedHandler(c.Writer, c.Request)
	}
}

// WrapRoleMiddleware converts a role-based middleware that returns a middleware function
func WrapRoleMiddleware(middlewareFunc func(string) func(http.HandlerFunc) http.HandlerFunc, role string) gin.HandlerFunc {
	middleware := middlewareFunc(role)
	return WrapMiddleware(middleware)
}
