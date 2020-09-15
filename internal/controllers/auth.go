package controllers

import "github.com/gin-gonic/gin"

var TOKEN = ""

 func authMiddle(c *gin.Context) {
		t := c.Request.Header.Get("AppToken")
		if t != TOKEN {
			c.AbortWithStatus(403)
		} else {
			c.Next()
		}
}
