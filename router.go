package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/controller"
	"github.com/ikuraoo/fastdouyin/middleware"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", middleware.TokenParse(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middleware.TokenParse(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.TokenParse(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.TokenParse(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.TokenParse(), controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.TokenParse(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.TokenParse(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.TokenParse(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.TokenParse(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
}
