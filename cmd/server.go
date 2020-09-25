package cmd

import (
	"autoclick/internal/controllers"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	hook "github.com/robotn/gohook"
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

	go func() {

		channel := hook.Start()
		for ev := range channel {
			if ev.Kind == hook.MouseDown {
				log.WithFields(log.Fields{
					"x": ev.X,
					"y": ev.Y,
				}).Info("mouse down")
			}
		}
	}()

	if isDebug {
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stdout)
	} else {
		controllers.TOKEN = token
	}

	router := gin.Default()

	router.Use(glog.Logger(log.StandardLogger()), gin.Recovery())

	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "world")
	})

	controllers.ConfigRouter(router)
	/*router.GET("/screeninfo", routehandler.ScreenInfoHandler)
	router.GET("/screencapture/:num", routehandler.ScreenCapture)
	router.GET("/imgcapture", routehandler.ImageCapture)*/
	return router.Run("127.0.0.1:" + port)
}
