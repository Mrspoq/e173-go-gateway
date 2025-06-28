// Add this to the v1 API group in main.go to fix the blacklist page:

// Blacklist endpoint
v1.GET("/blacklist", func(c *gin.Context) {
    // Return empty blacklist for now
    c.Header("Content-Type", "text/html")
    c.String(http.StatusOK, `
    <tr>
        <td colspan="6" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">
            <div class="py-8">
                <svg class="w-12 h-12 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>
                </svg>
                <p class="text-lg font-medium">No blocked numbers</p>
                <p class="mt-1 text-sm">Add numbers to the blacklist to prevent spam calls</p>
            </div>
        </td>
    </tr>`)
})