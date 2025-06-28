package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// Add this to your v1 API group in main.go:
/*
// Blacklist endpoint
v1.GET("/blacklist", func(c *gin.Context) {
    // Return empty blacklist for now
    c.Header("Content-Type", "text/html")
    c.String(http.StatusOK, `
    <tr>
        <td colspan="6" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">
            No blocked numbers yet
        </td>
    </tr>`)
})
*/