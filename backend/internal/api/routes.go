package api

import (
	"file-vault/backend/internal/database" // Import database package for AppClients
	"file-vault/backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	// Import postgrest-go
)

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, clients *database.AppClients) { // Accept AppClients
	// Initialize the rate limiter: 2 requests per second with a burst of 4.
	limiter := NewUserRateLimiter(rate.Limit(2), 4)

	// Group routes under /api/v1
	v1 := router.Group("/api/v1")
	v1.Use(RateLimitMiddleware(limiter, clients)) // Apply the rate limiting middleware to all v1 routes
	{
		// User routes
		v1.POST("/register", handlers.RegisterUser(clients.Postgrest)) // Pass Postgrest client
		v1.POST("/login", handlers.LoginUser(clients.Postgrest))       // Pass Postgrest client
		v1.POST("/verify-otp", handlers.VerifyOTP(clients.Postgrest))
		v1.POST("/resend-otp", handlers.ResendOTP(clients.Postgrest))

		// Authenticated user routes
		user := v1.Group("/user")
		user.Use(AuthMiddleware(clients)) // Protect user routes
		{
			user.GET("/quota", handlers.GetUserQuota(clients))
			user.POST("/password", handlers.UpdatePassword(clients))
			user.POST("/files/:id/share", handlers.ShareFile(clients))
		}
		// Publicly shared files route (no authentication required)
		v1.GET("/user/shared-publicly", handlers.ListPubliclySharedFiles(clients))

		// File routes (these should be accessible via v1, not user group)
		v1.POST("/upload", handlers.UploadFile(clients, "balkanid-file-storage")) // Pass the entire clients object and bucket name
		v1.GET("/files", handlers.ListFiles(clients))
		v1.GET("/files/:id", handlers.GetFile(clients, "balkanid-file-storage"))
		v1.DELETE("/files/:id", handlers.DeleteFile(clients, "balkanid-file-storage"))

		// Search route (accessible via v1)
		v1.GET("/search", handlers.SearchFiles(clients))
		v1.GET("/stats", handlers.GetStats(clients))

		// Public sharing routes (no authentication required for GetPublicShare and DownloadPublicShare)
		router.GET("/share/:token", handlers.GetPublicShare(clients, "balkanid-file-storage"))
		router.GET("/share/:token/download", handlers.DownloadPublicShare(clients, "balkanid-file-storage"))

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(AdminAuthMiddleware(clients)) // Protect admin routes
		{
			admin.GET("/files", handlers.AdminListFiles(clients))
			admin.POST("/config", handlers.UpdateConfig(clients))
		}
	}
}
