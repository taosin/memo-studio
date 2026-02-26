package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"memo-studio/backend/database"
	"memo-studio/backend/handlers"
	"memo-studio/backend/middleware"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 前端静态文件（SvelteKit adapter-static 产物会被同步到 backend/public）
//
//go:embed public/*
var publicFiles embed.FS

func main() {
	// 生产默认使用 release（也可通过 GIN_MODE 覆盖）
	if strings.TrimSpace(os.Getenv("GIN_MODE")) == "" && strings.TrimSpace(os.Getenv("MEMO_ENV")) == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 创建 Gin 路由（生产环境禁用控制台颜色与调试）
	r := gin.New()
	r.Use(gin.Recovery())
	if os.Getenv("GIN_MODE") != "release" {
		r.Use(gin.Logger())
	}

	// 安全响应头
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Robots-Tag", "noindex, nofollow")
		c.Next()
	})

	// 配置 CORS
	config := cors.DefaultConfig()
	if origins := strings.TrimSpace(os.Getenv("MEMO_CORS_ORIGINS")); origins != "" {
		parts := strings.Split(origins, ",")
		allow := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				allow = append(allow, p)
			}
		}
		if len(allow) > 0 {
			config.AllowOrigins = allow
		} else {
			config.AllowAllOrigins = true
		}
	} else {
		// 开发环境默认放开，生产环境建议设置
		if os.Getenv("MEMO_ENV") == "production" {
			log.Printf("[WARNING] 生产环境未设置 MEMO_CORS_ORIGINS，建议配置以提高安全性")
		}
		config.AllowAllOrigins = true
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 健康检查端点（公开，无速率限制）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "memo-studio-backend", "version": "v1"})
	})

	// 附件静态服务（/uploads -> 本地存储目录）
	storageDir := os.Getenv("MEMO_STORAGE_DIR")
	if strings.TrimSpace(storageDir) == "" {
		storageDir = "./storage"
	}
	r.Static("/uploads", storageDir)

	// ===== API v1 =====
	v1 := r.Group("/api/v1")
	{
		// 公开路由（登录/注册）- 带速率限制
		v1.Use(middleware.RateLimitMiddleware())
		{
			v1.POST("/auth/login", handlers.Login)
			v1.POST("/auth/register", handlers.Register)
		}

		// 需要认证的路由
		api := v1.Group("/")
		api.Use(middleware.AuthMiddleware())
		{
			api.GET("/auth/me", handlers.GetCurrentUser)
			api.GET("/users/me", handlers.GetMe)
			api.PUT("/users/me", handlers.UpdateMe)
			api.PUT("/users/me/password", handlers.ChangeMyPassword)

			api.GET("/memos", handlers.ListMemos)
			api.POST("/memos", handlers.CreateMemo)
			api.PUT("/memos/:id", handlers.UpdateMemo)
			api.DELETE("/memos/:id", handlers.DeleteMemo)

			api.GET("/notes", handlers.GetNotes)
			api.POST("/notes", handlers.CreateNote)
			api.GET("/notes/:id", handlers.GetNote)
			api.PUT("/notes/:id", handlers.UpdateNote)
			api.DELETE("/notes/:id", handlers.DeleteNote)
			api.DELETE("/notes/batch", handlers.DeleteNotes)
			api.GET("/search", handlers.SearchNotes)

			api.GET("/tags", handlers.GetTags)
			api.POST("/tags", handlers.CreateTag)
			api.PUT("/tags/:id", handlers.UpdateTag)
			api.DELETE("/tags/:id", handlers.DeleteTag)
			api.POST("/tags/merge", handlers.MergeTags)

			api.GET("/review/random", handlers.RandomReview)

			api.GET("/resources", handlers.ListResources)
			api.POST("/resources", handlers.UploadResource)
			api.POST("/resources/transcribe", handlers.UploadResourceAndTranscribe)
			api.DELETE("/resources/:id", handlers.DeleteResourceHandler)

			// 语音转文本（独立端点）
			api.POST("/speech-to-text", handlers.SpeechToTextOnly)

			api.GET("/notebooks", handlers.ListNotebooks)
			api.GET("/notebooks/:id", handlers.GetNotebook)
			api.POST("/notebooks", handlers.CreateNotebook)
			api.PUT("/notebooks/:id", handlers.UpdateNotebook)
			api.DELETE("/notebooks/:id", handlers.DeleteNotebook)
			api.GET("/notebooks/:id/notes", handlers.ListNotebookNotes)

			api.GET("/stats", handlers.GetStats)
			api.GET("/export", handlers.ExportNotes)
			api.POST("/import", handlers.ImportNotes)

			// AI 洞察与总结
			api.POST("/insights", handlers.GetInsight)
			api.POST("/insights/:type", handlers.GetInsightByType)
			api.POST("/insights/compare", handlers.CompareInsights)
			api.POST("/summarize", handlers.SummarizeNote)
			api.POST("/summarize/batch", handlers.BatchSummarize)

			// AI 辅助（润色、续写、摘要等）
			api.POST("/ai/assist", handlers.AIAssist)

			// 大模型管理
			api.GET("/models", handlers.GetModels)
			api.GET("/models/cloud", handlers.GetCloudModels)
			api.GET("/models/local", handlers.GetLocalModels)
			api.GET("/models/available", handlers.GetAvailableModels)
			api.GET("/models/config", handlers.GetModelConfig)
			api.POST("/models/active", handlers.SetActiveModel)
			api.POST("/models/local", handlers.AddLocalModel)
			api.POST("/models/local/health", handlers.CheckLocalHealth)
			api.POST("/models/test", handlers.TestModelConnection)

			// 位置管理
			api.PUT("/memos/:id/location", handlers.UpdateNoteLocation)
			api.POST("/memos/:id/detect-location", handlers.DetectNoteLocation)
			api.POST("/memos/:id/detect-and-save", handlers.SaveDetectedLocation)
			api.GET("/notes/by-location", handlers.GetNotesByLocation)
			api.GET("/locations/stats", handlers.GetLocationsStats)
			api.POST("/locations/batch-detect", handlers.BatchDetectLocations)

			// 股票分析
			api.GET("/stocks/search", handlers.SearchStocks)
			api.GET("/stocks/hot", handlers.GetHotStocks)
			api.GET("/stocks/:code", handlers.GetStockInfo)
			api.GET("/stocks/:code/history", handlers.GetStockHistory)
			api.POST("/stocks/analyze", handlers.AnalyzeStock)

			// 用户管理（管理员）
			admin := api.Group("/users")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("", handlers.AdminListUsers)
				admin.POST("", handlers.AdminCreateUser)
				admin.PUT("/:id", handlers.AdminUpdateUser)
				admin.DELETE("/:id", handlers.AdminDeleteUser)
			}
		}
	}

	// ===== 旧 API 兼容（已废弃，建议迁移到 /api/v1）=====
	// 登录/注册（无需认证，供前端 /api 前缀使用）
	legacyAuth := r.Group("/api")
	legacyAuth.Use(middleware.RateLimitMiddleware())
	{
		legacyAuth.POST("/auth/login", handlers.Login)
		legacyAuth.POST("/auth/register", handlers.Register)
	}
	// 其余旧 API（需要认证）
	legacy := r.Group("/api")
	legacy.Use(middleware.AuthMiddleware())
	{
		legacy.GET("/auth/me", handlers.GetCurrentUser)
		legacy.GET("/users/me", handlers.GetMe)
		legacy.PUT("/users/me", handlers.UpdateMe)
		legacy.PUT("/users/me/password", handlers.ChangeMyPassword)

		legacy.GET("/memos", handlers.ListMemos)
		legacy.POST("/memos", handlers.CreateMemo)
		legacy.PUT("/memos/:id", handlers.UpdateMemo)
		legacy.DELETE("/memos/:id", handlers.DeleteMemo)

		legacy.GET("/notes", handlers.GetNotes)
		legacy.POST("/notes", handlers.CreateNote)
		legacy.GET("/notes/:id", handlers.GetNote)
		legacy.PUT("/notes/:id", handlers.UpdateNote)
		legacy.DELETE("/notes/:id", handlers.DeleteNote)
		legacy.DELETE("/notes/batch", handlers.DeleteNotes)
		legacy.GET("/search", handlers.SearchNotes)

		legacy.GET("/tags", handlers.GetTags)
		legacy.POST("/tags", handlers.CreateTag)
		legacy.PUT("/tags/:id", handlers.UpdateTag)
		legacy.DELETE("/tags/:id", handlers.DeleteTag)
		legacy.POST("/tags/merge", handlers.MergeTags)

		legacy.GET("/review/random", handlers.RandomReview)

		legacy.GET("/resources", handlers.ListResources)
		legacy.POST("/resources", handlers.UploadResource)
		legacy.POST("/resources/transcribe", handlers.UploadResourceAndTranscribe)
		legacy.DELETE("/resources/:id", handlers.DeleteResourceHandler)

		// 语音转文本（独立端点）
		legacy.POST("/speech-to-text", handlers.SpeechToTextOnly)

		legacy.GET("/notebooks", handlers.ListNotebooks)
		legacy.GET("/notebooks/:id", handlers.GetNotebook)
		legacy.POST("/notebooks", handlers.CreateNotebook)
		legacy.PUT("/notebooks/:id", handlers.UpdateNotebook)
		legacy.DELETE("/notebooks/:id", handlers.DeleteNotebook)
		legacy.GET("/notebooks/:id/notes", handlers.ListNotebookNotes)

		legacy.GET("/stats", handlers.GetStats)
		legacy.GET("/export", handlers.ExportNotes)
		legacy.POST("/import", handlers.ImportNotes)

		// AI 洞察与总结
		legacy.POST("/insights", handlers.GetInsight)
		legacy.POST("/summarize", handlers.SummarizeNote)
		legacy.POST("/summarize/batch", handlers.BatchSummarize)

		// AI 辅助（润色、续写、摘要等）
		legacy.POST("/ai/assist", handlers.AIAssist)

		// 位置管理
		legacy.PUT("/memos/:id/location", handlers.UpdateNoteLocation)
		legacy.POST("/memos/:id/detect-location", handlers.DetectNoteLocation)
		legacy.POST("/memos/:id/detect-and-save", handlers.SaveDetectedLocation)
		legacy.GET("/notes/by-location", handlers.GetNotesByLocation)
		legacy.GET("/locations/stats", handlers.GetLocationsStats)
		legacy.POST("/locations/batch-detect", handlers.BatchDetectLocations)

		// 股票分析
		legacy.GET("/stocks/search", handlers.SearchStocks)
		legacy.GET("/stocks/hot", handlers.GetHotStocks)
		legacy.GET("/stocks/:code", handlers.GetStockInfo)
		legacy.GET("/stocks/:code/history", handlers.GetStockHistory)
		legacy.POST("/stocks/analyze", handlers.AnalyzeStock)

		admin := legacy.Group("/users")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("", handlers.AdminListUsers)
			admin.POST("", handlers.AdminCreateUser)
			admin.PUT("/:id", handlers.AdminUpdateUser)
			admin.DELETE("/:id", handlers.AdminDeleteUser)
		}
	}

	// 静态文件托管
	staticFS, err := fs.Sub(publicFiles, "public")
	if err != nil {
		log.Fatal("静态文件目录初始化失败:", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	// SPA fallback
	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api") || p == "/health" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		clean := strings.TrimPrefix(path.Clean("/"+p), "/")
		if clean == "" || clean == "." {
			clean = "index.html"
		}

		if f, err := staticFS.Open(clean); err == nil {
			_ = f.Close()
			c.Request.URL.Path = "/" + clean
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		c.Request.URL.Path = "/index.html"
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})

	// 启动服务器
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "9000"
	}

	// 生产环境检查必要配置
	if os.Getenv("MEMO_ENV") == "production" {
		if os.Getenv("MEMO_JWT_SECRET") == "" {
			log.Printf("[WARNING] 生产环境未设置 MEMO_JWT_SECRET")
		}
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		log.Printf("🚀 Memo Studio 服务器启动在 :%s (API v1)", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("服务器启动失败:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}
	log.Println("服务器已退出")
}
