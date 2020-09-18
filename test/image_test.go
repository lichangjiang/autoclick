package test

import (
	"autoclick/model"
	"autoclick/pkg/fsutil"
	"autoclick/pkg/imageutil"
	"path"
	"testing"

	"github.com/Nr90/imgsim"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func TestJpegDecode(t *testing.T) {
	_, err := fsutil.ReadImage(path.Join("img", "event1600356746000-313.jpeg"))
	assert.NilError(t, err)
}

func TestSimiliar(t *testing.T) {
	event1 := model.Event{
		Name: "125",
		Axis: model.Axis{
			Left:   125,
			Right:  145,
			Bottom: 90,
			Top:    70,
		},
	}

	event2 := model.Event{
		Name: "314",
		Axis: model.Axis{
			Left:   314,
			Right:  334,
			Bottom: 551,
			Top:    531,
		},
	}

	img1, err := imageutil.CaptureImage(event1.Left,
		event1.Top,
		event1.Right,
		event1.Bottom)

	assert.NilError(t, err)


	img1Hash := imgsim.AverageHash(img1)
	t.Log(img1Hash)

	img2, err := imageutil.CaptureImage(event2.Left,
		event2.Top,
		event2.Right,
		event2.Bottom)

	assert.NilError(t, err)

	img2Hash := imgsim.AverageHash(img2)
	t.Log(img2Hash)

	fimg1, err := fsutil.ReadImage(path.Join("img", "event1600402492000-125.jpeg"))
	assert.NilError(t, err)

	img3Hash := imgsim.AverageHash(fimg1)
	t.Log(img3Hash)

	fimg2, err := fsutil.ReadImage(path.Join("img", "event1600402497000-314.jpeg"))
	assert.NilError(t, err)

	img4Hash := imgsim.AverageHash(fimg2)
	t.Log(img4Hash)

	t.Log(imgsim.Distance(img1Hash,img3Hash))
	t.Log(imgsim.Distance(img2Hash,img4Hash))
	
	assert.Assert(t, cmp.Equal(true, imageutil.ImageSimiliar(img1, img1)))
	assert.Assert(t, cmp.Equal(true, imageutil.ImageSimiliar(img2, img2)))
	assert.Assert(t, cmp.Equal(true, imageutil.ImageSimiliar(fimg1, fimg1)))
	assert.Assert(t, cmp.Equal(true, imageutil.ImageSimiliar(fimg2, fimg2)))

	assert.Assert(t, cmp.Equal(true, imageutil.ImageSimiliar(img1, fimg1)))
	assert.Assert(t, cmp.Equal(true, imageutil.ImageSimiliar(img2, fimg2)))
}
