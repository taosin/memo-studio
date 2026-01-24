package main

import (
	"log"
	"memo-studio/backend/database"
	"memo-studio/backend/handlers"
	"memo-studio/backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 公开路由（不需要认证）
	public := r.Group("/api")
	{
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/register", handlers.Register)
	}

	// 需要认证的路由
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		api.GET("/auth/me", handlers.GetCurrentUser)

		// 笔记相关
		api.GET("/notes", handlers.GetNotes)
		api.POST("/notes", handlers.CreateNote)
		api.GET("/notes/:id", handlers.GetNote)
		api.PUT("/notes/:id", handlers.UpdateNote)
		api.DELETE("/notes/:id", handlers.DeleteNote)
		api.DELETE("/notes/batch", handlers.DeleteNotes)

		// 标签相关
		api.GET("/tags", handlers.GetTags)
		api.PUT("/tags/:id", handlers.UpdateTag)
		api.DELETE("/tags/:id", handlers.DeleteTag)
		api.POST("/tags/merge", handlers.MergeTags)
	}

	// 启动服务器
	log.Println("服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
