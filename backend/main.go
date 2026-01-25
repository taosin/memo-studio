package main

import (
	"embed"
	"io/fs"
	"log"
	"memo-studio/backend/database"
	"memo-studio/backend/handlers"
	"memo-studio/backend/middleware"
	"net/http"
	"os"
	"path"
	"strings"

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

	// 创建 Gin 路由
	r := gin.Default()

	// 配置 CORS
	// - 生产建议显式设置 MEMO_CORS_ORIGINS（逗号分隔）
	// - 未设置时默认放开（方便自部署），但可按需收紧
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
		config.AllowAllOrigins = true
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "memo-studio-backend"})
	})

	// 附件静态服务（/uploads -> 本地存储目录）
	// 默认 ./storage，可通过 MEMO_STORAGE_DIR 配置
	storageDir := os.Getenv("MEMO_STORAGE_DIR")
	if strings.TrimSpace(storageDir) == "" {
		storageDir = "./storage"
	}
	r.Static("/uploads", storageDir)

	// 公开路由（仅登录/注册）
	public := r.Group("/api")
	{
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/register", handlers.Register)
	}

	// 需要认证的路由（保留旧能力）
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/auth/me", handlers.GetCurrentUser)
		// 用户信息（新接口）
		api.GET("/users/me", handlers.GetMe)
		api.PUT("/users/me", handlers.UpdateMe)
		api.PUT("/users/me/password", handlers.ChangeMyPassword)

		// memos（新接口：需要登录）
		api.GET("/memos", handlers.ListMemos)
		api.POST("/memos", handlers.CreateMemo)
		api.PUT("/memos/:id", handlers.UpdateMemo)
		api.DELETE("/memos/:id", handlers.DeleteMemo)

		// tags（需要登录）
		api.GET("/tags", handlers.GetTags)
		api.POST("/tags", handlers.CreateTag)
		api.PUT("/tags/:id", handlers.UpdateTag)
		api.DELETE("/tags/:id", handlers.DeleteTag)
		api.POST("/tags/merge", handlers.MergeTags)

		// review（需要登录）
		api.GET("/review/random", handlers.RandomReview)

		// resources（附件上传）
		api.POST("/resources", handlers.UploadResource)

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

	// 静态文件托管（用于部署：Go 服务直接提供前端）
	staticFS, err := fs.Sub(publicFiles, "public")
	if err != nil {
		log.Fatal("静态文件目录初始化失败:", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	// SPA fallback：非 /api 路径都回退到 index.html
	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "/health") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// 尝试直接命中静态资源
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

		// 兜底：返回 index.html（给前端路由）
		c.Request.URL.Path = "/index.html"
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})

	// 启动服务器
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "9000"
	}
	log.Println("服务器启动在 :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
