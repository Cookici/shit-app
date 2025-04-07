package api

import (
	"record-project/interfaces/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, userHandler *UserHandler, recordHandler *RecordHandler, tagHandler *TagHandler, poopTypeHandler *PoopTypeHandler, authHandler *AuthHandler, fileHandler *FileHandler, rankingHandler *RankingHandler, friendHandler *FriendHandler) {
	// API版本
	v1 := r.Group("/api/v1")

	// 认证相关路由 - 登录接口不需要认证
	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.WechatLogin)

		// 需要认证的路由
		authProtected := authRoutes.Group("")
		authProtected.Use(middleware.JWTAuthMiddleware())
		{
			authProtected.PUT("/user/:id/info", authHandler.UpdateUserInfo)
			authProtected.POST("/upload", fileHandler.UploadFile)
		}
	}

	// 用户相关路由 - 需要认证
	userRoutes := v1.Group("/users")
	userRoutes.Use(middleware.JWTAuthMiddleware())
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
		userRoutes.GET("/openid/:openid", userHandler.GetUserByOpenID)
		// 添加用户搜索接口
		userRoutes.GET("/search", userHandler.SearchUsers)
	}

	// 记录相关路由 - 需要认证
	recordRoutes := v1.Group("/records")
	recordRoutes.Use(middleware.JWTAuthMiddleware())
	{
		recordRoutes.POST("", recordHandler.CreateRecord)
		recordRoutes.GET("/:id", recordHandler.GetRecord)
		recordRoutes.PUT("/:id", recordHandler.UpdateRecord)
		recordRoutes.DELETE("/:id", recordHandler.DeleteRecord)
		recordRoutes.GET("/user/:user_id", recordHandler.GetRecordsByUserID)
		recordRoutes.GET("/date-range", recordHandler.GetRecordsByDateRange)
		recordRoutes.GET("/daily-stats", recordHandler.GetUsersDailyRecordStats)
	}

	rankingRoutes := v1.Group("/rankings")
	rankingRoutes.Use(middleware.JWTAuthMiddleware())
	{
		rankingRoutes.GET("", rankingHandler.GetRanking)
		rankingRoutes.GET("/friends", rankingHandler.GetFriendRanking)
	}

	// 标签相关路由 - 需要认证
	tagRoutes := v1.Group("/tags")
	tagRoutes.Use(middleware.JWTAuthMiddleware())
	{
		tagRoutes.POST("", tagHandler.CreateTag)
		tagRoutes.GET("/:id", tagHandler.GetTag)
		tagRoutes.PUT("/:id", tagHandler.UpdateTag)
		tagRoutes.DELETE("/:id", tagHandler.DeleteTag)
		tagRoutes.GET("", tagHandler.GetAllTags)
	}

	// 屎的类型相关路由 - 需要认证
	poopTypeRoutes := v1.Group("/poop-types")
	poopTypeRoutes.Use(middleware.JWTAuthMiddleware())
	{
		poopTypeRoutes.GET("", poopTypeHandler.GetAllPoopTypes)
		poopTypeRoutes.GET("/:id", poopTypeHandler.GetPoopType)
	}

	// 好友相关路由 - 需要认证
	friendRoutes := v1.Group("/friends")
	friendRoutes.Use(middleware.JWTAuthMiddleware())
	{
		friendRoutes.GET("", friendHandler.GetFriends)
		friendRoutes.POST("", friendHandler.AddFriend)
		friendRoutes.PUT("/:id", friendHandler.UpdateFriendStatus)
		friendRoutes.DELETE("/:id", friendHandler.DeleteFriend)
		// 添加获取好友申请的路由
		friendRoutes.GET("/requests", friendHandler.GetFriendRequests)
	}
}
