package fsutil

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func ReadImage(imagePath string) (image.Image, error) {
	if imagePath != "" {
		file, err := os.Open(imagePath)
		defer file.Close()
		if err != nil {
			return nil, err
		}
		i, _, err := image.Decode(file)
		return i, err
	} else {
		return nil, fmt.Errorf("image path is empty")
	}
}

func WriteJpegImageToFile(path string, img image.Image) error {

	buffer := new(bytes.Buffer)

	err := jpeg.Encode(buffer, img, nil)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return nil

}
