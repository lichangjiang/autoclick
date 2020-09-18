package controllers

import "github.com/gin-gonic/gin"

type Controller interface {
	routeGroup() string
	config(rg *gin.RouterGroup)
}

var controllerList = []Controller{
	&ScreenController{},
	&ImageController{},
	&ActionController{},
}

func ConfigRouter(router *gin.Engine) {
	router.NoMethod(HandleNotFound)
	router.NoRoute(HandleNotFound)

	for _, controller := range controllerList {
		rg := router.Group(controller.routeGroup())
		controller.config(rg)
	}
}
