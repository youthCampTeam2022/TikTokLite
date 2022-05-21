package router

import (
	"TikTokLite/controller"
	"TikTokLite/middleware"
	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	//public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")
	//实现了用户注册，登录，信息的接口
	uc := controller.NewUserController()
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", uc.UserInfo)
	//apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/register/", uc.Register)
	apiRouter.POST("/user/login/", uc.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	rc := controller.NewRelationController()
	apiRouter.POST("/relation/action/", middleware.ValidDataTokenMiddleWare, rc.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.ValidDataTokenMiddleWare, rc.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.ValidDataTokenMiddleWare, rc.FollowerList)
}
