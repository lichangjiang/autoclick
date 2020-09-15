package controllers

import (
	"autoclick/pkg/imageutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScreenController struct {
}

func (*ScreenController) routeGroup() string {
	return "/screen"
}

func (*ScreenController) config(rg *gin.RouterGroup) {
	rg.Use(authMiddle)
	rg.GET("/info", wrapper(screenInfoHandler))
	rg.GET("/count",wrapper(screenDisplayCount))
}

func screenInfoHandler(c *gin.Context) error {
	numStr := c.Query("num")
	if numStr == "" {
		c.JSON(http.StatusOK, imageutil.GetAllScreenInfo())
		return nil
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return ParameterError("parameter num=" + numStr + " is not a number")
	}
	c.JSON(http.StatusOK, imageutil.GetScreenInfo(num))
	return nil
}

func screenDisplayCount(c *gin.Context) error {
	count := imageutil.GetScreenNum()
	result := map[string]int{
		"count": count,
	}
	c.JSON(http.StatusOK, result)
	return nil
}


