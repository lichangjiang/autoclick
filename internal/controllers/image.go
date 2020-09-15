package controllers

import (
	"autoclick/pkg/imageutil"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ImageController struct {
}

func (*ImageController) routeGroup() string {
	return "/img"
}

func (*ImageController) config(rg *gin.RouterGroup) {
	rg.GET("/screencapture/:num",wrapper(screenCapture))
	rg.GET("/imagecapture",wrapper(imageCapture))
}

func screenCapture(c *gin.Context) error {
	numstr := c.Param("num")
	num, err := strconv.Atoi(numstr)
	if err != nil {
		return ParameterError("screen number is not int")
	}

	count := imageutil.GetScreenNum()
	if num < 0 || num >= count {
		err := fmt.Errorf("screen num:%d is out of range ", num)
		return ParameterError(err.Error())
	}

	img, err := imageutil.CaptureScreen(num)
	if err != nil {
		return err
	}

	return responseImage(img, c)
}

func imageCapture(c *gin.Context) error {
	l, errl := strconv.Atoi(c.Query("l"))
	t, errt := strconv.Atoi(c.Query("t"))

	if errl != nil || errt != nil {
		err := fmt.Errorf("left or top should be int")
		return ParameterError(err.Error())
	}

	r, b := 0, 0

	if c.Query("size") != "" {
		size, errz := strconv.Atoi(c.Query("size"))
		if errz != nil {
			err := fmt.Errorf("size should be int")
			return ParameterError(err.Error())
		}
		r = l + size
		b = t + size
	} else if c.Query("r") != "" && c.Query("b") != "" {
		right, errr := strconv.Atoi(c.Query("r"))
		bottom, errb := strconv.Atoi(c.Query("b"))
		if errr != nil || errb != nil {
			err := fmt.Errorf("r or b should be int")
			return ParameterError(err.Error())
		}
		r = right
		b = bottom
	} else {
		err := fmt.Errorf("size not set")
		return ParameterError(err.Error())
	}

	img, err := imageutil.CaptureImage(l, t, r, b)

	if err != nil {
		return err
	}
	return responseImage(img, c)
}

func responseImage(img image.Image, c *gin.Context) error {
	buffer := new(bytes.Buffer)

	err := jpeg.Encode(buffer, img, nil)
	if err != nil {
		return errors.Wrap(err, "error encode jpeg")
	}

	c.Writer.Header().Set("Content-Type", "image/jpeg")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := c.Writer.Write(buffer.Bytes()); err != nil {
		return errors.Wrap(err, "error response jpeg")
	}
	return nil
}
