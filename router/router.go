package router

import (
	"TikTokLite/controller"
	"TikTokLite/middleware"
	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	//public directory is used to serve static resources
	//r.Static("/static", "./public")
	r.GET("/static/videos", controller.Videos)
	r.GET("/static/covers", controller.Covers)

	apiRouter := r.Group("/douyin")
	//实现了用户注册，登录，信息的接口
	uc := controller.NewUserController()
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", uc.UserInfo)
	//apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/register/", uc.Register)
	apiRouter.POST("/user/login/", uc.Login)

	publishGroup := apiRouter.Group("/publish", middleware.ValidDataTokenMiddleWare)
	publishGroup.POST("/action/", controller.Publish)
	publishGroup.GET("/list/", controller.PublishList)

	// extra apis - I
	favoriteGroup := apiRouter.Group("/favorite", middleware.ValidDataTokenMiddleWare)
	favoriteGroup.POST("/action/", controller.FavoriteAction)
	favoriteGroup.GET("/list/", controller.FavoriteList)

	commentGroup := apiRouter.Group("/comment", middleware.ValidDataTokenMiddleWare)
	commentGroup.POST("/action/", controller.CommentAction)
	commentGroup.GET("/list/", controller.CommentList)

	// extra apis - II
	rc := controller.NewRelationController()
	apiRouter.POST("/relation/action/", middleware.ValidDataTokenMiddleWare, rc.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.ValidDataTokenMiddleWare, rc.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.ValidDataTokenMiddleWare, rc.FollowerList)
}
