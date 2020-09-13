package cmd

import (
	"autoclick/internal/routehandler"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	glog "github.com/toorop/gin-logrus"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "autoclick server",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var token = ""
var port = "9000"
var isDebug bool

func init() {
	ServerCmd.RunE = handleServer
	ServerCmd.Flags().BoolVarP(&isDebug, "debug", "d", false, "is debug environment")
	ServerCmd.Flags().StringVarP(&port, "port", "p", "9000", "set http server listen port")
	ServerCmd.Flags().StringVarP(&token, "token", "t", "token", "set http server token header")

}

func handleServer(cmd *cobra.Command, args []string) error {
	if isDebug {
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stdout)
	}

	logger := log.New()
	file, err := os.OpenFile("gin-logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.SetOutput(os.Stdout)
		logger.Info("Failed to log to gin log file")
	}
	if isDebug {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	router := gin.Default()

	router.Use(glog.Logger(logger), gin.Recovery())

	if !isDebug {
		router.Use(func(c *gin.Context) {
			t := c.Request.Header.Get("AppToken")
			if t != token {
				c.AbortWithStatus(403)
			} else {
				c.Next()
			}
		})
	}

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "world")
	})

	router.GET("/screeninfo", routehandler.ScreenInfoHandler)
	router.GET("/screencapture/:num", routehandler.ScreenCapture)
	return router.Run("127.0.0.1:" + port)
}
