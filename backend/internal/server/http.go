package server

import (
	"backend/internal/handler"
	"backend/internal/pkg/middleware"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"github.com/gin-gonic/gin"
)

func NewServerHTTP(
	logger *log.Logger,
	jwt *jwt.JWT,
	userHandler handler.UserHandler,
	musicHandler handler.MusicHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	// No route group has permission
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.POST("/register", userHandler.Register)
		noAuthRouter.POST("/login", userHandler.Login)
	}
	// Non-strict permission routing group
	noStrictAuthRouter := r.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
	{
		noStrictAuthRouter.GET("/user", userHandler.GetProfile)
	}

	// 音乐相关
	noStrictAuthRouter = r.Group("/music").Use(middleware.NoStrictAuth(jwt, logger))
	{
		// 创建音乐
		noStrictAuthRouter.POST("/create", musicHandler.Create)
		// 上传文件
		noStrictAuthRouter.POST("/upload", musicHandler.Upload)
		// 音乐列表
		noStrictAuthRouter.GET("/list", musicHandler.List)
		// 音乐商店
		noStrictAuthRouter.GET("/detail/:id", musicHandler.GetDetail)
		// 识别音乐
		noStrictAuthRouter.POST("/recognize", musicHandler.Recognize)
		// 创建NFT
		noStrictAuthRouter.POST("/nft/:id", musicHandler.CreateNFT)
		// 获取NFT
		noStrictAuthRouter.GET("/nft/:id", musicHandler.GetNFTs)
		// 获取详情
		noStrictAuthRouter.GET("/nftDetail/:id", musicHandler.GetNFTDetail)
	}

	// Strict permission routing group
	strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(jwt, logger))
	{
		strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
	}

	return r
}
