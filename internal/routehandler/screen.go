package routehandler

import (
	"autoclick/model"
	"bytes"
	"image/jpeg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kbinani/screenshot"
	"github.com/pingcap/errors"
	log "github.com/sirupsen/logrus"
)

func ScreenInfoHandler(c *gin.Context) {
	disPlayNum := screenshot.NumActiveDisplays()

	screens := []model.ScreenInfo{}
	for i := 0; i < disPlayNum; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		screen := model.ScreenInfo{
			Num:    i,
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
			Axis: model.Axis{
				Left:   bounds.Min.X,
				Top:    bounds.Min.Y,
				Right:  bounds.Max.X,
				Bottom: bounds.Max.Y,
			},
		}

		screens = append(screens, screen)
	}
	c.JSON(http.StatusOK, screens)
}

func ScreenCapture(c *gin.Context) {
	numstr := c.Param("num")
	num, err := strconv.Atoi(numstr)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Errorln("error screen numer")
		c.AbortWithError(403, errors.Wrap(err, "error screen number"))
		return
	}

	img, err := screenshot.CaptureDisplay(num)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Errorln("error capture screen")
		c.AbortWithError(403, errors.Wrap(err, "error capture display"))
		return
	}

	buffer := new(bytes.Buffer)

	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Errorln("error ")
		c.AbortWithError(403, errors.Wrap(err, "error encode jpeg"))
		return
	}

	c.Writer.Header().Set("Content-Type", "image/jpeg")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := c.Writer.Write(buffer.Bytes()); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Errorln("error ")
		c.AbortWithError(403, errors.Wrap(err, "error response jpeg"))
	}

}
