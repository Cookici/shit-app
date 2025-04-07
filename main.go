package main

import (
	"fmt"
	"log"
	"record-project/application/service"
	"record-project/infrastructure/auth"
	"record-project/infrastructure/config"
	"record-project/infrastructure/persistence"
	"record-project/infrastructure/persistence/repository"
	"record-project/infrastructure/storage"
	"record-project/infrastructure/wechat"
	"record-project/interfaces/api"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化JWT密钥
	auth.InitJWT(cfg.JWT.Secret, time.Duration(cfg.JWT.ExpirationHours)*time.Hour)

	// 初始化数据库
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化数据
	if err := db.InitData(); err != nil {
		log.Fatalf("初始化数据失败: %v", err)
	}

	// 初始化仓储
	userRepo := repository.NewUserRepository(db.DB)
	recordRepo := repository.NewRecordRepository(db.DB)
	tagRepo := repository.NewTagRepository(db.DB)
	poopTypeRepo := repository.NewPoopTypeRepository(db.DB)
	recordTagRepo := repository.NewRecordTagRepository(db.DB)
	friendRepo := repository.NewFriendRepository(db.DB)

	// 初始化微信服务
	wechatService := wechat.NewWechatService(cfg.Wechat.AppID, cfg.Wechat.AppSecret)

	// 初始化阿里云OSS服务
	ossService, err := storage.NewOSSService(&cfg.Aliyun)
	if err != nil {
		log.Fatalf("初始化阿里云OSS服务失败: %v", err)
	}

	// 初始化应用服务
	userService := service.NewUserService(userRepo)
	recordService := service.NewRecordService(recordRepo, recordTagRepo)
	tagService := service.NewTagService(tagRepo, recordTagRepo)
	poopTypeService := service.NewPoopTypeService(poopTypeRepo)
	authService := service.NewAuthService(userService, wechatService)
	fileService := service.NewFileService(ossService)
	friendService := service.NewFriendService(friendRepo)

	// 初始化API处理器
	userHandler := api.NewUserHandler(userService, authService, friendService)
	recordHandler := api.NewRecordHandler(recordService, userService, tagService, poopTypeService)
	tagHandler := api.NewTagHandler(tagService)
	poopTypeHandler := api.NewPoopTypeHandler(poopTypeService)
	authHandler := api.NewAuthHandler(authService)
	fileHandler := api.NewFileHandler(fileService)
	rankingHandler := api.NewRankingHandler(recordService, authService, userService, friendService)
	friendHandler := api.NewFriendHandler(friendService, authService, userService)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	api.RegisterRoutes(r, userHandler, recordHandler, tagHandler, poopTypeHandler, authHandler, fileHandler, rankingHandler, friendHandler)

	// 启动服务器
	log.Printf("服务器启动在 :%d 端口", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
