package controllers

import (
	"autoclick/model"
	"autoclick/pkg/fsutil"
	"autoclick/pkg/imageutil"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var isStarted atomic.Value
var startNum atomic.Value

func init() {
	isStarted.Store(false)
	startNum.Store(0)
}

type ActionController struct {
}

func (*ActionController) routeGroup() string {
	return "/action"
}

func (*ActionController) config(rg *gin.RouterGroup) {

	rg.Use(authMiddle)
	rg.POST("/start", wrapper(startAction))
	rg.POST("/stop", wrapper(stopAction))
}

func startAction(c *gin.Context) error {
	started := isStarted.Load()
	num := startNum.Load().(int)
	if started == false || started == nil {
		isStarted.Store(true)
		startNum.Store(num + 1)
		startMsg := model.StartMsg{}
		if err := c.BindJSON(&startMsg); err != nil {
			return err
		}

		go action(startMsg.EventGroups, startMsg.EventGroupTimeInterval, startMsg.EventTimeInterval, startMsg.ShowMouse, num+1)
	}
	c.JSON(http.StatusOK, struct{}{})
	return nil
}

func stopAction(c *gin.Context) error {
	started := isStarted.Load()
	if started == true {
		isStarted.Store(false)
	}
	c.JSON(http.StatusOK, struct{}{})
	return nil
}

func action(eventGroups [][]model.Event,
	eventGroupTimeInterval, eventTimeInterval int,
	showMouse bool, oriStartNum int) {
	logrus.WithFields(logrus.Fields{
		"eventTimeInterval":      eventTimeInterval,
		"eventGroupTimeInterval": eventGroupTimeInterval,
		"eventGreoups":           eventGroups,
		"showMouse":              showMouse,
	}).Info("action request arguments")
	newEventGroups := [][]model.Event{}

	for _, eventGroup := range eventGroups {
		newEventGroup := []model.Event{}
		for _, event := range eventGroup {
			if event.Path != "" {
				img, err := fsutil.ReadImage(event.Path)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"name":  event.Name,
						"path":  event.Path,
						"error": err,
					}).Error("read image error")
				} else {
					event.Image = img
				}
			} else {
				logrus.WithFields(logrus.Fields{
					"name": event.Name,
				}).Error("file path is empty")
			}

			newEventGroup = append(newEventGroup, event)
		}
		newEventGroups = append(newEventGroups, newEventGroup)
	}

	for {
		started := isStarted.Load()
		startNum := startNum.Load()
		if started == true && startNum == oriStartNum {
			if err := imageutil.StartImageCheck(newEventGroups,
				eventGroupTimeInterval, eventTimeInterval,
				showMouse); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("image check error")
			}
		} else {
			logrus.Info("stop image check action")
			return
		}
	}
}
